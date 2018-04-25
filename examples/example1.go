package main

import (
	"fmt"
	"github.com/hyfather/pipeline"
	"time"
)

func add1(inObj interface{}) interface{} {
	if v, ok := inObj.(int); ok {
		return v + 1
	}
	return nil
}

func slowSquare(inObj interface{}) interface{} {
	time.Sleep(1 * time.Second)
	if v, ok := inObj.(int); ok {
		return v * v
	}
	return nil
}

func slowPrint(inObj interface{}) interface{} {
	time.Sleep(1 * time.Second)
	if v, ok := inObj.(int); ok {
		fmt.Println(v)
		return v * v
	}
	return nil
}

func main() {
	fmt.Println("Starting")
	ch := make(chan interface{}, 100)
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)

	p := pipeline.New()
	p.AddStageWithFanOut(add1, 1)
	p.AddStageWithFanOut(slowSquare, 7)
	p.AddStageWithFanOut(slowPrint, 10)

	fmt.Println("Waiting")
	<-p.Run(ch)
	fmt.Println("Done")
}
