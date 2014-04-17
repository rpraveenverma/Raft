package raft

import (
	"fmt"
	"github.com/rpraveenverma/cluster"
	"math/rand"
	"time"
	"encoding/json"

)
var noofserverup int
const (
	FOLLOWER  = -1
	CANDIDATE = 0
	LEADER    = 1
)

type Raft interface {
	Term() int
	isLeader() bool
}

type raft struct {
	sid       int
	term      int
	state     int
	serverobj cluster.Servernode
	timeout   time.Duration
	votecount int
}

type VoteReq struct {
	Typ  string
	Sid  int
	Term int
}

type VoteRes struct {
	Typ     string
	Voterid int
	Vote    bool
}

type AppendEntry struct {
	Typ  string
	Sid  int
	Term int
}

func (r raft) Term() int {
	return r.term
}

func (r raft) isLeader() bool {
	return true
}
func New(sid int,filename string) raft {
	var raftinstance raft
	raftinstance.serverobj = cluster.New(sid, filename)
	raftinstance.term = 0
	raftinstance.sid = sid
	raftinstance.timeout = random(300,600)
	raftinstance.state = FOLLOWER
	noofserverup += 1
	go raftinstance.serverresponse()
	return raftinstance
}

func (r raft) serverresponse() {
	for {
		if r.state == FOLLOWER {
			select {
			case <-time.After(r.timeout * time.Millisecond):
				var votereq VoteReq
				votereq.Typ = "votereq"
				votereq.Sid = r.sid
				votereq.Term = r.term+1
				r.term += 1
				fmt.Println(votereq, "sent with term", r.term)
		 		msg, err := json.Marshal(votereq)
				if err != nil {
					panic(err)
				}
			    r.serverobj.Outbox() <- &cluster.Envelope{Pid: cluster.BROADCAST, Msgtype:votereq.Typ, Msg:msg}
				r.state = CANDIDATE
			case instanceenv := (<-r.serverobj.Inbox()):
				switch {
				case instanceenv.Msgtype == "votereq":
					var votereq VoteReq

					err1 := json.Unmarshal(instanceenv.Msg,&votereq)
					if err1 != nil {
						panic("panic in Unmarshling the data")
					}
					if votereq.Term > r.term {
						r.term = votereq.Term
						var voteres VoteRes
						voteres.Typ = "voteres"
						voteres.Voterid = r.sid
						voteres.Vote = true
						fmt.Println(voteres, " sent with term", r.term)
						msg, err := json.Marshal(voteres)
						if err != nil {
							panic(err)
						}

						r.serverobj.Outbox() <- &cluster.Envelope{Pid: votereq.Sid, Msgtype:voteres.Typ, Msg:msg}

					} else if votereq.Term == r.term {
						var voteres VoteRes
						voteres.Typ = "voteres"
						voteres.Voterid = r.sid
						voteres.Vote = false

						msg, err := json.Marshal(voteres)
						if err != nil {
							panic(err)
						}
						fmt.Println(voteres, " sent with term", r.term)
						r.serverobj.Outbox() <- &cluster.Envelope{Pid: votereq.Sid, Msgtype:voteres.Typ, Msg: msg}
					}
				case instanceenv.Msgtype == "heartbeat":
					var appendentry AppendEntry
					err1 := json.Unmarshal(instanceenv.Msg, &appendentry)
					if err1 != nil {
						panic("panic in Unmarshling the data")
					}
					if appendentry.Term >= r.term {
						r.term = appendentry.Term
					}
				}
			case <-time.After(time.Second * 50):
				fmt.Println("channel is stuck")
			}
		}
		if r.state == CANDIDATE {
			select {
			case <-time.After(r.timeout * time.Millisecond):
				r.term += 1
				var votereq VoteReq
				votereq.Typ = "votereq"
				votereq.Sid = r.sid
				votereq.Term = r.term
				msg, err := json.Marshal(votereq)
				if err != nil {
					panic(err)
				}
			r.serverobj.Outbox() <- &cluster.Envelope{Pid: cluster.BROADCAST, Msgtype:votereq.Typ, Msg: msg}

			case instanceenv := (<-r.serverobj.Inbox()):
				switch {
				case instanceenv.Msgtype == "votereq":
					var votereq VoteReq
					err1 := json.Unmarshal(instanceenv.Msg, &votereq)
					if err1 != nil {
						panic("panic in Unmarshling the data")
					}
					if votereq.Term > r.term {
						r.term = votereq.Term
						r.state = FOLLOWER
					} else {
						var voteres VoteRes
						voteres.Typ = "voteres"
						voteres.Voterid = r.sid
						voteres.Vote = false

						msg, err := json.Marshal(voteres)
						if err != nil {
							panic(err)
						}
						fmt.Println(voteres, " sent with term", r.term)
						r.serverobj.Outbox() <- &cluster.Envelope{Pid: votereq.Sid, Msgtype:voteres.Typ, Msg: msg}
					}
				case instanceenv.Msgtype == "voteres":
					var voteres VoteRes
					err1 := json.Unmarshal(instanceenv.Msg, &voteres)
					if err1 != nil {
						panic("panic in Unmarshling the data")
					}
					fmt.Println(voteres, "Recieved on candidate")
					if voteres.Vote == true {
						r.votecount += 1
					}
					if r.votecount > (noofserverup/2) {
						fmt.Print(r.sid, " is leader now")
						r.state = LEADER
					}
				case instanceenv.Msgtype == "heartbeat":
					var appendentry AppendEntry
					err1 := json.Unmarshal(instanceenv.Msg, &appendentry)
					if err1 != nil {
						panic("panic in Unmarshling the data")
					}
					if appendentry.Term >= r.term {
						fmt.Println(r.term, " now Candidate to follower")
						r.state = FOLLOWER
					}
				}
			case <-time.After(time.Second * 50):
				fmt.Println("channel is stuck")
			}
		}
		if r.state == LEADER {
			select {
			case <-time.After(50 * time.Millisecond):
				var appendentry AppendEntry
				appendentry.Typ = "heartbeat"
				appendentry.Sid = r.sid
				appendentry.Term = r.term
				msg, err := json.Marshal(appendentry)
				if err != nil {
					panic(err)
				}
		//	fmt.Println(appendentry," sent with term",appendentry.Term)
			r.serverobj.Outbox() <- &cluster.Envelope{Pid: cluster.BROADCAST, Msgtype:appendentry.Typ, Msg: msg}

			case instanceenv := (<-r.serverobj.Inbox()):
				switch {
				case instanceenv.Msgtype == "heartbeat":
					var appendentry AppendEntry
					err1 := json.Unmarshal(instanceenv.Msg, &appendentry)
					if err1 != nil {
						panic("panic in Unmarshling the data")
					}
					if appendentry.Term > r.term {
						fmt.Println(r.term," now Leader to follower")
						r.state = FOLLOWER
					}
				}
			case <-time.After(time.Second * 50):
				fmt.Println("channel is stuck")
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func random(min, max int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Int63n(int64(max-min)))+time.Duration(min)
}
