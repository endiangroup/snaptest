package snaptest

import (
	"testing"

	"github.com/spf13/afero"
)

// A utility function that assumes the default path and a local filesystem
func Snapshot(t *testing.T, objects ...interface{}) {
	snap := NewFilesystemSnapshotter(t, afero.NewOsFs())
	snap.Snapshot(objects...)
}
