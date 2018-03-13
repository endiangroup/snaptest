package snaptest

type Snapshotter interface {
	Snapshot(object ...interface{})
}
