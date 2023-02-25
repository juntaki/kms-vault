package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func templateDataID(path string) string {
	// remove ext and base
	out := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
	// lower
	out = strings.ToLower(out)
	return out
}

func convertToTemplateData(raw map[string][]byte) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for filename, plainText := range raw {
		if strings.HasSuffix(filename, ".yaml") ||
			strings.HasSuffix(filename, ".yml") ||
			strings.HasSuffix(filename, ".json") {
			d := make(map[interface{}]interface{})
			err := yaml.Unmarshal(plainText, d)
			if err != nil {
				return nil, err
			}

			if _, ok := result[templateDataID(filename)]; ok {
				return nil, fmt.Errorf("duplicate filename: %s", templateDataID(filename))
			}
			result[templateDataID(filename)] = d
		} else {
			if _, ok := result[templateDataID(filename)]; ok {
				return nil, fmt.Errorf("duplicate filename: %s", templateDataID(filename))
			}
			result[templateDataID(filename)] = string(plainText)
		}
	}
	return result, nil
}
