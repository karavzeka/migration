package migration

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "testfile.txt")
	if err := createFile(file); err != nil {
		t.Error(err)
	}
	assert.FileExists(t, file)
}

func TestLastSeqVersion(t *testing.T) {
	var err error

	dir := filepath.Join(t.TempDir(), "tmpDir")
	if err = createTmpFiles(dir, "000065_test.up.sql", "000066_test.up.sql"); err != nil {
		t.Error(err)
	}

	// secondly, get las version
	v, err := LastSeqVersion(dir)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "000066", v)
}

func TestNextSeqVersion(t *testing.T) {
	var err error

	dir := filepath.Join(t.TempDir(), "tmpDir")
	if err = createTmpFiles(dir, "000065_test.up.sql", "000066_test.up.sql"); err != nil {
		t.Error(err)
	}

	// secondly, get las version
	v, err := nextSeqVersion(dir)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "000067", v)
}

func TestNewUpDown(t *testing.T) {
	var upPath, downPath string
	var err error

	dir := filepath.Join(t.TempDir(), "migrations")

	// first call
	upPath, downPath, err = NewUpDown(dir, "test")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, filepath.Join(dir, "000001_test.up.sql"), upPath)
	assert.Equal(t, filepath.Join(dir, "000001_test.down.sql"), downPath)
	assert.FileExists(t, upPath)
	assert.FileExists(t, downPath)

	// second call
	upPath, downPath, err = NewUpDown(dir, "test")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, filepath.Join(dir, "000002_test.up.sql"), upPath)
	assert.Equal(t, filepath.Join(dir, "000002_test.down.sql"), downPath)
	assert.FileExists(t, upPath)
	assert.FileExists(t, downPath)
}

func TestNewSnapshot(t *testing.T) {
	dirSnap := filepath.Join(t.TempDir(), "snapshots")
	snapshotPath, _ := NewSnapshot(dirSnap, 66)

	assert.Equal(t, filepath.Join(dirSnap, "000066.sql"), snapshotPath)
}

func TestNearestSnapshot(t *testing.T) {
	var err error
	var snapshotPath string
	var snapshotVersion uint

	dir := filepath.Join(t.TempDir(), "tmpDir")
	if err = createTmpFiles(dir, "000002.sql", "000004.sql"); err != nil {
		t.Error(err)
	}

	_, _, err = NearestSnapshot(dir, 1)
	assert.ErrorIs(t, err, ErrSnapshotNotExist)

	snapshotPath, snapshotVersion, err = NearestSnapshot(dir, 2)
	assert.Equal(t, filepath.Join(dir, "000002.sql"), snapshotPath)
	assert.Equal(t, uint(2), snapshotVersion)

	snapshotPath, snapshotVersion, err = NearestSnapshot(dir, 3)
	assert.Equal(t, uint(2), snapshotVersion)
}

func createTmpFiles(dir string, fileNames ...string) error {
	var err error
	if err = os.RemoveAll(dir); err != nil {
		return err
	}
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, fileName := range fileNames {
		f, err := os.Create(filepath.Join(dir, fileName))
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()
	}

	return nil
}
