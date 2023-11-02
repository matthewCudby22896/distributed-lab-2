package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/rpc"
	"os"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	fmt.Println("Server: ", *server)

	//Need to connect to RPC server and send the request
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	/*request := stubs.Request{Message: "Hello"}
	response := new(stubs.Response)

	client.Call(stubs.PremiumReverseHandler, request, response)
	fmt.Println("Responded: " + response.Message)*/

	//Read word from list
	//Send them off to client to be reversed

	file, _ := os.Open("wordlist")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		//fmt.Println(text)
		request := stubs.Request{Message: text}
		response := new(stubs.Response)

		//Call() waits for it's completion
		client.Call(stubs.PremiumReverseHandler, request, response)
		fmt.Println("Responded: " + response.Message)
	}
}
