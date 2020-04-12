package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

func fillCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:  "fill",
		Usage: "Fill template by files",
		Flags: append([]cli.Flag{
			cli.StringFlag{
				Name:     "template",
				Usage:    "",
				Required: true,
			},
			cli.StringFlag{
				Name:  "output",
				Usage: "",
			},
		}, kmsFlags...),
		Action: fillAction,
	}
}

func fillAction(c *cli.Context) error {
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

	result, err := convertToTemplateData(raw)
	if err != nil {
		return err
	}

	tmplString, err := ioutil.ReadFile(c.String("template-file"))
	if err != nil {
		return err
	}
	tmpl, err := template.New("vault-template").Parse(string(tmplString))
	if err != nil {
		return err
	}

	output := os.Stdout
	if outputFile := c.String("output"); outputFile != "" {
		_, err := os.Stat(outputFile)
		if err == nil {
			fmt.Printf("File exist, overwrite? (y/N): ")
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()
			text := stdin.Text()
			if !(len(text) > 0 && strings.ToLower(strings.TrimSpace(text))[0] == 'y') {
				log.Println("Aborted")
				return nil
			}
		}

		fp, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer fp.Close()
		output = fp
	}

	err = tmpl.Execute(output, result)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

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
				return nil, xerrors.Errorf("duplicate filename: %s", templateDataID(filename))
			}
			result[templateDataID(filename)] = d
		} else {
			if _, ok := result[templateDataID(filename)]; ok {
				return nil, xerrors.Errorf("duplicate filename: %s", templateDataID(filename))
			}
			result[templateDataID(filename)] = string(plainText)
		}
	}
	return result, nil
}

func getPlainText(name, filename string) ([]byte, error) {
	fp, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, xerrors.Errorf("open: %w", err)
	}
	defer fp.Close()

	plainText, err := decryptFile(name, fp)
	if err != nil {
		return nil, xerrors.Errorf("decryptFileAndPrint: %w", err)
	}
	return plainText, nil
}
