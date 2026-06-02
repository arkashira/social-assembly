package moderation

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/axentx/axentx-social-assembly/pkg/config"
)

func GetRulesPath() string {
	return filepath.Join(config.GetConfigPath(), "rules.yaml")
}

func LoadRules() error {
	rulesPath := GetRulesPath()
	return NewRuleEngine().LoadRules(rulesPath)
}

func GetRule(key string) (interface{}, bool) {
	return NewRuleEngine().GetRule(key)
}

func UpdateRule(key string, value interface{}) error {
	return NewRuleEngine().UpdateRule(key, value)
}

func TriggerAction(key string, action string) error {
	return NewRuleEngine().TriggerAction(key, action)
}