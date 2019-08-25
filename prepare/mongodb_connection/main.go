package mongodb_connection

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	Client *mongo.Client
)

func init() {
	opts := options.Client()
	opts = opts.ApplyURI("mongodb://49.235.1.29").SetConnectTimeout(time.Second * 5).SetMaxPoolSize(16)
	var err error
	Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
}
func main() {
	//选择数据库

	db := Client.Database("my_db")
	//选择表
	collection := db.Collection("my_collection")

	fmt.Println(collection.Name())
}
