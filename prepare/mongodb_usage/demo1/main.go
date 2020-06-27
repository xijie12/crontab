package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"fmt"
)

func main() {
	var (
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
	)
	//1.建立连接
	//连接是一个网络操作，如果网络不畅通，context.TODO()可用于取消连接，
	client, err = mongo.Connect(context.TODO(),"mongodb://120.78.140.135:27017",clientopt.ConnectTimeout(5 * time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	//2.选择数据库my_db
	database = client.Database("my_db")

	//3.选择表my_collection
	collection = database.Collection("my_collection")

	fmt.Println(collection)
}
