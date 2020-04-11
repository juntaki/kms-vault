package main

import (
	cloudkms "cloud.google.com/go/kms/apiv1"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func encrypt(projectID, location, ringID, keyID string, plaintext []byte) (string, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectID, location, ringID, keyID)

	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", xerrors.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      name,
		Plaintext: []byte(plaintext),
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return "", xerrors.Errorf("encrypt: %v", err)
	}

	return base64.StdEncoding.EncodeToString(resp.Ciphertext), nil
}

func main() {
	var (
		suffix string
	)
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

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "suffix, s",
			Value:       "!!!",
			Usage:       "text after speaking something",
			Destination: &suffix,
			EnvVar:      "SUFFIX",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "kms-encrypt",
			Usage: "fdsafs",
			Flags: append([]cli.Flag{
				cli.StringFlag{
					Name:     "out-dir",
					Usage:    "Output directory path to encrypt",
					Required: true,
				},
			}, kmsFlags...),

			Action: func(c *cli.Context) error {
				outDir := c.String("out-dir")
				fstat, err := os.Stat(outDir)
				if err != nil {
					log.Fatal(err)
				}
				if !fstat.IsDir() {
					log.Fatal("out-dir is not a directory")
				}

				for _, filename := range c.Args() {
					// Skip *.enc
					if strings.HasSuffix(filename, ".enc") {
						continue
					}
					// Skip dir
					fstat, err := os.Stat(filename)
					if err != nil {
						log.Fatal(err)
					}
					if fstat.IsDir() {
						continue
					}
					file, err := ioutil.ReadFile(filename)
					if err != nil {
						log.Fatal(err)
					}

					val, err := encrypt(
						c.String("project"),
						c.String("location"),
						c.String("keyring"),
						c.String("key"),
						file,
					)
					if err != nil {
						log.Fatal(err)
					}

					outputPath := filepath.Join(outDir, filepath.Base(filename)+".enc")
					ioutil.WriteFile(outputPath, []byte(val), fstat.Mode())

					log.Println(val)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
