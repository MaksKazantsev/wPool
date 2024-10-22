# wPool - Golang realization of a simple worker pool.
## Usage:
### Import package with:
go get github.com/MaksKazantsev/wPool
### Getting started:
```go
package main

import(
	"github.com/MaksKazantsev/wPool"
	"fmt"
	"time"
)

func main() {
	pool := wpool.NewPool(5) // Init new worker pull, function accepts workers capacity.
	
	go func() {
		time.Sleep(time.Second * 5) // Give time to worker pool to task all functions.
		pool.Stop() // Stop worker pool.
        }()
	
	for i := 0; i < 15; i++ { 
		pool.Task(func() error { // Task new function.
			// Your task.
		    return nil
		})
        }
	
	// Catch all errors from the pool.
	for err := range pool.CatchError() {
		log.Println(err)
        }
}
