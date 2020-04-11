package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

func formatEncrypted(ciphertext []byte) []byte {
	return []byte(fmt.Sprintf("$VAULT;0.1.0;CLOUD_KMS\n%s", base64.StdEncoding.EncodeToString(ciphertext)))
}

func isEncrypted(file []byte) bool {
	// TODO: Parse version, if file format changes.a
	return bytes.HasPrefix(file, []byte("$VAULT;"))
}
