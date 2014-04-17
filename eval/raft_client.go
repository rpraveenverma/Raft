package main

import (
	"github.com/rpraveenverma/raft"
	"time"
)

func main() {
	raft.New(101,"configurationfile.json")
	raft.New(102,"configurationfile.json")
	raft.New(103,"configurationfile.json")
	raft.New(104,"configurationfile.json")
	raft.New(105,"configurationfile.json")
	raft.New(106,"configurationfile.json")
	raft.New(107,"configurationfile.json")
	raft.New(108,"configurationfile.json")
	raft.New(109,"configurationfile.json")
	raft.New(110,"configurationfile.json")
	time.Sleep(time.Second*5000)
}
