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
	NumBottles int
}
type Response struct {
	Finished bool
}
type Operations struct{}

func main() {
	//FLAGS
	buddyNum = flag.String("buddyNum", "NIL", "Number in sequence of instances")
	bottles := flag.Int("bottles", 0, "Bottles of Beer (launches song if not 0)")
	//port := flag.String("l_address", "8030", "port that this instance needs to listen on")
	nextIP = flag.String("nextIP", "", "port that this instance calls")
	flag.Parse()

	//DEBUGGING
	fmt.Println("bottles: " + strconv.Itoa(*bottles))
	fmt.Println("buddyNum: " + *buddyNum)
	fmt.Println("nextIP: " + *nextIP)

	//Starts server which is taking connections via gate 8030

	listener, err := net.Listen("tcp", ":8030")
	fmt.Println("Listening on port :8030")
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
	err = rpc.Register(new(Operations))
	if err != nil {
		log.Fatal("Error with rpc.Register()")
	}

	//Constantly be listening for remote calls
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			fmt.Println("Listening on Listener...")
			rpc.Accept(listener)
		}
		wg.Done()
		return
	}()

	//Start if not already started
	//Get connection to next client in chain...
	fmt.Println("bottles: " + strconv.Itoa(*bottles))
	if *bottles != 0 {
		client, err := rpc.Dial("tcp", *nextIP)
		if err != nil {
			log.Fatal("Error with Dial()")
		}

		request := Request{NumBottles: *bottles}
		response := Response{}
		err = client.Call("Operations.CallNextInChain", request, response)
		fmt.Println("Attempted RPC...")
		if err != nil {
			fmt.Println(err)
			log.Fatal("Error with Call()")
		}

		if response.Finished {
			defer os.Exit(0)
		}
	}
	wg.Wait()
	return
}

func (s *Operations) CallNextInChain(req Request, res *Response) (err error) {
	if req.NumBottles == 0 {
		fmt.Println("No bottles of beer on the wall :D")
		os.Exit(0)
	}

	n := strconv.Itoa(req.NumBottles)
	fmt.Println("Buddy " + *buddyNum + ": " + n + " bottles of beer on the wall, " + n + " bottles of beer. Take one down, pass it around...")

	//Get connection to next client in chain...
	client, err := rpc.Dial("tcp", *nextIP)
	if err != nil {
		log.Fatal("Error with Dial()")
	}
	fmt.Println("Connection made to client...")

	//Make RPC to next in chain
	response := Response{}
	err = client.Call("Operations.CallNextInChain", Request{NumBottles: req.NumBottles - 1}, response)
	fmt.Println("Attempted RPC...")
	if err != nil {
		log.Fatal("Error with Call()")
	}

	if response.Finished {
		res.Finished = true
		defer os.Exit(0)
	}
	return nil
}
