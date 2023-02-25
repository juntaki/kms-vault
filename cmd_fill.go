package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/urfave/cli"
)

func fillCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:  "fill",
		Usage: "Fill template by files",
		Flags: append([]cli.Flag{
			cli.StringFlag{
				Name:     "template",
				Usage:    "File path of template file.",
				Required: true,
			},
			cli.StringFlag{
				Name:  "output",
				Usage: "File path of output.",
			},
		}, kmsFlags...),
		Before: initializeKMS,
		Action: fillAction,
	}
}

func fillAction(c *cli.Context) error {
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

		plainText, err := getPlainText(filename)
		if errors.Is(err, ErrorInvalidFormat) {
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
		return fmt.Errorf("specify at least one file")
	}

	result, err := convertToTemplateData(raw)
	if err != nil {
		return err
	}

	output := os.Stdout
	if outputFile := c.String("output"); outputFile != "" {
		err := checkOverwrite(outputFile)
		if err != nil {
			return err
		}
		fp, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer fp.Close()
		output = fp
	}

	tmplString, err := ioutil.ReadFile(c.String("template"))
	if err != nil {
		return err
	}
	tmpl, err := template.New("vault-template").Parse(string(tmplString))
	if err != nil {
		return err
	}
	err = tmpl.Execute(output, result)
	if err != nil {
		return err
	}

	return nil
}
