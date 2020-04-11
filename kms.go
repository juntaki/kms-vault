package main

import (
	cloudkms "cloud.google.com/go/kms/apiv1"
	"context"
	"fmt"
	"golang.org/x/xerrors"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

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
		Plaintext: []byte(plaintext),
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, xerrors.Errorf("kmsEncrypt: %w", err)
	}

	return resp.Ciphertext, nil
}
