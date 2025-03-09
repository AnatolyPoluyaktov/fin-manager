package test_util

import (
	"context"
	"fin-manager/internal/config"
	"fin-manager/internal/storage/pg"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type TestDb struct {
	Container testcontainers.Container
	DB        *pg.DB
}

func SetupTestDb(cfg config.Config) (*TestDb, error) {
	ctx := context.Background()
	request := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.Storage.DbUser,
			"POSTGRES_PASSWORD": cfg.Storage.DbPassword,
			"POSTGRES_DB":       cfg.Storage.DbName,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}
	host, err := container.Host(ctx)

	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}
	time.Sleep(2 * time.Second)

	cfg.Storage.DbPort = port.Port()
	cfg.Storage.DbHost = host
	db := pg.New(ctx, cfg)
	err = db.Migrate(cfg.MigrationsPath)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}
	return &TestDb{
		Container: container,
		DB:        db,
	}, nil
}

func (t *TestDb) CleanUp() {
	t.DB.Close()
	t.Container.Terminate(context.Background())
}
