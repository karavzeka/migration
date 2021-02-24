package main

import (
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

	if ce.MigrateCmd == nil {
		log.Fatal("Migration command is not specified")
	}

	dsn, err := ce.GetDbEnv().GetDsn()
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New("file://"+ce.MigrationDir, dsn)
	if err != nil {
		log.Fatal(err)
	}

	m.Log = &migration.Log{}

	v, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}
	if dirty {
		log.Fatal("Dirty database version " + strconv.Itoa(int(v)) + ". Fix and force version")
	}

	switch ce.MigrateCmd.Cmd {
	case "goto":
		if ce.MigrateCmd.Version == nil {
			log.Fatal("Migration version for goto is not specified")
		}

		if err == migrate.ErrNilVersion {
			err := applyImport(&ce, m, *ce.MigrateCmd.Version)
			if err != nil && err != migration.ErrSnapshotNotExist {
				log.Fatal(err)
			}
		}

		err = m.Migrate(*ce.MigrateCmd.Version)
	case "up":
		if err == migrate.ErrNilVersion {
			var targetVersion uint
			if ce.MigrateCmd.Version != nil {
				targetVersion = *ce.MigrateCmd.Version
			} else {
				lastSeqSrt, err := migration.LastSeqVersion(ce.MigrationDir)
				if err != nil {
					log.Fatal(err)
				}
				lastSeq, _ := strconv.ParseUint(lastSeqSrt, 10, 64)
				targetVersion = uint(lastSeq)
			}
			err := applyImport(&ce, m, targetVersion)
			if err != nil && err != migration.ErrSnapshotNotExist {
				log.Fatal(err)
			}
		}

		if ce.MigrateCmd.Version != nil && *ce.MigrateCmd.Version > v {
			err = m.Migrate(*ce.MigrateCmd.Version)
		} else {
			err = m.Up()
		}
	case "down":
		if ce.MigrateCmd.Version != nil && *ce.MigrateCmd.Version < v {
			err = m.Migrate(*ce.MigrateCmd.Version)
		} else {
			err = m.Down()
		}
	default:
		log.Fatal("Undefined migration command '" + ce.MigrateCmd.Cmd + "'")
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println(err)
		} else {
			log.Fatal(err)
		}
	}
}

func applyImport(ce *env.CmdEnv, m *migrate.Migrate, version uint) error {
	var err error
	snapshotPath, snapshotVersion, err := migration.NearestSnapshot(ce.SnapshotDir, version)
	if err == migration.ErrSnapshotNotExist {
		return nil
	} else if err != nil {
		return err
	}

	de := ce.GetDbEnv()
	if de.Engine == env.EnginePg {
		err = db_utils.ImportPgsql(de, snapshotPath)
	} else {
		err = db_utils.ImportMysql(de, snapshotPath)
	}
	if err != nil {
		return err
	}

	fmt.Println("Snapshot is applied: " + snapshotPath)

	// Force version according to snapshot
	if snapshotVersion > 0 {
		if err = m.Force(int(snapshotVersion)); err != nil {
			log.Fatal(err)
		}
	}

	return err
}
