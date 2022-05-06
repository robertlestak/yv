package rules

import (
	"bytes"
	"fmt"

	gyaml "github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	uniqueValues map[string][]string
)

func validateGloballyUniqueValues(y map[string]any, rule any) error {
	l := log.WithFields(log.Fields{
		"rule":       "globally-unique-values",
		"constraint": rule,
		"y":          y,
	})
	l.Debug("validating globally-unique-values rule")
	var r GloballyUniqueRule
	yd, err := yaml.Marshal(rule)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yd, &r)
	if err != nil {
		return err
	}
	l.WithFields(log.Fields{
		"path": r.Path,
	}).Debug("validating globally-unique-values rule")
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
		return err
	}
	l = l.WithFields(log.Fields{
		"yamlValue": pv,
	})
	l.Debug("validating globally-unique-values rule")
	if uniqueValues == nil {
		uniqueValues = make(map[string][]string)
	}
	for _, v := range uniqueValues[r.Path] {
		if v == pv.(string) {
			return fmt.Errorf("globally-unique-values violation: %v", rule)
		}
	}
	if _, ok := uniqueValues[r.Path]; !ok {
		l.Debug("setting globally-unique-values rule")
		uniqueValues[r.Path] = []string{
			pv.(string),
		}
	} else {
		uniqueValues[r.Path] = append(uniqueValues[r.Path], pv.(string))
	}
	return nil
}
