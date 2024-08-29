package simqueue

import (
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	//获取实例
	qe := GetEntry()
	//加入延时队列
	qe.Schedule(func() {
		println("hi!")
	}, 3*time.Second, "id_001")
	if time.Now().Year() == 2014 {
		//取消队列
		qe.Cancel("id_001")
	}
	//保持运行
	qe.Run()
}
