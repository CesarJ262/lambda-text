package main

import (
    "fmt"
    "context"
)

func handler(ctx context.Context, event interface{}) (interface{}, error) {
    fmt.Println("Hello from Lambda!")
    return map[string]interface{}{
        "statusCode": 200,
        "body":       "Hello Addfly",
    }, nil
}