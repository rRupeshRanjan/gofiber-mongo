package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gofiber-mongo/domain"
	"log"
	"os"
)

var (
	ServerPort          string
	MongoServerAddress  string
	MongoDatabaseName   string
	MongoCollectionName string
	LogFile             *os.File
	AppLogger           *zap.Logger
	FiberLogFormat      string
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Error reading config file: " + err.Error())
	}

	ServerPort = domain.ColonPort + viper.GetString(domain.ServerPort)
	MongoServerAddress = viper.GetString(domain.MongoServerAddress)
	MongoDatabaseName = viper.GetString(domain.MongoDatabaseName)
	MongoCollectionName = viper.GetString(domain.MongoCollectionName)
	FiberLogFormat = viper.GetString(domain.FiberLogFormat)
	LogFile, _ = os.OpenFile(viper.GetString(domain.AppLogLocation), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(LogFile),
		zap.InfoLevel)
	AppLogger = zap.New(core)
}
