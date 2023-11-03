package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

var nextIP *string
var buddyNum *string

type Request struct {
	numBottles int
}
type Response struct {
	finished bool
}
type Operations struct{}

func main() {
	//FLAGS
	buddyNum = flag.String("buddy", "NIL", "Number in sequence of instances")
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	port := flag.String("l_address", "8030", "port that this instance needs to listen on")
	nextIP = flag.String("c_address", "", "port that this instance calls")
	flag.Parse()

	//TODO: Up to you from here! Remember, you'll need to both listen for
	//Get connection to next client in chain...

	//Starts server which is taking connections via gate 8030
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal("Error with net.Listen()")
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatal("Error with net.Close()")
		}
	}()

	//Register Operations struct and it's methods with rpc
	err = rpc.RegisterName("Operations", new(Operations))
	if err != nil {
		log.Fatal("Error with rpc.RegisterName()")
	}

	//Constantly be listening for remote calls
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			rpc.Accept(listener)
		}
		wg.Done()
		return
	}()

	//Start if not already started
	//Get connection to next client in chain...

	if *bottles != 0 {
		client, err := rpc.Dial("tcp", *nextIP)
		if err != nil {
			log.Fatal("Error with Dial()")
		}

		request := Request{numBottles: *bottles}
		response := Response{}
		err = client.Call("Operations.callNextInChain", request, response)
		if err != nil {
			log.Fatal("Error with Call()")
		}
		if response.finished {
			defer os.Exit(0)
		}
	}
	wg.Wait()
	return
}

func (s *Operations) CallNextInChain(req *Request, res *Response) error {
	if req.numBottles == 0 {
		fmt.Println("No bottles of beer on the wall :D")
		os.Exit(0)
	}

	n := strconv.Itoa(req.numBottles)
	fmt.Println("Buddy " + *buddyNum + ": " + n + " bottles of beer on the wall, " + n + " bottles of beer. Take one down, pass it around...")

	//Get connection to next client in chain...
	client, err := rpc.Dial("tcp", *nextIP)
	if err != nil {
		log.Fatal("Error with Dial()")
	}

	//Make RPC to next in chain
	response := Response{}
	err = client.Call("Operations.callNextInChain", Request{numBottles: req.numBottles - 1}, response)
	if err != nil {
		log.Fatal("Error with Call()")
	}

	if response.finished {
		res.finished = true
		defer os.Exit(0)
	}
	return nil
}
