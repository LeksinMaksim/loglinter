package analyzer

// Имена правил, которые можно отключить через конфигурацию.
const (
	RuleLowercase     = "lowercase"
	RuleEnglishOnly   = "english-only"
	RuleSpecialChars  = "special-chars"
	RuleSensitiveData = "sensitive-data"
)

// disabledRules содержит список отключённых правил.
// Заполняется из CLI-флага или настроек golangci-lint.
var disabledRules []string

// isRuleEnabled возвращает true, если правило не находится в списке отключённых.
func isRuleEnabled(ruleName string) bool {
	for _, r := range disabledRules {
		if r == ruleName {
			return false
		}
	}
	return true
}
