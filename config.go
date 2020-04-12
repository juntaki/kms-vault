package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func configCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "config",
		Usage:  "Create .kmsvault.yaml config file",
		Flags:  kmsFlags,
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

	fmt.Println(string(val))
	return nil
}
