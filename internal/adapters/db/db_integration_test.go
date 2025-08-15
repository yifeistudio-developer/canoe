package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/yifeistudio-developer/canoe/internal/ports"
)

type DatabaseSuite struct {
	suite.Suite
	dbPort ports.DbPort
}

func (s *DatabaseSuite) SetupSuite() {
	ctx := context.Background()
	port := "5432/tcp"
	dbUrl := func(host string, p nat.Port) string {
		return fmt.Sprintf("postgres://canoe:canoe@%s:%s/canoe?sslmode=disable&TimeZone=Asia/Shanghai", host, p.Port())
	}
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "canoe",
			"POSTGRES_DB":       "canoe",
			"POSTGRES_USER":     "canoe",
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "postgres", dbUrl).WithStartupTimeout(time.Second * 30),
	}
	testContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal("failed to start test container", err)
	}
	endpoint, err := testContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatal("failed to connect to test container", err)
	}
	datasourceUrl := fmt.Sprintf("postgres://canoe:canoe@%s/canoe?sslmode=disable&TimeZone=Asia/Shanghai", endpoint)
	err = os.Setenv("IRIS_APPLICATION_NAME", "canoe")
	err = os.Setenv("DATA_SOURCE_URL", datasourceUrl)
	adapter, err := NewAdapter()
	s.dbPort = adapter
}

func (s *DatabaseSuite) Test_Save_User() {
	userDbPort := s.dbPort.UserDbPort()
	user, err := userDbPort.GetById(1)
	s.Nil(err)
	log.Println(user)
}

func (s *DatabaseSuite) TearDownSuite() {}

func TestIntegrationDB(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}
