package mr

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"os"
	"plugin"
	"time"
)

//KeyValue is knfdlk
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//Worker is
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	mapf, reducef = loadPlugin(os.Args[1])
	// Your worker implementation here.

	jac := true
	xa := ""

	for jac == true {

		jac, xa = mapcall(xa, mapf)
		fmt.Printf("recived %v and %v\n", xa, jac)
		time.Sleep(time.Second)

	}
	// uncomment to send the Example RPC to the coordinator.

	fmt.Println("********************MAP is Done***************")

}

func mapcall(x string, mapf func(string, string) []KeyValue) (bool, string) {

	//intermediate := []KeyValue{}

	filenam := Saynext(x)

	filename := filenam.Next

	/*file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := mapf(filename, string(content))
	intermediate = append(intermediate, kva...)*/
	fmt.Println(filename)
	return filenam.Bol, filenam.Next

}

//Saynext is lol
func Saynext(x string) NextReply {

	// declare an argument structure.
	args := NextArgs{}

	// fill in the argument(s).
	args.Done = x

	// declare a reply structure.
	reply := NextReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.Nextex", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Ne %v\n", reply.Next)
	fmt.Printf("reply.Nnol %v\n", reply.Bol)

	return reply
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func loadPlugin(filename string) (func(string, string) []KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v", filename)
	}
	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	mapf := xmapf.(func(string, string) []KeyValue)
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatalf("cannot find Reduce in %v", filename)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}
