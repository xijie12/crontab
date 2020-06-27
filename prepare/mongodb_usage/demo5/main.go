package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"fmt"
)

//startTime小于某时间
//{"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}
//{"timePoint.startTime":{"$lt:timestamp}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond	`bson:"timePoint.startTime"`
}

func main() {
	var (
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		delCond *DeleteCond
		delResult *mongo.DeleteResult
	)

	// 1, 建立连接
	client,err = mongo.Connect(context.TODO(),"mongodb://120.78.140.135:27017",clientopt.ConnectTimeout(5 * time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2, 选择数据库my_db
	database = client.Database("cron")

	// 3, 选择表my_collection
	collection = database.Collection("log")

	//4.要删除开始时间早于当前时间的所有日志($lt是less than 小于操作符)
	//delete({"timePoint.startTime":{"$lt":当前时间}}) //将定义的结构体转化为bson对象
	delCond = &DeleteCond{BeforeCond:TimeBeforeCond{Before:time.Now().Unix()}}

	//5.执行删除
	delResult,err = collection.DeleteMany(context.TODO(),delCond)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("删除的行数：",delResult.DeletedCount)
}