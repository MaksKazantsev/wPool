# wPool - Golang realization of a simple worker pull.
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
		time.Sleep(time.Second * 5) // Give time to worker pool to execute.
		pool.Stop() // 
        }()
	
	for i := 0; i < 10; i++ { 
		pool.Task(func() error { // Task new function.
			// Your task.
		    return nil
		})
        }
	
	// Catch all errors from the pool.
	for err := range pool.AwaitError() {
		fmt.Println(err)
        }
}
