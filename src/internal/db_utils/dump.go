package db_utils

import (
	"bufio"
	"bytes"
	"errors"
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/env"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

const migrationTable = "schema_migrations"

// Call 'pg_dump' command and put result to specified file
func DumpPgsql(de *env.DbEnv, filePath string) (err error) {
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
		"pg_dump",
		"-h", de.Host,
		"-p", de.PortAsString(),
		"-U", de.User,
		"-d", de.Database,
		"-T", migrationTable,
		"-s",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out) + "\n" + err.Error())
	}

	// Manual filtering output and writing to the file
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, os.FileMode(644))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "SELECT pg_catalog.set_config('search_path', '', false);" {
			continue
		}
		_, err = f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Call 'mysqldump' command and put result to specified file
func DumpMysql(de *env.DbEnv, filePath string) (err error) {
	passFilePath, existed, err := createPasswordFileMysql(de)
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
		"mysqldump",
		"-h", de.Host,
		"-P", de.PortAsString(),
		"-u", de.User,
		"--ignore-table="+de.Database+"."+migrationTable,
		"--no-tablespaces",
		"--no-data",
		"--skip-add-drop-table",
		de.Database,
	)

	outfile, err := os.OpenFile(filePath, os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer func() {
		err = outfile.Close()
	}()

	cmd.Stdout = outfile
	err = cmd.Run()

	return err
}

// Writes DROP TABLE statements in file according to CREATE TABLE patterns in dump file
func DropTables(dumpFile, dropTableFile string, dbEngine string) (err error) {
	tf, err := os.OpenFile(dropTableFile, os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer func() {
		err = tf.Close()
	}()

	if dbEngine == env.EnginePg {
		_, err = tf.WriteString("START TRANSACTION;\n")
		if err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(dumpFile)
	if err != nil {
		return err
	}

	r, err := regexp.Compile("(?i)create table (if not exists )?`?([\\w.]+)`?\\s*\\(")
	if err != nil {
		return err
	}

	matches := r.FindAllStringSubmatch(string(b), -1)
	for _, match := range matches {
		if len(match) >= 2 {
			tableName := match[2]
			if dbEngine == env.EngineMy {
				tableName = "`" + tableName + "`"
			}
			_, err = tf.WriteString("DROP TABLE " + tableName + ";\n")
			if err != nil {
				return err
			}
		}
	}

	if dbEngine == env.EnginePg {
		_, err = tf.WriteString("COMMIT;\n")
		if err != nil {
			return err
		}
	}

	return err
}
