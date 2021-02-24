package main

import (
	"flag"
	"fmt"
	"git.a.kluatr.ru/e.karavskii/migration-example/internal/migration"
	"log"
	"os"
)

func main() {
	dir := flag.String("dir", "", "directory for migration files")
	name := flag.String("name", "noname", "name of migration")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help || *dir == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	upPath, downPath, err := migration.NewUpDown(*dir, *name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Migration files are created:\n%s\n%s\n", upPath, downPath)
}
