package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// 任务的执行时间点
type TimePoint struct {
	StartTime int64	`bson:"startTime"`
	EndTime int64	`bson:"endTime"`
}

// 一条日志
type LogRecord struct {
	JobName string	`bson:"jobName"` // 任务名
	Command string `bson:"command"` // shell命令
	Err string `bson:"err"` // 脚本错误
	Content string `bson:"content"`// 脚本输出
	TimePoint TimePoint `bson:"timePoint"`// 执行时间点
}

//jobName过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"`	//JobName赋值为job10
}

func main() {
	//mongodb读取回来的是bson，需要反序列化为LogRecord对象
	var (
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		cond *FindByJobName
		cursor mongo.Cursor	//游标与mysql游标相似
		record *LogRecord
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

	//4.按照jobName字段过滤，想找出jobName=job10，找出2条
	cond = &FindByJobName{JobName:"job10"}	//{"jobName":"job10"}

	//5.查询(过滤+翻页参数)
	cursor,err = collection.Find(context.TODO(),cond,findopt.Skip(0),findopt.Limit(2))
	if err != nil {
		fmt.Println(err)
		return
	}

	//6.遍历结果集
	for cursor.Next(context.TODO()){
		//定义一个日志对象
		record = &LogRecord{}

		//反序列化bson到对象
		if err = cursor.Decode(record);err != nil {
			fmt.Println(err)
			return
		}
		//把日志行打印出来
		fmt.Println(*record)
	}

	//释放游标
	defer cursor.Close(context.TODO())
}
