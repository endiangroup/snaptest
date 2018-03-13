package snaptest

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sanity-io/litter"
	"github.com/spf13/afero"
)

const defaultPath = ".snapshot"

type filesystemSnapshotter struct {
	t    *testing.T
	fs   afero.Fs
	path string
}

func NewFilesystemSnapshotter(t *testing.T, fs afero.Fs, paths ...string) Snapshotter {

	path := defaultPath

	if len(paths) > 0 {
		path = paths[0]
	}

	return &filesystemSnapshotter{
		t:    t,
		fs:   fs,
		path: path,
	}
}

func (f *filesystemSnapshotter) Snapshot(objects ...interface{}) {

	abort := false

	for _, object := range objects {
		path := f.pathForObject(object)
		code := litter.Sdump(object)

		data, err := afero.ReadFile(f.fs, path)

		// Write the file if it doesn't exist
		if err != nil {
			f.t.Logf("Snapshot file %s doesn't exist; creating", path)
			abort = true

			if err := f.fs.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				f.t.Fatalf("Failed to create dir %s: %s", filepath.Dir(path), err)
			}

			if err := afero.WriteFile(f.fs, path, []byte(code), os.ModePerm); err != nil {
				f.t.Fatalf("Failed to write %s: %s", path, err)
			}

			continue
		}

		name := reflect.TypeOf(object).String()

		f.t.Logf("Snapshot of %s:\n\n%s", name, code)

		if string(data) != code {
			f.t.Fatalf("Snapshot comparison failed for %s\n\nExpected:\n\n%s\n\nGot:\n\n%s\n\n%s", name, data, code, path)
		}
	}

	if abort {
		f.t.Fatalf("Missing snapshots; run go test again")
	}
}

func (f *filesystemSnapshotter) pathForObject(object interface{}) string {

	typ := reflect.TypeOf(object)

	name := fmt.Sprintf("%x", md5.Sum([]byte(typ.String())))
	return filepath.Join(f.path, f.t.Name(), name)
}
