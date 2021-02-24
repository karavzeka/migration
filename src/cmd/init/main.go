package main

import (
	"errors"
	"fmt"
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/db_utils"
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/env"
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/migration"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"io/ioutil"
	"log"
	"os"
)

const dryRun = "dry"

func main() {
	ce := env.InitCmdEnv()

	files, err := ioutil.ReadDir(ce.MigrationDir)
	if err == nil && len(files) > 0 {
		log.Fatal("Database '" + ce.DbAlias + "' already has migrations, can't make init migration.")
	}

	upPath, downPath, err := migration.NewUpDown(ce.MigrationDir, "init")
	if err != nil {
		log.Fatal(err)
	}

	err = fillMigrationFiles(ce.GetDbEnv(), upPath, downPath)
	if err != nil {
		if err := os.Remove(upPath); err != nil {
			log.Println(err)
		}
		if err := os.Remove(downPath); err != nil {
			log.Println(err)
		}
		log.Fatal(err)
	}

	// Set version 1 if not dry run
	err = setFirstVersion(&ce)
	if err != nil {
		log.Fatal(err)
	}
}

// Fill 'up' file by DB dump and 'down' file by table truncates
func fillMigrationFiles(de *env.DbEnv, upPath, downPath string) error {
	var err error
	switch de.Engine {
	case env.EnginePg:
		err = db_utils.DumpPgsql(de, upPath)
		if err != nil {
			return err
		}
		err = db_utils.DropTables(upPath, downPath, env.EnginePg)
	case env.EngineMy:
		err = db_utils.DumpMysql(de, upPath)
		if err != nil {
			return err
		}
		err = db_utils.DropTables(upPath, downPath, env.EngineMy)
	default:
		err = errors.New("Db engine " + de.Engine + " is not defined")
	}

	if err != nil {
		return err
	}

	fmt.Printf("First migration initialized\n%s\n%s\n", upPath, downPath)

	return nil
}

// Set version 1 if not dry run
func setFirstVersion(ce *env.CmdEnv) error {
	if ce.MigrateCmd != nil {
		if ce.MigrateCmd.Cmd != dryRun {
			return errors.New(fmt.Sprintf("Command %q is not valid", ce.MigrateCmd.Cmd))
		}
		return nil
	}

	dsn, err := ce.GetDbEnv().GetDsn()
	if err != nil {
		return err
	}

	m, err := migrate.New("file://"+ce.MigrationDir, dsn)
	if err != nil {
		return err
	}

	err = m.Force(1)
	if err != nil {
		return err
	}

	fmt.Println("Database is set to version 1")

	return nil
}
