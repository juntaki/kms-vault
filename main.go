package main

import (
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

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
			Usage:    "Google cloud project",
			Value:    config.Project,
			Required: config.Project == "",
		},
	}

	app := cli.NewApp()
	app.Name = "kmsvault"
	app.Usage = "Manage configuration file that partially contain confidential information in a repository using Cloud KMS."
	app.Version = "0.1.0"
	app.EnableBashCompletion = true

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

const configName = ".kmsvault.yaml"

func loadConfig() (config *VaultConfig) {
	config = &VaultConfig{
		Location: "global",
	}

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for ; ; dir = filepath.Dir(dir) {
		_, err := os.Stat(filepath.Join(dir, configName))
		if err == nil {
			break
		}
		if dir == filepath.Dir(dir) {
			return
		}
	}

	fp, err := os.Open(filepath.Join(dir, configName))
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
