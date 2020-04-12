package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const vaultConfigFilename = ".kms-vault.yaml"

func loadConfig() (config *VaultConfig) {
	config = &VaultConfig{
		Location: "global",
	}

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for ; ; dir = filepath.Dir(dir) {
		_, err := os.Stat(filepath.Join(dir, vaultConfigFilename))
		if err == nil {
			break
		}
		if dir == filepath.Dir(dir) {
			return
		}
	}

	fp, err := os.Open(filepath.Join(dir, vaultConfigFilename))
	if err != nil {
		return
	}
	defer fp.Close()
	d := yaml.NewDecoder(fp)
	err = d.Decode(config)
	if err != nil {
		return
	}
	return
}
