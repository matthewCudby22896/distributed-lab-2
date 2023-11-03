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
		"18.232.83.36:8030"}

	connections := make([]Connection, len(serverIP))

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

		connection := Connection{client, make(chan string), make(chan Output)}
		connections = append(connections, connection)
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

	aggChan := make(chan Output)

	//Start worker go routines
	for _, c := range connections {
		go worker(c, aggChan)
		scanner.Scan()
		//Send first piece of work
		text := scanner.Text()
		c.In <- text
	}

	//Process all words in text doc
	for scanner.Scan() {
		output := <-aggChan
		output.In <- scanner.Text()
		fmt.Println("Output: " + output.Result)
	}

	//Send close message to
	for _, c := range connections {
		c.In <- "QUIT"
		<-c.Out
	}

	fmt.Println("FINISHED")

}
func worker(c Connection, aggChan chan Output) {
	fmt.Println("Worker Started...")
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
