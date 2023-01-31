package main

import (
	"context"
	"fmt"
	"github.com/bxcodec/dbresolver"
	"github.com/getsentry/sentry-go"
	envConfig "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"poc-order-service/app/handlers"
	"poc-order-service/global/utils/helper"
	kafkadbo "poc-order-service/global/utils/kafka"
	"poc-order-service/global/utils/mongodb"
	opensearch_dbo "poc-order-service/global/utils/opensearch_dbo"
	"poc-order-service/global/utils/redisdb"
	"poc-order-service/global/utils/sqldb"
	"runtime"
	"strings"
	"sync"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	arg := os.Args[1]

	switch arg {
	case "main":
		mainWithoutArg()
		break
	case "command":
		args := []interface{}{}
		args = append(args, os.Args[1], os.Args[2], os.Args[3], os.Args[4])
		commands(args)
	case "consumer":
		args := []interface{}{}
		args = append(args, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5])
		consumers(args)
	default:
		mainWithoutArg()
	}
}

func mainWithoutArg() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := envConfig.Load(".env"); err != nil {
		errStr := fmt.Sprintf(".env not load properly %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	ctx := context.Background()

	//redisdb
	redisDb := redisdb.InitRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DATABASE"))

	//mysql write
	mysqlWrite, err := sqldb.InitSql("mysql", os.Getenv("MYSQL_WRITE_HOST"), os.Getenv("MYSQL_WRITE_PORT"), os.Getenv("MYSQL_WRITE_USERNAME"), os.Getenv("MYSQL_WRITE_PASSWORD"), os.Getenv("MYSQL_WRITE_DATABASE"))
	if err != nil {
		errStr := fmt.Sprintf("Error mysql write connection %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	//mysql read
	mysqlRead, err := sqldb.InitSql("mysql", os.Getenv("MYSQL_READ_01_HOST"), os.Getenv("MYSQL_READ_01_PORT"), os.Getenv("MYSQL_READ_01_USERNAME"), os.Getenv("MYSQL_READ_01_PASSWORD"), os.Getenv("MYSQL_READ_01_DATABASE"))
	if err != nil {
		errStr := fmt.Sprintf("Error mysql read onnection %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	dbConnection := dbresolver.WrapDBs(mysqlWrite.DB(), mysqlRead.DB())
	//mongoDb
	mongoDb := mongodb.InitMongoDB(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_PORT"), ctx)

	//opensearch_dbo
	openSearchHosts := []string{os.Getenv("OPENSEARCH_HOST_01")}
	openSearchClient := opensearch_dbo.InitOpenSearchClientInterface(openSearchHosts, os.Getenv("OPENSEARCH_USERNAME"), os.Getenv("OPENSEARCH_PASSWORD"), ctx)

	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")
	kafkaClient := kafkadbo.InitKafkaClientInterface(ctx, kafkaHosts)

	defer dbConnection.Close()
	defer redisDb.Client().Close()
	defer kafkaClient.GetController().Close()
	defer kafkaClient.GetConnection().Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Printf("Starting POC Order Service HTTP Handler\n")
		handlers.MainHttpHandler(dbConnection, redisDb, mongoDb, openSearchClient, kafkaClient, ctx)
	}()

	wg.Wait()
}

func commands(args []interface{}) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := envConfig.Load(".env"); err != nil {
		errStr := fmt.Sprintf(".env not load properly %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	ctx := context.Background()

	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")
	kafkaClient := kafkadbo.InitKafkaClientInterface(ctx, kafkaHosts)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(args []interface{}) {
		defer wg.Done()
		fmt.Printf("Starting Command Handler\n")
		handlers.MainCommandHandler(kafkaClient, ctx, args)
	}(args)

	wg.Wait()
}

func consumers(args []interface{}) {
	fmt.Println(args)
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := envConfig.Load(".env"); err != nil {
		errStr := fmt.Sprintf(".env not load properly %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	ctx := context.Background()

	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")
	kafkaClient := kafkadbo.InitKafkaClientInterface(ctx, kafkaHosts)

	//opensearch_dbo
	openSearchHosts := []string{os.Getenv("OPENSEARCH_HOST_01")}
	openSearchClient := opensearch_dbo.InitOpenSearchClientInterface(openSearchHosts, os.Getenv("OPENSEARCH_USERNAME"), os.Getenv("OPENSEARCH_PASSWORD"), ctx)

	//redisdb
	redisDb := redisdb.InitRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DATABASE"))

	//mysql write
	mysqlWrite, err := sqldb.InitSql("mysql", os.Getenv("MYSQL_WRITE_HOST"), os.Getenv("MYSQL_WRITE_PORT"), os.Getenv("MYSQL_WRITE_USERNAME"), os.Getenv("MYSQL_WRITE_PASSWORD"), os.Getenv("MYSQL_WRITE_DATABASE"))
	if err != nil {
		errStr := fmt.Sprintf("Error mysql write connection %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	//mysql read
	mysqlRead, err := sqldb.InitSql("mysql", os.Getenv("MYSQL_READ_01_HOST"), os.Getenv("MYSQL_READ_01_PORT"), os.Getenv("MYSQL_READ_01_USERNAME"), os.Getenv("MYSQL_READ_01_PASSWORD"), os.Getenv("MYSQL_READ_01_DATABASE"))
	if err != nil {
		errStr := fmt.Sprintf("Error mysql read onnection %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		panic(err)
	}

	dbConnection := dbresolver.WrapDBs(mysqlWrite.DB(), mysqlRead.DB())

	//mongoDb
	mongoDb := mongodb.InitMongoDB(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_PORT"), ctx)

	defer dbConnection.Close()
	defer redisDb.Client().Close()
	defer kafkaClient.GetController().Close()
	defer kafkaClient.GetConnection().Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(args []interface{}) {
		defer wg.Done()
		fmt.Printf("Starting Consumer Handler\n")
		handlers.MainConsumerHandler(kafkaClient, mongoDb, openSearchClient, dbConnection, redisDb, ctx, args)
	}(args)

	wg.Wait()
}
