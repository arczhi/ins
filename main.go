package main

import (
	"fmt"
	"time"
)

func main() {
	i := New()

	fmt.Println(i.Get("name"))

	i.Set("name", "tom")
	fmt.Println(i.Get("name"))
	time.Sleep(2 * time.Second)
	fmt.Println(i.Get("name"))

	time.Sleep(10 * time.Second)
}
