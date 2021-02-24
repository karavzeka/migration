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
	"log"
	"strconv"
)

func main() {
	ce := env.InitCmdEnv()

	dsn, err := ce.GetDbEnv().GetDsn()
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New("file://"+ce.MigrationDir, dsn)
	if err != nil {
		log.Fatal(err)
	}

	v, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}
	if dirty {
		log.Fatal("Current version " + strconv.Itoa(int(v)) + " is dirty. Please fix current version manually before migration")
	}

	snapshotPath, err := migration.NewSnapshot(ce.SnapshotDir, v)
	if err != nil {
		log.Fatal(err)
	}

	err = fillSnapshotFile(ce.GetDbEnv(), snapshotPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Snapshot created:\n%s\n", snapshotPath)
}

func fillSnapshotFile(de *env.DbEnv, snapshotPath string) error {
	var err error
	switch de.Engine {
	case env.EnginePg:
		err = db_utils.DumpPgsql(de, snapshotPath)
	case env.EngineMy:
		err = db_utils.DumpMysql(de, snapshotPath)
	default:
		err = errors.New("Db engine " + de.Engine + " is not defined")
	}

	return err
}
