package db_utils

import (
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/env"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Creates pgsql credentials file
func createPasswordFilePgsql(de *env.DbEnv) (path string, existed bool, err error) {
	pgConnString := strings.Join([]string{de.Host, de.PortAsString(), de.Database, de.User, de.Password}, ":")
	path, existed, err = createPasswordFile(".pgpass", pgConnString)
	return
}

// Creates mysql credentials file
func createPasswordFileMysql(de *env.DbEnv) (path string, existed bool, err error) {
	myCongString := "[client]\npassword=" + de.Password
	path, existed, err = createPasswordFile(".my.cnf", myCongString)
	return
}

// Creates file to keep database connection password
// It returns full file path and whether the file existed
func createPasswordFile(fileName string, content string) (path string, existed bool, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	path = filepath.Join(usr.HomeDir, fileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return "", false, err
		}
		if _, err := file.WriteString(content); err != nil {
			return "", false, err
		}
		if err = file.Close(); err != nil {
			return "", false, err
		}
	}

	return path, !os.IsNotExist(err), nil
}
