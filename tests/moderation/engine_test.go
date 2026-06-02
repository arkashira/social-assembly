package moderation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleEngine_LoadRules(t *testing.T) {
	ruleEngine := NewRuleEngine()
	rulesPath := GetRulesPath()
	err := ruleEngine.LoadRules(rulesPath)
	assert.NoError(t, err)
}

func TestRuleEngine_GetRule(t *testing.T) {
	ruleEngine := NewRuleEngine()
	rule, ok := ruleEngine.GetRule("test_rule")
	assert.True(t, ok)
	assert.NotNil(t, rule)
}

func TestRuleEngine_UpdateRule(t *testing.T) {
	ruleEngine := NewRuleEngine()
	err := ruleEngine.UpdateRule("test_rule", "new_value")
	assert.NoError(t, err)
}

func TestRuleEngine_TriggerAction(t *testing.T) {
	ruleEngine := NewRuleEngine()
	err := ruleEngine.TriggerAction("test_rule", "hide_post")
	assert.NoError(t, err)
}