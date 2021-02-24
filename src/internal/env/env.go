package env

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const EnginePg = "postgres"
const EngineMy = "mysql"

// Contains db credentials
type DbEnv struct {
	Engine   string
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// Return port as string
func (de *DbEnv) PortAsString() string {
	return strconv.Itoa(de.Port)
}

// Return appropriate dsn connection string to database which is used by external migration library
func (de *DbEnv) GetDsn() (string, error) {
	switch de.Engine {
	case "postgres", "postgresql":
		return de.Engine + "://" + de.User + ":" + de.Password + "@" + de.Host + ":" + de.PortAsString() + "/" + de.Database + "?sslmode=disable", nil
	case "mysql":
		return de.Engine + "://" + de.User + ":" + de.Password + "@tcp(" + de.Host + ":" + de.PortAsString() + ")/" + de.Database, nil
	default:
		return "", errors.New("Undefined database engine '" + de.Engine + "'")
	}
}

// Describe command for migration
type MigrateCmd struct {
	Cmd     string
	Version *uint
}

// Contains common info about environment
type CmdEnv struct {
	DbAlias      string
	DbDir        string
	MigrationDir string
	SnapshotDir  string
	MigrateCmd   *MigrateCmd
	dbEnv        *DbEnv
}

func (sp *CmdEnv) GetDbEnv() *DbEnv {
	if sp.dbEnv == nil {
		dbEnvPath := filepath.Join(sp.DbDir, ".env")

		dbEnv, err := godotenv.Read(dbEnvPath)
		if err != nil {
			log.Fatal(err)
		}

		port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
		sp.dbEnv = &DbEnv{
			Engine:   os.Getenv("DB_ENGINE"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			Database: os.Getenv("DB_DATABASE"),
		}

		if val, ok := dbEnv["DB_ENGINE"]; sp.dbEnv.Engine == "" && ok {
			sp.dbEnv.Engine = val
		}
		if val, ok := dbEnv["DB_USER"]; sp.dbEnv.User == "" && ok {
			sp.dbEnv.User = val
		}
		if val, ok := dbEnv["DB_PASSWORD"]; sp.dbEnv.Password == "" && ok {
			sp.dbEnv.Password = val
		}
		if val, ok := dbEnv["DB_HOST"]; sp.dbEnv.Host == "" && ok {
			sp.dbEnv.Host = val
		}
		if val, ok := dbEnv["DB_PORT"]; sp.dbEnv.Port == 0 && ok {
			val, _ := strconv.Atoi(val)
			sp.dbEnv.Port = val
		}
		if val, ok := dbEnv["DB_DATABASE"]; sp.dbEnv.Database == "" && ok {
			sp.dbEnv.Database = val
		}
	}

	return sp.dbEnv
}

// Initialize parameters and environment based on cmd line
func InitCmdEnv() CmdEnv {
	if len(os.Args) < 2 {
		log.Fatal("Database alias should be defined as first command line argument")
	}

	sp := CmdEnv{}

	sp.DbAlias = os.Args[1]
	if len(os.Args) > 2 {
		sp.MigrateCmd = &MigrateCmd{Cmd: os.Args[2]}
		if len(os.Args) > 3 {
			v64, err := strconv.ParseUint(os.Args[3], 10, 0)
			if err != nil {
				log.Fatal(err)
			}
			v := uint(v64)
			sp.MigrateCmd.Version = &v
		}
	}

	// Reading of environment
	binDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	sp.DbDir = filepath.Join(filepath.Dir(binDir), "databases", sp.DbAlias)
	sp.MigrationDir = filepath.Join(sp.DbDir, "migrations")
	sp.SnapshotDir = filepath.Join(sp.DbDir, "snapshots")

	return sp
}
