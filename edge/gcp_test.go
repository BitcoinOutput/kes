// Copyright 2023 - MinIO, Inc. All rights reserved.
// Use of this source code is governed by the AGPLv3
// license that can be found in the LICENSE file.

package edge_test

import (
	"flag"
	"os"
	"testing"

	"github.com/minio/kes/edge"
)

var gcpConfigFile = flag.String("gcp.config", "", "Path to a KES config file with GCP SecretManager config")

func TestGCP(t *testing.T) {
	if *gcpConfigFile == "" {
		t.Skip("GCP tests disabled. Use -gcp.config=<FILE> to enable them")
	}
	file, err := os.Open(*gcpConfigFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	config, err := edge.ReadServerConfigYAML(file)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := config.KeyStore.(*edge.GCPSecretManagerKeyStore); !ok {
		t.Fatalf("Invalid Keystore: want %T - got %T", config.KeyStore, &edge.GCPSecretManagerKeyStore{})
	}

	ctx, cancel := testingContext(t)
	defer cancel()

	store, err := config.KeyStore.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Create", func(t *testing.T) { testCreate(ctx, store, t) })
	t.Run("Set", func(t *testing.T) { testSet(ctx, store, t) })
	t.Run("Get", func(t *testing.T) { testGet(ctx, store, t) })
	t.Run("Status", func(t *testing.T) { testStatus(ctx, store, t) })
}
