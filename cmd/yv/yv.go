package main

import (
	"flag"
	"os"

	"github.com/robertlestak/yv/pkg/rules"
	log "github.com/sirupsen/logrus"
)

func init() {
	ll, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}

func main() {
	var ruleFilePath string
	var yamlsFilePath string
	flag.StringVar(&ruleFilePath, "r", "", "path to the rule file(s)")
	flag.StringVar(&yamlsFilePath, "f", "", "path to the yaml file(s)")
	flag.Parse()
	if ruleFilePath == "" || yamlsFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	r, err := rules.ParseRulesFromFile(ruleFilePath)
	if err != nil {
		log.Fatal(err)
	}
	y, err := rules.ParseYamlFiles(yamlsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, yaml := range y {
		for _, rule := range r {
			if err := rules.Validate(yaml, rule); err != nil {
				log.Fatal(err)
			}
		}
	}
}
