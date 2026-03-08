package main

import (
    "fmt"
)

func cleanupOrder() {
    defer fmt.Println("defer 1")
    defer fmt.Println("defer 2")
    defer fmt.Println("defer 3")
    fmt.Println("body")
}

func recoverAtBoundary(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    fn()
    return nil
}

func main() {
    cleanupOrder()

    err := recoverAtBoundary(func() {
        panic("boom")
    })
    fmt.Println(err)

    err = recoverAtBoundary(func() {
        fmt.Println("normal execution")
    })
    fmt.Println(err)
}
