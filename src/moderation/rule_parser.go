// Package moderation provides content moderation rule parsing and validation
package moderation

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"gopkg.in/yaml.v3"
)

// RuleType enumerates supported rule types
type RuleType string

const (
	RuleTypeWordFilter RuleType = "word_filter"
	RuleTypeRateLimit  RuleType = "rate_limit"
	RuleTypeUserBan    RuleType = "user_ban"
)

// ActionType enumerates automatic actions on violation
type ActionType string

const (
	ActionHidePost  ActionType = "hide_post"
	ActionTempMute  ActionType = "temp_mute"
	ActionPermBan   ActionType = "perm_ban"
)

// WordFilterRule defines a word filter rule
type WordFilterRule struct {
	Patterns      []string `yaml:"patterns"`
	CaseSensitive bool     `yaml:"case_sensitive"`
	Action        ActionType `yaml:"action"`
	Duration      string    `yaml:"duration"` // e.g., "24h"
}

// RateLimitRule defines a rate limit rule
type RateLimitRule struct {
	MaxActions int       `yaml:"max_actions"`
	Window     string    `yaml:"window"` // e.g., "5m"
	Action     ActionType `yaml:"action"`
	Duration   string    `yaml:"duration"`
}

// UserBanRule defines a user ban rule
type UserBanRule struct {
	Usernames []string  `yaml:"usernames"`
	Reason    string    `yaml:"reason"`
	Action    ActionType `yaml:"action"`
	Duration  string    `yaml:"duration"`
}

// Ruleset represents the top-level YAML ruleset
type Ruleset struct {
	Version string `yaml:"version"`
	Rules   []Rule `yaml:"rules"`
}

// Rule is a union type for all rule kinds
type Rule struct {
	Type        RuleType       `yaml:"type"`
	WordFilter  *WordFilterRule `yaml:"word_filter,omitempty"`
	RateLimit   *RateLimitRule  `yaml:"rate_limit,omitempty"`
	UserBan     *UserBanRule    `yaml:"user_ban,omitempty"`
}

// ParseYAML parses a YAML file at path into a Ruleset
func ParseYAML(path string) (*Ruleset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file %q: %w", path, err)
	}
	return ParseYAMLBytes(data)
}

// ParseYAMLBytes parses YAML bytes into a Ruleset
func ParseYAMLBytes(data []byte) (*Ruleset, error) {
	var rs Ruleset
	if err := yaml.Unmarshal(data, &rs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules YAML: %w", err)
	}

	if rs.Version != "1" {
		return nil, errors.New("unsupported ruleset version; expected '1'")
	}

	if len(rs.Rules) == 0 {
		return nil, errors.New("ruleset contains no rules")
	}

	for i, r := range rs.Rules {
		if r.Type == "" {
			return nil, fmt.Errorf("rule at index %d missing 'type'", i)
		}

		switch r.Type {
		case RuleTypeWordFilter:
			if r.WordFilter == nil {
				return nil, fmt.Errorf("rule at index %d type=word_filter missing word_filter block", i)
			}
			if len(r.WordFilter.Patterns) == 0 {
				return nil, fmt.Errorf("rule at index %d word_filter.patterns must be non-empty", i)
			}
			if r.WordFilter.Action == "" {
				return nil, fmt.Errorf("rule at index %d word_filter.action must be set", i)
			}
			if _, err := parseDuration(r.WordFilter.Duration); err != nil {
				return nil, fmt.Errorf("rule at index %d invalid word_filter.duration: %v", i, err)
			}

		case RuleTypeRateLimit:
			if r.RateLimit == nil {
				return nil, fmt.Errorf("rule at index %d type=rate_limit missing rate_limit block", i)
			}
			if r.RateLimit.MaxActions <= 0 {
				return nil, fmt.Errorf("rule at index %d rate_limit.max_actions must be > 0", i)
			}
			if r.RateLimit.Window == "" {
				return nil, fmt.Errorf("rule at index %d rate_limit.window must be set", i)
			}
			if _, err := parseDuration(r.RateLimit.Window); err != nil {
				return nil, fmt.Errorf("rule at index %d invalid rate_limit.window: %v", i, err)
			}
			if r.RateLimit.Action == "" {
				return nil, fmt.Errorf("rule at index %d rate_limit.action must be set", i)
			}
			if _, err := parseDuration(r.RateLimit.Duration); err != nil {
				return nil, fmt.Errorf("rule at index %d invalid rate_limit.duration: %v", i, err)
			}

		case RuleTypeUserBan:
			if r.UserBan == nil {
				return nil, fmt.Errorf("rule at index %d type=user_ban missing user_ban block", i)
			}
			if len(r.UserBan.Usernames) == 0 {
				return nil, fmt.Errorf("rule at index %d user_ban.usernames must be non-empty", i)
			}
			if r.UserBan.Action == "" {
				return nil, fmt.Errorf("rule at index %d user_ban.action must be set", i)
			}
			if _, err := parseDuration(r.UserBan.Duration); err != nil {
				return nil, fmt.Errorf("rule at index %d invalid user_ban.duration: %v", i, err)
			}

		default:
			return nil, fmt.Errorf("rule at index %d has unknown type %q", i, r.Type)
		}
	}

	return &rs, nil
}

// parseDuration parses a duration string (e.g., "5m", "1h", "24h") into time.Duration
func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("duration cannot be empty")
	}

	if strings.HasSuffix(s, "m") {
		val := strings.TrimSuffix(s, "m")
		if val == "" {
			return 0, errors.New("invalid minutes value")
		}
		return time.ParseDuration(val + "m")
	}

	if strings.HasSuffix(s, "h") {
		val := strings.TrimSuffix(s, "h")
		if val == "" {
			return 0, errors.New("invalid hours value")
		}
		return time.ParseDuration(val + "h")
	}

	if strings.HasSuffix(s, "d") {
		val := strings.TrimSuffix(s, "d")
		if val == "" {
			return 0, errors.New("invalid days value")
		}
		return time.ParseDuration(val + "h") * 24
	}

	return time.ParseDuration(s)
}

// ParseReader parses rules from an io.Reader
func ParseReader(r io.Reader) (*Ruleset, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseYAMLBytes(data)
}