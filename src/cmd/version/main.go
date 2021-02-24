package main

import (
	"fmt"
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/env"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

const forceParam = "force"

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

	if ce.MigrateCmd != nil && ce.MigrateCmd.Cmd == forceParam {
		if ce.MigrateCmd.Version != nil {
			v = *ce.MigrateCmd.Version
		}
		if err = m.Force(int(v)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database forced to version ", v)
	}

	state := "clean"
	if dirty {
		state = "dirty"
	}

	fmt.Println(v, state)
}
