package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func ReverseString(s string, i int) string {
	time.Sleep(time.Duration(rand.Intn(i)) * time.Second)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type SecretStringOperations struct{}

func (s *SecretStringOperations) Reverse(req stubs.Request, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 10)
	return
}

func (s *SecretStringOperations) FastReverse(req stubs.Response, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 2)
	return
}
func main() {
	//Adds optional flag with default value
	pAddr := flag.String("port", "8030", "Port to listen to")
	flag.Parse()

	//Initialising random
	rand.Seed(time.Now().UnixNano())

	//Register the SecretStringOperations struct and it's methods with rpc
	rpc.Register(&SecretStringOperations{})

	//Listen() function creates servers
	listener, _ := net.Listen("tcp", ":"+*pAddr)

	defer listener.Close()
	//Accept accepts connections on the listener
	rpc.Accept(listener)

}
