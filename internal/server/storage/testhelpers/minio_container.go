package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MinioContainer struct {
	Container testcontainers.Container
	Endpoint  string
	AccessKey string
	SecretKey string
}

func CreateMinioContainer(ctx context.Context) (*MinioContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Cmd:          []string{"server", "/data", "--console-address", ":9001"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
		},
		WaitingFor: wait.ForLog("API: http://"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "9000")
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s:%s", host, port.Port())

	// Ждем немного, чтобы MinIO полностью запустился
	time.Sleep(2 * time.Second)

	return &MinioContainer{
		Container: container,
		Endpoint:  endpoint,
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
	}, nil
}

func (m *MinioContainer) Terminate(ctx context.Context) error {
	return m.Container.Terminate(ctx)
}
