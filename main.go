package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Vault"
	app.Usage = "Manage configuration file that partially contain confidential information in a repository using Cloud KMS."
	app.Version = "0.1.0"
	app.EnableBashCompletion = true

	kmsFlags := []cli.Flag{
		cli.StringFlag{
			Name:     "key",
			Usage:    "The key to use for encryption.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "keyring",
			Usage:    "Key ring of the key.",
			Required: true,
		},
		cli.StringFlag{
			Name:  "location",
			Usage: "Location of the keyring.",
			Value: "global",
		},
		cli.StringFlag{
			Name:     "project",
			Usage:    "Google cloud project",
			Required: true,
		},
	}

	app.Flags = append([]cli.Flag{}, kmsFlags...)
	app.Commands = []cli.Command{
		encryptCommand,
		decryptCommand,
	}
	app.Run(os.Args)
}
