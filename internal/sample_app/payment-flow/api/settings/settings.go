package settings

import (
	"os"
	"sync"
)

const (
	ApiName     = "payment-flow-api"
	UserID      = "b55fa78a-67c5-43ad-b28f-e97e1aad7247"
	ServiceName = "payment-flow-api"
	CDCSlotName = "paymentflowcdc"
)

var (
	ServerPort = 9876
	StackName  string
	LogFormat  string
	LogLevel   string
)

var settingsOnce sync.Once

func Configure() {
	settingsOnce.Do(func() {
		StackName = GetStringOrDefault("STACK_NAME", "local")
		LogFormat = os.Getenv("LOG_FORMAT")
		LogLevel = os.Getenv("LOG_LEVEL")
	})
}

func GetStringOrDefault(envName, defaultVal string) string {
	if os.Getenv(envName) != "" {
		return os.Getenv(envName)
	}
	return defaultVal
}
