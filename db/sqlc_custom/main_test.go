package sqlc_custom

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	sqlc "github.com/stuneak/simplebank/db/sqlc"
	"github.com/stuneak/simplebank/util"
)

var testQueries *sqlc.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../../")

	if err != nil {
		log.Fatal("cannot load config", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	testQueries = sqlc.New(testDB)

	os.Exit(m.Run())
}
