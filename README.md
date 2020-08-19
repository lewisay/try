# Try [![Go Report Card](https://goreportcard.com/badge/github.com/lewisay/try)](https://goreportcard.com/report/github.com/lewisay/try) [![GitHub release](https://img.shields.io/github/release/lewisay/try.svg?style=flat-square)](https://github.com/lewisay/try/releases) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/lewisay/try?tab=doc)




Idiomatic Go retry package.

## Example
```go
package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/lewisay/try"
)

func main() {
	// before
	var err error
	for i := 0; i < 3; i++ {
		err = doSomeThing()
		if err != nil {
			break
		}
	}
	if err != nil {
		log.Println("error:", err)
	}

	// after
	err = try.Do(context.TODO(), func(attempt int) (retry bool, err error) {
		retry = attempt < 3 // try 3 times
		return retry, doSomeThing()
	})
	if err != nil {
		log.Println("error:", err)
	}
}

func doSomeThing() error {
	time.Sleep(1 * time.Second)
	return errors.New("something error")
}
```
