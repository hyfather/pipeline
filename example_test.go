package pipeline_test

import (
	"fmt"
	"github.com/hyfather/pipeline"
)

func printStage(inObj interface{}) interface{} {
	fmt.Println(inObj)
	return inObj
}

func squareStage(inObj interface{}) interface{} {
	if v, ok := inObj.(int); ok {
		return v * v
	}
	return nil
}

var pipelineChan chan interface{}

func Example() {
	p := pipeline.New()
	p.AddStageWithFanOut(squareStage, 1)
	p.AddStage(printStage)

	pipelineChan = make(chan interface{}, 10)
	pipelineChan <- 2
	pipelineChan <- 3
	close(pipelineChan)

	<-p.Run(pipelineChan)
	// Output: 4
	// 9
}
