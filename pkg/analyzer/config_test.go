package analyzer

import "testing"

func TestIsRuleEnabled(t *testing.T) {
	// Сохраняем и восстанавливаем исходное состояние
	origDisabled := disabledRules
	defer func() { disabledRules = origDisabled }()

	t.Run("all rules enabled by default", func(t *testing.T) {
		disabledRules = nil

		rules := []string{RuleLowercase, RuleEnglishOnly, RuleSpecialChars, RuleSensitiveData}
		for _, r := range rules {
			if !isRuleEnabled(r) {
				t.Errorf("expected rule %q to be enabled, but it was disabled", r)
			}
		}
	})

	t.Run("disabled rule returns false", func(t *testing.T) {
		disabledRules = []string{RuleLowercase, RuleSensitiveData}

		if isRuleEnabled(RuleLowercase) {
			t.Error("expected rule lowercase to be disabled")
		}
		if isRuleEnabled(RuleSensitiveData) {
			t.Error("expected rule sensitive-data to be disabled")
		}
	})

	t.Run("non-disabled rule returns true", func(t *testing.T) {
		disabledRules = []string{RuleLowercase}

		if !isRuleEnabled(RuleEnglishOnly) {
			t.Error("expected rule english-only to be enabled")
		}
		if !isRuleEnabled(RuleSpecialChars) {
			t.Error("expected rule special-chars to be enabled")
		}
	})
}
