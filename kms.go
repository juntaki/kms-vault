package main

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func kmsNameFromContext(c *cli.Context) string {
	return kmsName(
		c.String("project"),
		c.String("location"),
		c.String("keyring"),
		c.String("key"),
	)
}

func kmsName(projectID, location, ringID, keyID string) string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectID, location, ringID, keyID)
}

func kmsEncrypt(name string, plaintext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("cloudkms.NewKeyManagementClient: %w", err)
	}

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      name,
		Plaintext: plaintext,
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, xerrors.Errorf("kmsEncrypt: %w", err)
	}

	return resp.Ciphertext, nil
}

func kmsDecrypt(name string, ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("cloudkms.NewKeyManagementClient: %w", err)
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       name,
		Ciphertext: ciphertext,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, xerrors.Errorf("decrypt: %w", err)
	}
	return resp.Plaintext, nil
}
