package startup

import (
	"kitbook/pkg/logger"
)

func InitLogger() logger.Logger {

	return logger.NewNopLogger()
	//return ioc.InitLogger()
}
