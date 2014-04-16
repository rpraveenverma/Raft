package main

import (
	"github.com/rpraveenverma/raft"
	"time"
)

func main() {
	raft.New(101,5, "configurationfile.json")
	raft.New(102,40, "configurationfile.json")
	raft.New(103, 40,"configurationfile.json")
	raft.New(104, 40,"configurationfile.json")
	raft.New(105, 40,"configurationfile.json")
	time.Sleep(50*time.Second)
}
