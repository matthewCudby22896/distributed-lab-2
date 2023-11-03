package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

type Connection struct {
	Client *rpc.Client
	In     chan string
	Out    chan Output
}

type Output struct {
	In     chan string
	Result string
}

func main() {
	//Private IP addresses of other aws instances
	serverIP := []string{
		"18.206.200.153:8030",
		"18.232.83.36:8030",
	}

	var connections []Connection
	fmt.Println(connections)
	//Open file & initialise scanner
	file, err := os.Open("wordlist")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error with os.Open()")
	}
	scanner := bufio.NewScanner(file)

	fmt.Println("FILE OPENED")

	//Get connection to every aws instance
	for _, ip := range serverIP {
		client, err := rpc.Dial("tcp", ip)
		if err != nil {
			fmt.Println(err)
			log.Fatal("Error with Dial()")
		}
		fmt.Println("Connection made...")

		connection := Connection{Client: client, In: make(chan string), Out: make(chan Output)}
		fmt.Println(connection)
		connections = append(connections, connection)
		fmt.Println(connections)
	}

	//Close all connections when method has finished
	defer func() {
		for _, connection := range connections {
			err := connection.Client.Close()
			if err != nil {
				fmt.Println(err)
				log.Fatal("Error with Close()")
			}
		}
	}()

	aggChan := make(chan Output, 10)

	//Start worker go routines
	fmt.Println(connections)
	for _, c := range connections {
		fmt.Println(c)
		go worker(c, aggChan)
		scanner.Scan()
		//Send first piece of work
		text := scanner.Text()
		fmt.Println("text: " + text)
		c.In <- text
		fmt.Println("Input sent...")
	}

	//Process all words in text doc
	for scanner.Scan() {
		output := <-aggChan
		output.In <- scanner.Text()
		fmt.Println("Output: " + output.Result)
	}

	var output Output
	var notEmpty bool = false
	//Empty out aggregate channel
	for notEmpty {
		select {
		case output = <-aggChan:
			fmt.Println("Output: " + output.Result)
		default:
			notEmpty = false
		}

	}
	fmt.Println("Aggregate channel emptied...")
	//Send close message to each connection
	for _, c := range connections {
		fmt.Println(c)
		c.In <- "QUIT"
		<-c.Out
	}

	fmt.Println("FINISHED")

}
func worker(c Connection, aggChan chan Output) {
	fmt.Println("Worker started...")
	fmt.Println(c)
	for {
		input := <-c.In
		fmt.Println("Input received...")
		if input == "QUIT" {
			c.Out <- Output{}
			fmt.Println("Worker returning...")
			return
		}

		request := stubs.Request{Message: input}
		response := new(stubs.Response)

		fmt.Println("Sending:" + request.Message)

		err := c.Client.Call(stubs.PremiumReverseHandler, request, response)
		if err != nil {
			fmt.Println(err)
			log.Fatal("Error with Call()")
		}

		output := Output{In: c.In, Result: response.Message}
		aggChan <- output
	}
}
