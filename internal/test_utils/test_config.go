package test_util

import (
	"os"
)

// SetupTestConfig заполняет переменные окружения для тестирования
func SetupTestConfig(env string) {
	os.Setenv("ENV", env)
	err := os.Setenv("AUTH_TOKEN", "test-token")
	if err != nil {
		return
	}
	err = os.Setenv("DB_USER", "test_user")
	if err != nil {
		return
	}
	err = os.Setenv("DB_PASSWORD", "test_password")
	if err != nil {
		return
	}
	err = os.Setenv("DB_HOST", "localhost")
	if err != nil {
		return
	}
	err = os.Setenv("DB_PORT", "5432")
	if err != nil {
		return
	}
	os.Setenv("DB_NAME", "test_db")
	err = os.Setenv("SERVER_ADDRESS", "localhost:8080")
	if err != nil {
		return
	}
	err = os.Setenv("SERVER_TIMEOUT", "5s")
	if err != nil {
		return
	}
	err = os.Setenv("MIGRATIONS_PATH", "migrations")
	if err != nil {
		return
	}
}

// CleanupTestConfig удаляет переменные окружения после тестов
func CleanupTestConfig() {
	err := os.Unsetenv("ENV")
	if err != nil {
		return
	}
	os.Unsetenv("AUTH_TOKEN")
	err = os.Unsetenv("DB_USER")
	if err != nil {
		return
	}
	err = os.Unsetenv("DB_PASSWORD")
	if err != nil {
		return
	}
	err = os.Unsetenv("DB_HOST")
	if err != nil {
		return
	}
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("SERVER_TIMEOUT")
	os.Unsetenv("MIGRATIONS_PATH")

}
