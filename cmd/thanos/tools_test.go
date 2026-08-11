// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package main

import (
	"os"
	"testing"

	"github.com/go-kit/log"

	"github.com/efficientgo/core/testutil"
)

func Test_CheckRules(t *testing.T) {
	validFiles := []string{
		"./testdata/rules-files/valid.yaml",
	}

	invalidFiles := [][]string{
		{"./testdata/rules-files/non-existing-file.yaml"},
		{"./testdata/rules-files/invalid-yaml-format.yaml"},
		{"./testdata/rules-files/invalid-rules-data.yaml"},
		{"./testdata/rules-files/invalid-unknown-field.yaml"},
	}

	logger := log.NewNopLogger()
	testutil.Ok(t, checkRulesFiles(logger, &validFiles))

	for _, fn := range invalidFiles {
		testutil.NotOk(t, checkRulesFiles(logger, &fn), "expected err for file %s", fn)
	}
}

func Test_CheckRules_Glob(t *testing.T) {
	// regex path
	files := &[]string{"./testdata/rules-files/valid*.yaml"}
	logger := log.NewNopLogger()
	testutil.Ok(t, checkRulesFiles(logger, files))

	// direct path
	files = &[]string{"./testdata/rules-files/valid.yaml"}
	testutil.Ok(t, checkRulesFiles(logger, files))

	// invalid path
	files = &[]string{"./testdata/rules-files/*.yamlaaa"}
	testutil.NotOk(t, checkRulesFiles(logger, files), "expected err for file %s", files)

	// Unreadable path — skip when chmod is not effective (root bypasses permissions,
	// non-owner in container can't chmod, read-only FS, etc.)
	files = &[]string{"./testdata/rules-files/unreadable_valid.yaml"}
	filename := (*files)[0]
	if err := os.Chmod(filename, 0000); err != nil {
		t.Skipf("skipping unreadable-file test: chmod failed: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(filename, 0777) })
	if os.Getuid() == 0 {
		t.Skip("skipping unreadable-file test: root bypasses file permissions")
	}
	testutil.NotOk(t, checkRulesFiles(logger, files), "expected err for file %s", files)
	testutil.Ok(t, os.Chmod(filename, 0777), "failed to restore file permissions of %s", filename)
}
