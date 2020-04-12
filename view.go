package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

func viewCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:  "view",
		Usage: "View decrypted files",
		Flags: append([]cli.Flag{
			cli.BoolFlag{
				Name:  "yaml",
				Usage: "View the parsed data to fill template",
			},
		}, kmsFlags...),
		Action: viewAction,
	}
}

func viewAction(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return xerrors.New("Specify at least one file")
	}

	name := kmsNameFromContext(c)

	raw := make(map[string][]byte)
	for _, filename := range c.Args() {
		// Skip dir
		fstat, err := os.Stat(filename)
		if err != nil {
			return err
		}
		if fstat.IsDir() {
			continue
		}

		plainText, err := getPlainText(name, filename)
		if xerrors.Is(err, InvalidFormatError) {
			plainText, err = ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		raw[filename] = plainText
	}
	if len(raw) == 0 {
		return xerrors.New("Specify at least one file")
	}

	if c.Bool("yaml") {
		result, err := convertToTemplateData(raw)
		if err != nil {
			return err
		}
		val, err := yaml.Marshal(result)
		if err != nil {
			return err
		}
		os.Stdout.Write(val)
		return nil
	}
	if len(raw) == 1 {
		for _, val := range raw {
			fmt.Print(string(val))
		}
	} else {
		for filename, val := range raw {
			fmt.Println(filename)
			fmt.Println(string(val))
		}
	}

	return nil
}
