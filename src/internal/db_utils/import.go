package db_utils

import (
	"errors"
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/env"
	"os"
	"os/exec"
)

// Imports database from dump file (postgresql)
func ImportPgsql(de *env.DbEnv, snapshotPath string) (err error) {
	passFilePath, existed, err := createPasswordFilePgsql(de)
	if err != nil {
		return
	}
	if !existed {
		// if file created by script, remove it after using
		defer func() {
			err = os.Remove(passFilePath)
		}()
	}

	cmd := exec.Command(
		"psql",
		"-h", de.Host,
		"-p", de.PortAsString(),
		"-U", de.User,
		"-f", snapshotPath,
		"-d", de.Database,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out) + "\n" + err.Error())
	}

	return nil
}

// Imports database from dump file (mysql)
func ImportMysql(de *env.DbEnv, snapshotPath string) (err error) {
	cmd := exec.Command(
		"mysqlimport",
		"-h", de.Host,
		"-P", de.PortAsString(),
		"-u", de.User,
		de.Database,
		snapshotPath,
	)
	err = cmd.Run()

	return err
}
