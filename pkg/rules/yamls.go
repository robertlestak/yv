package rules

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ParseYamls(yr io.Reader) ([]map[string]any, error) {
	var yamls []map[string]any
	d := yaml.NewDecoder(yr)
	for {
		spec := new(map[string]any)
		err := d.Decode(&spec)
		if spec == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return yamls, err
		}
		yamls = append(yamls, *spec)
	}
	return yamls, nil
}

func ParseYamlFiles(f string) ([]map[string]any, error) {
	// recursively parse all yaml files in f
	// if f is a directory, loop through all files and load them
	var yamls []map[string]any
	err := filepath.Walk(f,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".yaml" {
				return nil
			}
			fd, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fd.Close()
			yaml, err := ParseYamls(fd)
			if err != nil {
				return err
			}
			yamls = append(yamls, yaml...)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return yamls, nil
}
