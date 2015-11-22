package main

import (
	"fmt"
	"time"

	"github.com/mouadino/go-lymph"
)

var echo = lymph.Client("tcp://127.0.0.1:5555")

func main() {
	c := time.Tick(1 * time.Second)
	i := 0
	for _ = range c {
		text := fmt.Sprintf("foo_%d", i)
		result, err := echo.Call("upper", text)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("%s\n", result.(string))
		}
		i++
	}
}
