package testutil

import (
	"path/filepath"
	"testing"
)

type TestEnv struct {
	DatabasePath,
	DatabaseFileTemplate	string
}

func LoadTestEnv(t *testing.T) TestEnv {
	dir := t.TempDir()
	return TestEnv{
		DatabasePath:         dir,
		DatabaseFileTemplate: filepath.Join(dir, "test_%s.db"),
	}
}
