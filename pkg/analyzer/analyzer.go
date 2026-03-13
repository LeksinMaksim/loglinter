package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	register.Plugin("loglinter", New)
	Analyzer.Flags.StringVar(&cliCustomPatterns, "custom-patterns", "", "comma-separated list of custom patterns for sensitive data checking")
}

var (
	pluginCustomPatterns     []string
	cliCustomPatterns        string
	activeSensitiveDataRegex *regexp.Regexp
	regexOnce                sync.Once
)

type LogLinterPlugin struct{}

func New(conf any) (register.LinterPlugin, error) {
	// conf - это map[string]any, куда golangci-lint передает всё из блока linters-settings.
	// Мы пытаемся найти наш массив `custom-patterns` и сохранить паттерны.
	if confMap, ok := conf.(map[string]any); ok {
		if customPatterns, ok := confMap["custom-patterns"].([]any); ok {
			for _, p := range customPatterns {
				if strP, ok := p.(string); ok {
					pluginCustomPatterns = append(pluginCustomPatterns, strP)
				}
			}
		}
	}
	return &LogLinterPlugin{}, nil
}

func (p *LogLinterPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

func (p *LogLinterPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks log messages for complience with formatting and security rules",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var sensitiveDataRegex = regexp.MustCompile(`(?i)(password|api_key|token|secret)\s*[:=]`)
var allowedChars = regexp.MustCompile(`^[\p{L}0-9\s.,:=\-\[\]_]+$`)

func getSensitiveDataRegex() *regexp.Regexp {
	// Используем sync.Once, чтобы скомпилировать регулярное выражение
	// только один раз при первой проверке
	regexOnce.Do(func() {
		var patterns []string

		// 1. Добавляем паттерны, которые пришли из настроек .golangci.yml
		if len(pluginCustomPatterns) > 0 {
			patterns = append(patterns, pluginCustomPatterns...)
		}

		// 2. Добавляем паттерны, которые передали через командную строку (-custom-patterns="a,b")
		if cliCustomPatterns != "" {
			parts := strings.Split(cliCustomPatterns, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					patterns = append(patterns, p)
				}
			}
		}

		// Если собственных паттернов нет, используем дефолтные
		if len(patterns) == 0 {
			activeSensitiveDataRegex = sensitiveDataRegex
			return
		}

		// Если есть кастомные паттерны, собираем из них новую регулярку
		// вида: (?i)(pattern1|pattern2)\s*[:=]
		patternStr := `(?i)(` + strings.Join(patterns, "|") + `)\s*[:=]`
		activeSensitiveDataRegex = regexp.MustCompile(patternStr)
	})
	return activeSensitiveDataRegex
}

var logMethods = map[string]bool{
	"Info": true, "Infof": true,
	"Error": true, "Errorf": true,
	"Warn": true, "Warnf": true,
	"Debug": true, "Debugf": true,
	"Fatal": true, "Fatalf": true,
	"Panic": true, "Panicf": true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		if !logMethods[sel.Sel.Name] {
			return
		}

		// Линтер работает только с log/slog и go.uber.org/zap
		if obj, ok := pass.TypesInfo.Uses[sel.Sel].(*types.Func); ok {
			pkg := obj.Pkg()
			if pkg == nil {
				return
			}

			pkgPath := pkg.Path()

			if pkgPath != "log/slog" && pkgPath != "go.uber.org/zap" {
				return
			}
		} else {
			return
		}

		if len(call.Args) == 0 {
			return
		}

		var msgString string
		var errorPos token.Pos

		firstArg := call.Args[0]

		// Передается строковый литерал
		if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msgString, _ = strconv.Unquote(lit.Value)
			errorPos = lit.Pos()
		} else if bin, ok := firstArg.(*ast.BinaryExpr); ok && bin.Op == token.ADD { // Передается строковый литерал + переменная
			if lit, ok := bin.X.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				msgString, _ = strconv.Unquote(lit.Value)
				errorPos = bin.X.Pos()
			}
		}

		if msgString == "" {
			return
		}

		runes := []rune(msgString)

		// Лог-сообщения начинаются со строчной буквы
		firstChar := runes[0]
		if unicode.IsLetter(firstChar) && !unicode.IsLower(firstChar) {
			pass.Reportf(errorPos, "log message must start with a lowercase letter")
		}

		// Лог-сообщения должны быть только на английском языке
		for _, r := range runes {
			if unicode.IsLetter(r) && r > unicode.MaxASCII {
				pass.Reportf(errorPos, "log message must be in English only")
				break
			}
		}

		// Лог-сообщения не должны содержать спецсимволы или эмодзи
		hasSpecial := false
		if strings.Contains(msgString, "...") {
			hasSpecial = true
		}

		if !allowedChars.MatchString(msgString) {
			hasSpecial = true
		}

		if hasSpecial {
			pass.Reportf(errorPos, "log message must not contain special characters or emojis")
		}

		// Чувствительные данные
		if getSensitiveDataRegex().MatchString(msgString) {
			pass.Reportf(errorPos, "log message must not contain sensitive data")
		}
	})
	return nil, nil
}
