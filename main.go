package main

import (
	"fmt"
)

// Define a Cache struct to hold key-value pairs
type Cache struct {
	cache map[string]string
}

// create  a new cache instance
func NewCache() *Cache {
	return &Cache{cache: make(map[string]string)}
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

func main() {
	// create a new cache
	cache := NewCache()

	// simulate a distributed system with two nodes
	node1 := make(chan string)
	node2 := make(chan string)

	// Node 1 goroutine: retrieve values from the cache using the Get method

	go func() {
		for {
			key := <-node1              // Wait for a request on the node1 channel
			value, ok := cache.Get(key) // Retrieve the value from the cache
			if ok {
				fmt.Println("node 1: cache hit for key", key, "value is", value)
			} else {
				fmt.Println("node 1: cache miss for key", key)
			}
		}
	}()

	// node 2
	go func() {
		for {
			key := <-node2
			value, ok := cache.Get(key)
			if ok {
				fmt.Println("node 2: cache hit for key", key, "value is", value)
			} else {
				fmt.Println("node 2: cache miss for key", key)
			}
		}
	}()

	// add some data to the cache
	cache.Put("foo", "bar")
	cache.Put("hello", "world")

	// simulate requests to the cache
	//node1 <- "next"
	//node2 <- "hello"
	//node1 <- "hello"
	//node2 <- "foo"
	for {
		var key string
		fmt.Print("Enter a key to retrieve from the cache: ")
		fmt.Scanln(&key)
		node1 <- key // Request the key from node 1
	}
}
