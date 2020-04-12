package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const VaultVersion = "0.1.0"

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	config := loadConfig()
	kmsFlags := []cli.Flag{
		cli.StringFlag{
			Name:     "key",
			Usage:    "The key to use for encryption.",
			Value:    config.Key,
			Required: config.Key == "",
		},
		cli.StringFlag{
			Name:     "keyring",
			Usage:    "Key ring of the key.",
			Value:    config.KeyRing,
			Required: config.KeyRing == "",
		},
		cli.StringFlag{
			Name:     "location",
			Usage:    "Location of the keyring.",
			Value:    config.Location,
			Required: config.Location == "",
		},
		cli.StringFlag{
			Name:     "project",
			Usage:    "Google cloud project name.",
			Value:    config.Project,
			Required: config.Project == "",
		},
	}

	app := cli.NewApp()
	app.Name = "kms-vault"
	app.Usage = "Manage configuration file that partially contain confidential information in a repository using Cloud KMS."

	app.Version = VaultVersion
	app.Authors = []cli.Author{
		{
			Name:  "Jumpei Takiyasu",
			Email: "me@juntaki.com",
		},
	}
	app.Copyright = "(c) 2020 Jumpei Takiyasu"
	app.Commands = []cli.Command{
		encryptCommand(kmsFlags),
		decryptCommand(kmsFlags),
		viewCommand(kmsFlags),
		configCommand(kmsFlags),
		fillCommand(kmsFlags),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

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
