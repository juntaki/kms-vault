package main

import (
	"log"
	"os"

	"github.com/juntaki/kms-vault/kms"

	"github.com/urfave/cli"
)

const VaultVersion = "0.1.0"

type Crypter interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

var kmsClient Crypter

func initializeKMS(c *cli.Context) (err error) {
	kmsClient, err = kms.NewKMSClient(c)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	config := loadConfig()
	kmsFlags := kms.Flags(
		config.Key,
		config.KeyRing,
		config.Location,
		config.Project,
	)

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
