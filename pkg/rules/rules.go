package rules

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type RuleType string

const (
	RuleTypeRestrictedPathValues RuleType = "restricted-path-values"
	RuleTypeGloballyUniqueValues RuleType = "globally-unique-values"
)

type RestrictedPathRule struct {
	Path  string `yaml:"path"`
	Value string `yaml:"value"`
	Regex bool   `yaml:"regex"`
}

type GloballyUniqueRule struct {
	Path string `yaml:"path"`
}

type Rule struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Type        RuleType `yaml:"type"`
	Rule        any      `yaml:"rule"`
}

func ParseRulesFromFile(path string) ([]Rule, error) {
	var rules []Rule
	files, err := ParseYamlFiles(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		var rs struct {
			Rules []Rule `yaml:"rules"`
		}
		jd, err := yaml.Marshal(file)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(jd, &rs)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rs.Rules...)
	}
	return rules, nil
}

func Validate(y map[string]any, rule Rule) error {
	switch rule.Type {
	case RuleTypeRestrictedPathValues:
		return validateRestrictedPathValues(y, rule.Rule)
	case RuleTypeGloballyUniqueValues:
		return validateGloballyUniqueValues(y, rule.Rule)
	default:
		return errors.New("unknown rule type")
	}
}
