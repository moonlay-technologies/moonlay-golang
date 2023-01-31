package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
)

type mongoDb struct {
	client *mongo.Client
	ctx    context.Context
}

type MongoDBInterface interface {
	Client() *mongo.Client
	GetCtx() context.Context
}

func InitMongoDB(host, database, username, password, port string, context context.Context) MongoDBInterface {
	usedMongoReplica := os.Getenv("USED_MONGO_REPLICA")
	var mongoURL string

	if usedMongoReplica == "0" {
		portInt, _ := strconv.Atoi(port)
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", username, password, host, portInt, database)
	} else {
		mongoURL = fmt.Sprintf("mongodb+srv://%s:%s@%s", username, password, host)
	}

	fmt.Println(mongoURL)

	client, err := mongo.Connect(context, options.Client().ApplyURI(mongoURL))

	if err != nil {
		panic(err)
	}

	return &mongoDb{
		client: client,
		ctx:    context,
	}
}

func (m *mongoDb) Client() *mongo.Client {
	return m.client
}

func (m *mongoDb) GetCtx() context.Context {
	return m.ctx
}
