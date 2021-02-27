package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Coordinator struct {
	file       []string
	sent       []string
	done       []string
	number     int
	waitnumber int
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//

//Nextex s nxt
func (c *Coordinator) Nextex(args *NextArgs, reply *NextReply) error {
	fmt.Printf("*****************************")
	fmt.Printf("SENT %v \n", c.waitnumber)
	fmt.Printf("SUM num%v \n", c.number)
	fmt.Printf("SUM %v \n", len(c.file))
	fmt.Printf("Str %v \n", args.Done)
	if args.Done != "" {
		c.done = append(c.done, args.Done)
		fmt.Printf("Done with %v", args.Done)
		c.waitnumber = c.waitnumber - 1
		fmt.Println("DONE")

	}
	if c.number+c.waitnumber == 0 {
		reply.Bol = false
	} else if c.number == 0 {
		reply.Bol = true
	} else {

		reply.Next = c.file[0]
		c.sent = append(c.sent, c.file[0])
		c.waitnumber = c.waitnumber + 1
		c.number = c.number - 1
		c.file = c.file[1:]
		reply.Bol = true
		fmt.Println("SEND")

	}
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	c.file = files
	c.number = len(files)

	c.server()
	return &c
}
