package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

//任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime int64 `bson:"endTime"`
}

//一条日志
type LogRecord struct {
	JobName string `bson:"jobName"`//任务名
	Command string `bson:"command"`//shell命令
	Err string `bson:"err"`//脚本错误
	Content string `bson:"content"`//脚本输出
	TimePoint TimePoint `bson:"timePoint"`//执行时间点
}

func main() {
	var (
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		result *mongo.InsertOneResult
		docId objectid.ObjectID
	)
	//1.建立连接
	//连接是一个网络操作，如果网络不畅通，context.TODO()可用于取消连接，
	client, err = mongo.Connect(context.TODO(),"mongodb://120.78.140.135:27017",clientopt.ConnectTimeout(5 * time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	//2.选择数据库my_db
	database = client.Database("cron")

	//3.选择表my_collection
	collection = database.Collection("log")

	//4.插入记录(bson)
	record = &LogRecord{
		JobName:"job10",
		Command:"echo hello",
		Err:"",
		Content:"hello",
		TimePoint: TimePoint{StartTime:time.Now().Unix(),EndTime:time.Now().Unix() + 10},
	}
	result,err = collection.InsertOne(context.TODO(),record)
	if err != nil {
		fmt.Println(err)
		return
	}

	//_id:默认生成一个全局唯一ID，ObjectID：12字节的二进制
	docId = result.InsertedID.(objectid.ObjectID) //反射的断言
	fmt.Println("自增ID:", docId.Hex()) //转化为16进制
}
