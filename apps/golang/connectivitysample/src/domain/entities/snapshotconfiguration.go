package entities

type SnapshotConfiguration interface {
	SnapshotMaxWidth() int
	SnapshotMaxHeight() int
}

type snapshotConfiguration struct {
	snapshotMaxWidth  int // Max Width used when fetching for camera's snapshots.
	snapshotMaxHeight int // Max Height used when fetching for camera's snapshots.
}

func NewSnapshotConfiguration(snapshotMaxWidth int, snapshotMaxHeight int) SnapshotConfiguration {
	return &snapshotConfiguration{
		snapshotMaxWidth:  snapshotMaxWidth,
		snapshotMaxHeight: snapshotMaxHeight,
	}
}

func (s *snapshotConfiguration) SnapshotMaxWidth() int {
	return s.snapshotMaxWidth
}

func (s *snapshotConfiguration) SnapshotMaxHeight() int {
	return s.snapshotMaxHeight
}
