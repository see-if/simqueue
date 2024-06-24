package main

import (
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	qe := GetEntry()
	qe.Schedule(func() {
		println("hi!")
	}, 3*time.Second, "id_001")
	qe.Run()
}
