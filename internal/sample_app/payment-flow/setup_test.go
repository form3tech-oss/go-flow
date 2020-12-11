package payment_flow

import (
	"database/sql"
	"fmt"
	api2 "github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/settings"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/form3tech/go-docker-compose/dockercompose"
	"github.com/form3tech/go-logger/log"
	"github.com/form3tech/go-pact-testing/pacttesting"
	"github.com/giantswarm/retry-go"
	"github.com/jmoiron/sqlx"
	"github.com/phayes/freeport"
	"github.com/spf13/viper"
)

var ServerPort = settings.ServerPort
var database *sqlx.DB

func getServerPort() int {
	var serverPort int
	if _, err := os.Stat(".serverport"); os.IsNotExist(err) {
		serverPort, _ = freeport.GetFreePort()
		_ = ioutil.WriteFile(".serverport", []byte(fmt.Sprintf("%d", serverPort)), 0644)
	} else {
		portStr, _ := ioutil.ReadFile(".serverport")
		serverPort, _ = strconv.Atoi(string(portStr))
	}
	return serverPort
}

func TestMain(m *testing.M) {
	_ = os.Setenv("STACK_NAME", "local")

	dir, _ := os.Getwd()

	dynamicPorts, err := dockercompose.NewDynamicPorts(
		"POSTGRES_PORT:postgresql:5432",
	)

	if err != nil {
		panic(err)
	}

	ServerPort = getServerPort()
	settings.ServerPort = ServerPort

	var localAddress string
	if runtime.GOOS == "darwin" {
		localAddress = "docker.for.mac.host.internal"
	} else if os.Getenv("HOST_IP") != "" {
		localAddress = os.Getenv("HOST_IP")
	} else {
		localAddress = "localhost"
	}

	_ = os.Setenv("SQS_HOST", localAddress)
	_ = os.Setenv("SNS_HOST", localAddress)

	serviceNameNoHyphens := strings.Replace(settings.ServiceName, "-", "", -1)
	dc, err := dockercompose.NewDockerCompose(
		dockercompose.NewAwsEcrAuth("288840537196", "eu-west-1"),
		filepath.Join(dir, "dockercompose/docker-compose.yml"),
		serviceNameNoHyphens+"testing",
		dynamicPorts,
		"indoor_terraform")

	if err != nil {
		panic(err)
	}

	_ = os.Setenv("VAULT_ADDR", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("VAULT_PORT")))
	_ = os.Setenv("AWS_REGION", "eu-west-1")
	_ = os.Setenv("LOG_LEVEL", "debug")

	databaseUrl := fmt.Sprintf("postgres://postgres:password@localhost:%d?sslmode=disable", dc.GetDynamicContainerPort("POSTGRES_PORT"))
	viper.Set("DATABASE_HOST", "localhost")
	viper.Set("DATABASE_PORT", dc.GetDynamicContainerPort("POSTGRES_PORT"))
	viper.Set("database-ssl-mode", "disable")

	viper.Set("MessageVisibilityTimeout", 5)
	viper.Set("aws.default.sqs.endpoint", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("SQS_PORT")))
	viper.Set("aws.default.sns.endpoint", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("SNS_PORT")))
	viper.Set("aws.default.sts.endpoint", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("STS_PORT")))

	containerWaiter := dockercompose.
		WaitForContainersToStartWithTimeout(10*time.Minute).
		ContainerLogLine("postgresql_1", "database system is ready to accept connections").
		Container("postgresql_1", func(finish chan error) {
			err := retry.Do(func() error {

				_, err := sql.Open("postgres", databaseUrl)
				return err
			}, retry.MaxTries(500), retry.Sleep(time.Duration(500*time.Millisecond)))
			finish <- err
		})

	err = dc.Start(containerWaiter)

	if err != nil {
		panic(err)
	}

	viper.Set(settings.ServiceName+"-address", fmt.Sprintf("http://localhost:%d", ServerPort))

	api2.Configure()

	startedSignal := make(chan bool)
	stopServer := make(chan os.Signal)
	stopped := make(chan struct{})
	signal.Notify(stopServer, os.Interrupt)

	go func() { api2.Start(stopServer, startedSignal, stopped) }()
	<-startedSignal

	err = WaitForPort(ServerPort, time.Second*10)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
	log.Infof("Server port available (%d).", ServerPort)

	database = api2.ConnectToDatabase()
	result := m.Run()

	stopServer <- syscall.SIGKILL
	<-stopped

	dc.Stop()

	pacttesting.StopMockServers()

	os.Exit(result)
}

func WaitForPort(port int, timeout time.Duration) error {
	delay := time.Millisecond * 50
	var attempts uint

	if timeout < delay {
		attempts = 1
	} else {
		attempts = uint((timeout / delay) + 1)
	}

	if err := retry.Do(
		func() error {
			log.Info("Checking for server port.")
			conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return err
			}
			_ = conn.Close()
			return nil
		},
		retry.MaxTries(int(attempts)),
		retry.Sleep(timeout),
	); err != nil {
		return fmt.Errorf("timed out waiting for server to be ready at :%d", port)
	}

	return nil
}
func truncateAllTables(t *testing.T) {
	tables := []string{"Payment"}
	for _, v := range tables {
		truncateTable(v, t)
	}
}

func truncateTable(table string, t *testing.T) {
	_, err := database.Exec(fmt.Sprintf(`TRUNCATE table "%s" CASCADE`, table))
	if err != nil {
		t.Fatalf("Failed to truncate %s table: %s", table, err)
	}
}
