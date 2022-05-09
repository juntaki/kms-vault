package kms

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"github.com/urfave/cli"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type Client struct {
	client *cloudkms.KeyManagementClient
	name   string
}

func NewKMSClient(c *cli.Context) (*Client, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %w", err)
	}
	name := nameFromContext(c)
	return &Client{client: client, name: name}, nil
}

func nameFromContext(c *cli.Context) string {
	return name(
		c.String("project"),
		c.String("location"),
		c.String("keyring"),
		c.String("key"),
	)
}

func name(projectID, location, ringID, keyID string) string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectID, location, ringID, keyID)
}

func Flags(key, keyring, location, project string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "key",
			Usage:    "The key to use for encryption.",
			Value:    key,
			Required: key == "",
		},
		cli.StringFlag{
			Name:     "keyring",
			Usage:    "Key ring of the key.",
			Value:    keyring,
			Required: keyring == "",
		},
		cli.StringFlag{
			Name:     "location",
			Usage:    "Location of the keyring.",
			Value:    location,
			Required: location == "",
		},
		cli.StringFlag{
			Name:     "project",
			Usage:    "Google cloud project name.",
			Value:    project,
			Required: project == "",
		},
	}
}

func (c *Client) Encrypt(plaintext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %w", err)
	}

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      c.name,
		Plaintext: plaintext,
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Encrypt: %w", err)
	}

	return resp.Ciphertext, nil
}

func (c *Client) Decrypt(ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %w", err)
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       c.name,
		Ciphertext: ciphertext,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	return resp.Plaintext, nil
}
