package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks log messages for complience with formatting and security rules",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var sensitiveDataRegex = regexp.MustCompile(`(?i)(password|api_key|token|secret)\s*[:=]`)

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

		if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msgString, _ = strconv.Unquote(lit.Value)
			errorPos = lit.Pos()
		} else if bin, ok := firstArg.(*ast.BinaryExpr); ok && bin.Op == token.ADD {
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
			if unicode.Is(unicode.Cyrillic, r) {
				pass.Reportf(errorPos, "log message must be in English only")
				break
			}
		}

		// Лог-сообщения не должны содержать спецсимволы или эмодзи
		hasSpecial := false
		if strings.Contains(msgString, "...") {
			hasSpecial = true
		} else {
			for _, r := range runes {
				if r == '!' || r == '?' || unicode.Is(unicode.So, r) {
					hasSpecial = true
					break
				}
			}
		}
		if hasSpecial {
			pass.Reportf(errorPos, "log message must not contain special characters or emojis")
		}

		// Чувствительные данные
		if sensitiveDataRegex.MatchString(msgString) {
			pass.Reportf(errorPos, "log message must not contain sensitive data")
		}
	})
	return nil, nil
}
