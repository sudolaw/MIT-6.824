package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	//	"bytes"

	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	//	"6.824/labgob"
	"6.824/labrpc"
)

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

//Dol is
const Dol = 1

//Dprintf is
func Dprintf(f string, a ...interface{}) (n int, e error) {
	if Dol > 0 {
		log.Printf(f, a...)

	}
	return
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()
	isLeader  bool
	term      int
	LogD      []ApplyMsg
	timestamp time.Time
	timep     time.Duration
	state     string
	lockTerm  int

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

}

var leTer sync.Mutex

//GetState is ...
// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	ts.Lock()

	leTer.Lock()

	term := rf.term

	isleader := rf.isLeader
	leTer.Unlock()
	ts.Unlock()

	// Your code here (2A).
	return term, isleader

}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
//
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	// Your code here (2D).

	return true
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}

// AppendEntryArgs is
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type AppendEntryArgs struct {
	// Your data here (2A, 2B).
	Index       int
	Candidateid int
}

//AppendEntryReply is
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type AppendEntryReply struct {
	// Your data here (2A).
	Rep bool
}

var ts sync.Mutex
var inlock sync.Mutex

//AppendEntryReply is
func (rf *Raft) AppendEntryReply(args *AppendEntryArgs, reply *AppendEntryReply) {

	inlock.Lock()
	ts.Lock()
	if rf.term < args.Index {
		rf.isLeader = false
		rf.term = args.Index

	}
	rf.lockTerm = args.Index + 1
	rf.timestamp = time.Now()
	rand := (150 + rand.Intn(150))
	rf.timep = time.Duration(rand) * time.Millisecond
	rf.term = args.Index
	reply.Rep = true
	ts.Unlock()
	inlock.Unlock()

}

//AppendEntrymaster is
func (rf *Raft) AppendEntrymaster(server int, args *AppendEntryArgs, reply *AppendEntryReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntryReply", args, reply)
	return ok
}

//AppendEngine is
func (rf *Raft) AppendEngine(lab, laf int) {
	Dprintf("[%v] we have started appending", rf.me)
	for {

		for x := range rf.peers {

			func(xan int) {
				args := AppendEntryArgs{}

				args.Index = laf
				args.Candidateid = lab

				reply := AppendEntryReply{}
				rf.AppendEntrymaster(xan, &args, &reply)

			}(x)

		}
		time.Sleep(5 * time.Millisecond)
	}

}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term        int
	Candidateid int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Reply bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).

	Dprintf("[%v] my  Term is %v", rf.me, rf.term)

	inlock.Lock()

	rf.timestamp = time.Now()
	rand := (150 + rand.Intn(150) + 10000)
	rf.timep = time.Duration(rand) * time.Millisecond

	inlock.Unlock()

	if rf.term < args.Term || args.Candidateid == rf.me {
		reply.Reply = true
	} else {
		reply.Reply = false
	}

	Dprintf("[%v] voted for %v term %v", rf.me, args.Candidateid, args.Term)
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	Dprintf("[%v]:[%v] reply for vote is  : %v", rf.me, server, reply.Reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {

	Dprintf("Start started [%v]", rf.me)
	index := len(rf.LogD)
	term := rf.term
	isLeader := rf.isLeader

	// Your code here (2B).

	return index, term, isLeader
}

//
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
//
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	Dprintf("[%v] here we go lets start ticker ", rf.me)
	for rf.killed() == false {

		// Your code here to check if a leader election should
		// be started and to randomize sleeping time using
		// time.Sleep().
		func() {

			ts.Lock()
			jac := time.Now()
			diff := jac.Sub(rf.timestamp)

			if diff > rf.timep && rf.isLeader != true {

				jaE := time.Now()
				Dprintf("[%v] election beggins is : %v  \n", rf.me, time.Now())
				rf.callElection()
				Dprintf("[%v] election ends is : %v  \n", rf.me, time.Now())
				Dprintf("[%v] election timr is : %v  \n", rf.me, time.Now().Sub(jaE))
			}
			ts.Unlock()

		}()

		time.Sleep(time.Millisecond)
	}

}

var mu sync.Mutex

// callElection is nothing
func (rf *Raft) callElection() {
	leTer.Lock()
	rf.term = rf.term + 1
	rf.lockTerm = 0

	leTer.Unlock()
	count := 0
	finished := 0

	Dprintf("[%v] election started count is : %v  \n", rf.me, count)
	var elec sync.Mutex
	cond := sync.NewCond(&elec)
	for x := range rf.peers {
		go func(xan int) {
			args := RequestVoteArgs{}
			args.Term = rf.term
			args.Candidateid = rf.me
			reply := RequestVoteReply{}
			rf.sendRequestVote(xan, &args, &reply)

			elec.Lock()
			defer elec.Unlock()
			if reply.Reply {
				count++
			}

			finished++
			cond.Broadcast()

		}(x)

	}

	elec.Lock()
	for count < len(rf.peers)/2 || finished < len(rf.peers) {
		cond.Wait()
		fmt.Printf("waiting")
	}

	if count >= (len(rf.peers)/2)+1 && rf.lockTerm == 0 {
		Dprintf("[%v] check check check  %v", rf.me, len(rf.peers)/2)
		rf.isLeader = true
		rf.lockTerm = rf.term

		leTer.Lock()
		name := rf.me
		term := rf.term
		leTer.Unlock()
		go rf.AppendEngine(name, term)

		Dprintf("[%v] won elction_____%v with %v", rf.me, rf.term, rf.isLeader)

	} else {
		leTer.Lock()
		rf.term = rf.term - 1
		leTer.Unlock()
		Dprintf("[%v] lost elction", rf.me)

		ts.Lock()
		rf.term--
		rf.timestamp = time.Now()
		rand := (150 + rand.Intn(150))
		rf.timep = time.Duration(rand) * time.Millisecond

		ts.Unlock()
	}
	elec.Unlock()

}

//Make is jdsfsk
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {

	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.term = 0
	rf.timep = time.Duration((rf.me+1)*200) * time.Millisecond
	rf.timestamp = time.Now()
	rf.lockTerm = 0
	rf.isLeader = false

	// Your initialization code here (2A, 2B, 2C).

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	Dprintf("[%v] Raft started  \n", rf.me)

	// start ticker goroutine to start elections

	go rf.ticker()

	return rf
}
