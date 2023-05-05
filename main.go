package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Define a Cache struct to hold key-value pairs
type Cache struct {
	cache map[string]string
}

// Get a value from the cache given a key
func (c *Cache) Get(key string) (string, bool) {
	val, ok := c.cache[key]
	return val, ok
}

// Store a value in the cache given a key
func (c *Cache) Put(key string, value string) {
	c.cache[key] = value
}

// create new cache instance
func NewCache(filename string) (*Cache, error) {
	cache := make(map[string]string)

	//open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//include scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		//parse each line into key and value
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid data format in file: %s", line)
		}
		//add the key-value pair to the cache
		cache[parts[0]] = parts[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Cache{cache: cache}, nil

}

func main() {
	// create a new cache
	cache, err := NewCache("data.txt")
	if err != nil {
		log.Fatal(err)

	}

	// simulate a distributed system with two nodes
	node1 := make(chan string)
	node2 := make(chan string)

	// Node 1 goroutine: retrieve values from the cache using the Get method

	go func() {
		for {
			key := <-node1              // Wait for a request on the node1 channel
			value, ok := cache.Get(key) // Retrieve the value from the cache
			if ok {
				fmt.Println("\nnode 1: cache hit for key", key, "value", value)
			} else {
				fmt.Println("\nnode 1: cache miss for key", key)
			}
		}
	}()

	// node 2
	go func() {
		for {
			key := <-node2
			value, ok := cache.Get(key)
			if ok {
				fmt.Println("\nnode 2: cache hit for key", key, "value is", value)
			} else {
				fmt.Println("\nnode 2: cache miss for key", key)
			}
		}
	}()

	// // add some data to the cache
	// cache.Put("hello", "world")
	// cache.Put("country", "kenya")
	// cache.Put("food", "ugali")
	// cache.Put("goodbye", "world")

	//to simulate requests to the cache:
	//node1 <- "next"
	//node2 <- "hello"
	//node1 <- "hello"
	//node2 <- "foo"

	//to prompt for user input to the cache:

	for {
		var key string
		fmt.Print("Enter a key to retrieve from node 1: ")
		fmt.Scanln(&key)
		node1 <- key // Request the key from node 1

		if key == "exit" {
			// Signal node 1 to exit loop and terminate the program
			node1 <- key
			break
		}

		fmt.Print("Enter a key to retrive from node 2:")
		fmt.Scanln(&key)
		//request key from node 2
		node2 <- key

		if key == "exit" {
			//signal node 2 to end loop and terminate the program
			node2 <- key
			break
		}

	}

}
