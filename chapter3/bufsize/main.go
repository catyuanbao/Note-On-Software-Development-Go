package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    line, isPrefix, err := reader.ReadLine()
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
    if isPrefix {
        fmt.Println("Error: input line too long")
        return
    }
    fmt.Println("Input:", string(line))
}
