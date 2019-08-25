package main

import (
	"context"
	mClient "corntab/prepare/mongodb_connection"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

//任务执行时间点
type timePoint struct {
	StartTime int64 `bson:"start_time"`
	EndTime   int64 `bson:"end_time"`
}
type LogRecord struct {
	JobName string `bson:"job_name"` // 任务名
	Command string `bson:"command"`  //shell 命令
	ErrMsg  string `bson:"err_msg"`  //脚本错误
	Content string `bson:"content"`  //脚本输出

	TimePoint timePoint `bson:"time_point"`
}
type FindByName struct {
	JobName string `bson:"job_name"`
}

type timeBeforeCond struct {
	BeforeTime int64 `bson:"$lt"`
}
type DeleteByStartTime struct {
	BeforeCond timeBeforeCond `bson:"time_point.start_time"`
}

//mongodb 存储的是bson
func main() {
	//选择数据库
	db := mClient.Client.Database("my_db")
	//选择数据库表
	collection := db.Collection("my_collection")

	insertIntoCollection(collection)
	queryFromCollection(collection)
}

func insertIntoCollection(collection *mongo.Collection) {

	record := &LogRecord{
		JobName:   "test_insert",
		Command:   "ls",
		ErrMsg:    "",
		Content:   "xxx",
		TimePoint: timePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Add(time.Second * 5).Unix()},
	}

	//*InsertOneResult
	result, err := collection.InsertOne(context.TODO(), record)

	if err != nil {
		log.Printf("insert record error:%v", err)
	}

	//自增id，默认生成一个12字节的二进制
	recordId := result.InsertedID.(primitive.ObjectID)
	fmt.Println(recordId.Hex())
}

func queryFromCollection(collection *mongo.Collection) {

	//findCondition := &LogRecord{JobName:"test_insert"}//直接使用记录查询查不到结果集
	condition := &FindByName{JobName: "test_insert"}

	result, err := collection.Find(context.TODO(), condition)
	if err != nil {
		log.Printf("find result error:%v", err)
	}

	//循环遍历结果集
	for result.Next(context.TODO()) {
		record := &LogRecord{}
		err = result.Decode(record)
		if err != nil {
			fmt.Printf("decode failed,error:%v", err)
		}
		fmt.Println(record)
	}
}

func deleteByTime(collection *mongo.Collection) {
	dc := &DeleteByStartTime{
		BeforeCond: timeBeforeCond{BeforeTime: time.Now().Unix()},
	}

	result, err := collection.DeleteMany(context.TODO(), dc)
	if err != nil {
		fmt.Printf("delete error:%v", err)
	}

	fmt.Printf("delete count:%d", result)
}
