package migration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const ext = "sql"
const seqDigits = 6

var ErrSnapshotNotExist = errors.New("snapshot does not exist")

// Creates up/down empty files with new version prefix and specified name in the specified directory.
// If 'dry' is set as true, new file paths will be returned without creating real files.
func NewUpDown(dir string, name string) (upPath, downPath string, err error) {
	var version string
	var matches []string

	version, err = nextSeqVersion(dir)
	if err != nil {
		return "", "", err
	}

	versionGlob := filepath.Join(dir, version+"_*."+ext)
	matches, err = filepath.Glob(versionGlob)
	if err != nil {
		return "", "", err
	}

	if len(matches) > 0 {
		return "", "", fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", "", err
	}

	newPaths := [2]string{"up", "down"}
	for i, direction := range newPaths {
		basename := fmt.Sprintf("%s_%s.%s.%s", version, name, direction, ext)
		filePath := filepath.Join(dir, basename)

		if err = createFile(filePath); err != nil {
			return "", "", err
		}

		newPaths[i] = filePath
	}

	return newPaths[0], newPaths[1], nil
}

// Creates new snapshot file with specified version
func NewSnapshot(snapshotDir string, version uint) (path string, err error) {
	var matches []string

	versionStr := fmt.Sprintf("%0[2]*[1]d", version, seqDigits)
	snapshotPath := filepath.Join(snapshotDir, versionStr+"."+ext)
	matches, err = filepath.Glob(snapshotPath)
	if err != nil {
		return "", err
	}

	if len(matches) > 0 {
		return "", fmt.Errorf("duplicate snapshot version: %s", versionStr)
	}

	if err = os.MkdirAll(snapshotDir, os.ModePerm); err != nil {
		return "", err
	}

	if err = createFile(snapshotPath); err != nil {
		return "", err
	}

	return snapshotPath, nil
}

// Returns nearest snapshot path lower or equal than specified version
func NearestSnapshot(snapshotDir string, version uint) (path string, snapshotVersion uint, err error) {
	var lastVer uint

	files, err := ioutil.ReadDir(snapshotDir)
	if err != nil {
		return "", 0, ErrSnapshotNotExist
	}

	for _, file := range files {
		if len(file.Name()) < seqDigits {
			return "", 0, fmt.Errorf("malformed snapshot filename: %s/%s", snapshotDir, file.Name())
		}

		matchSeqStr := file.Name()[0:seqDigits]
		seq, err := strconv.ParseUint(matchSeqStr, 10, 64)
		if err != nil {
			return "", 0, err
		}
		if uint(seq) > version {
			break
		}

		lastVer = uint(seq)
		path = filepath.Join(snapshotDir, file.Name())
	}

	if lastVer == 0 {
		return "", 0, ErrSnapshotNotExist
	} else {
		return path, lastVer, nil
	}
}

// Analise migration files extracting last version
func LastSeqVersion(migrationDir string) (string, error) {
	var lastSeq uint64

	migrationDir = filepath.Clean(migrationDir)
	matches, _ := filepath.Glob(filepath.Join(migrationDir, "*."+ext))

	if len(matches) > 0 {
		filePath := matches[len(matches)-1]
		filename := filepath.Base(filePath)

		if len(filename) < seqDigits {
			return "", fmt.Errorf("malformed migration filename: %s", filePath)
		}

		var err error
		matchSeqStr := filename[0:seqDigits]
		lastSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)
		if err != nil {
			return "", err
		}
	}

	version := fmt.Sprintf("%0[2]*[1]d", lastSeq, seqDigits)
	return version, nil
}

// Gets list of paths, find next version based the last file prefix
func nextSeqVersion(migrationDir string) (string, error) {
	lastVersion, err := LastSeqVersion(migrationDir)
	if err != nil {
		return "", err
	}

	lastVersionUint, err := strconv.ParseUint(lastVersion, 10, 64)
	if err != nil {
		return "", err
	}

	lastVersionUint++

	version := fmt.Sprintf("%0[2]*[1]d", lastVersionUint, seqDigits)
	if len(version) > seqDigits {
		return "", fmt.Errorf("next sequence number %s too large, at most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}

func createFile(filename string) error {
	// create exclusive (fails if file already exists)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0664)

	if err != nil {
		return err
	}

	return f.Close()
}
