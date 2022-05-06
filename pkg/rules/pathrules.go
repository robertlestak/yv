package rules

import (
	"bytes"
	"fmt"
	"regexp"

	gyaml "github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (r RestrictedPathRule) checkValue(v any) error {
	if r.Regex {
		rx := regexp.MustCompile(r.Value)
		if rx.MatchString(v.(string)) {
			return fmt.Errorf("restricted path violation: %v", r)
		}
	} else {
		if v.(string) == r.Value {
			return fmt.Errorf("restricted path violation: %v", r)
		}
	}
	return nil
}

func validateRestrictedPathValues(y map[string]any, rule any) error {
	l := log.WithFields(log.Fields{
		"rule": "restricted-path-values",
		"y":    y,
	})
	l.Debug("validating restricted-path-values rule")
	var r RestrictedPathRule
	yd, err := yaml.Marshal(rule)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yd, &r)
	if err != nil {
		return err
	}
	l.WithFields(log.Fields{
		"path":  r.Path,
		"value": r.Value,
		"regex": r.Regex,
	}).Debug("validating restricted-path-values rule")
	path, err := gyaml.PathString("$." + r.Path)
	if err != nil {
		l.WithFields(log.Fields{
			"error": err,
		}).Error("failed to parse path")
		return err
	}
	yd, err = yaml.Marshal(y)
	if err != nil {
		return err
	}
	var pv any
	err = path.Read(bytes.NewReader(yd), &pv)
	if gyaml.IsNotFoundNodeError(err) {
		return nil
	}
	if err != nil {
		l.WithFields(log.Fields{
			"error": err,
		}).Error("failed to read path")
		return err
	}
	l = l.WithFields(log.Fields{
		"yamlValue": pv,
	})
	l.Debug("validating restricted-path-values rule")
	// if pv is an array, check each element
	if lpv, ok := pv.([]interface{}); ok {
		for _, v := range lpv {
			if err := r.checkValue(v); err != nil {
				return err
			}
		}
	} else {
		if err := r.checkValue(pv); err != nil {
			return err
		}
	}
	return nil
}
