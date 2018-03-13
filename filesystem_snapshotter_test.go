package snaptest

import (
	"os"
	"testing"

	"github.com/sanity-io/litter"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_AFilesystemSnapshotterCanIdentifyTheRightPath(t *testing.T) {

	snap := NewFilesystemSnapshotter(t, afero.NewMemMapFs()).(*filesystemSnapshotter)

	for _, test := range []struct {
		object interface{}
		path   string
	}{
		{
			object: "123",
			path:   ".snapshot/Test_AFilesystemSnapshotterCanIdentifyTheRightPath/b45cffe084dd3d20d928bee85e7b0f21",
		},
		{
			object: struct{}{},
			path:   ".snapshot/Test_AFilesystemSnapshotterCanIdentifyTheRightPath/1bacf627f3b49a23d8612d30db0b71dd",
		},
		{
			object: &struct{}{},
			path:   ".snapshot/Test_AFilesystemSnapshotterCanIdentifyTheRightPath/4a24a9e5477846dfa9a68f345f91d382",
		},
		{
			object: t,
			path:   ".snapshot/Test_AFilesystemSnapshotterCanIdentifyTheRightPath/9201dafe08a33bbb90680a051adde096",
		},
	} {
		require.Equal(t, test.path, snap.pathForObject(test.object))
	}

	t.Run("some subtest", func(t2 *testing.T) {
		snap2 := NewFilesystemSnapshotter(t2, afero.NewMemMapFs()).(*filesystemSnapshotter)
		require.Equal(t2, ".snapshot/Test_AFilesystemSnapshotterCanIdentifyTheRightPath/some_subtest/b45cffe084dd3d20d928bee85e7b0f21", snap2.pathForObject("123"))
	})
}

func Test_AFilesystemSnapshotterCanTestSnapshots(t *testing.T) {

	fs := afero.NewMemMapFs()

	snap := NewFilesystemSnapshotter(t, fs).(*filesystemSnapshotter)
	object := struct {
		Name string
	}{
		Name: "123",
	}

	// Write expected data to expected file
	afero.WriteFile(fs, snap.pathForObject(object), []byte(litter.Sdump(object)), os.ModePerm)

	// ... And test
	snap.Snapshot(object)
}
