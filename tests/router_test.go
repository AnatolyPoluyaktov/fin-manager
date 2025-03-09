package tests

import (
	"fin-manager/internal/config"
	test_util "fin-manager/internal/test_utils"
	"log"
	"os"
	"testing"
)

var testDB *test_util.TestDb

const TestToken = "my-secret-token"

func TestMain(m *testing.M) {
	var err error
	test_util.SetupTestConfig("test")
	cfg := config.MustLoadConfig()
	testDB, err = test_util.SetupTestDb(*cfg)
	if err != nil {
		log.Fatalf("error setting up test db: %v", err)
	}
	code := m.Run()
	testDB.CleanUp()
	test_util.CleanupTestConfig()
	os.Exit(code)
}
