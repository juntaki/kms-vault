package main

import (
	"os"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func configCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:  "config",
		Usage: "Create .kms-vault.yaml",
		Flags: append([]cli.Flag{
			cli.BoolFlag{
				Name:  "w",
				Usage: "Write to .kms-vault.yaml.",
			},
		}, kmsFlags...),
		Action: configAction,
	}
}

type VaultConfig struct {
	Project  string
	Location string
	KeyRing  string
	Key      string
}

func configAction(c *cli.Context) error {
	val, err := yaml.Marshal(VaultConfig{
		Project:  c.String("project"),
		Location: c.String("location"),
		KeyRing:  c.String("keyring"),
		Key:      c.String("key"),
	})

	if err != nil {
		return err
	}

	output := os.Stdout
	if c.Bool("w") {
		err := checkOverwrite(vaultConfigFilename)
		if err != nil {
			return err
		}
		fp, err := os.Create(vaultConfigFilename)
		if err != nil {
			return err
		}
		defer fp.Close()
		output = fp
	}

	_, err = output.Write(val)
	if err != nil {
		return err
	}
	return nil
}
