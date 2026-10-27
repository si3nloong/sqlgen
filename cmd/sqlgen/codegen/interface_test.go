package codegen

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Package initialization must work even before the caller has a valid module.
func TestInterfaceInitOutsideModule(t *testing.T) {
	if os.Getenv("SQLGEN_TEST_INTERFACE_INIT") == "1" {
		if typeOfTime != "time.Time" || goSqlValuer == nil || goSqlScanner == nil ||
			binaryMarshaler == nil || binaryUnmarshaler == nil ||
			textMarshaler == nil || textUnmarshaler == nil {
			t.Fatal("standard-library types were not initialized")
		}
		return
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("invalid module file\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestInterfaceInitOutsideModule$")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "SQLGEN_TEST_INTERFACE_INIT=1", "GOWORK=off")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("initialization in an invalid module failed: %v\n%s", err, output)
	}
}
