package main

import (
	"fmt"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

	// initialize the termui library
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	// create a title bar widget
	title := widgets.NewParagraph()
	title.Text = "Distributed Cache Simulation"
	title.SetRect(0, 0, 60, 3)

	// create two message boxes to display messages from each node
	node1Box := widgets.NewList()
	node1Box.Title = "Node 1"
	node1Box.SetRect(0, 3, 30, 15)

	node2Box := widgets.NewList()
	node2Box.Title = "Node 2"
	node2Box.SetRect(30, 3, 60, 15)

	// add some data to the cache
	cache.Put("foo", "bar")
	cache.Put("hello", "world")

	// start a goroutine for each node to handle message passing
	node1 := make(chan string)
	go func() {
		for {
			key := <-node1
			value, ok := cache.Get(key)
			if ok {
				node1Box.Rows = append(node1Box.Rows, fmt.Sprintf("cache hit for key '%s', value is '%s'\n", key, value))
			} else {
				node1Box.Rows = append(node1Box.Rows, fmt.Sprintf("cache miss for key '%s'\n", key))
			}
			termui.Render(node1Box) // update the node1 message box
		}
	}()

	node2 := make(chan string)
	go func() {
		for {
			key := <-node2
			value, ok := cache.Get(key)
			if ok {
				node2Box.Rows = append(node2Box.Rows, fmt.Sprintf("cache hit for key '%s', value is '%s'\n", key, value))
			} else {
				node2Box.Rows = append(node2Box.Rows, fmt.Sprintf("cache miss for key '%s'\n", key))
			}
			termui.Render(node2Box) // update the node2 message box
		}
	}()

	// create an input box for user input
	inputBox := widgets.NewParagraph()
	inputBox.Title = "Enter a key to retrieve from the cache:"
	inputBox.SetRect(0, 15, 60, 18)

	// create a main loop that listens for user input
	termui.Render(title, node1Box, node2Box, inputBox)
	termuiEvents := termui.PollEvents()
	for {
		select {
		case e := <-termuiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Enter>":
				// get the user input from the input box
				input := inputBox.Text[32:]

				// send the input to the appropriate node
				if input == "foo" || input == "hello" {
					node1 <- input
				} else {
					node2 <- input
				}

				// clear the input box
				inputBox.Text = "Enter a key to retrieve from the cache:"

				termui.Render(title, node1Box, node2Box, inputBox)
			}
		case key := <-node1:
			value, ok := cache.Get(key)
			if ok {
				node1Box.Rows = append(node1Box.Rows, fmt.Sprintf("cache hit for key '%s', value is '%s'", key, value))
			} else {
				node1Box.Rows = append(node1Box.Rows, fmt.Sprintf("cache miss for key '%s'", key))
			}
			termui.Render(node1Box)
		case key := <-node2:
			value, ok := cache.Get(key)
			if ok {
				node2Box.Rows = append(node2Box.Rows, fmt.Sprintf("cache hit for key '%s', value is '%s'", key, value))
			} else {
				node2Box.Rows = append(node2Box.Rows, fmt.Sprintf("cache miss for key '%s'", key))
			}
			termui.Render(node2Box)
		}
	}

}
