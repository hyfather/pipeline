package pipeline_test

import (
	"fmt"
	"github.com/hyfather/pipeline"
	"sort"
)

func ExampleMergeChannels() {
	inChan1 := make(chan interface{}, 10)
	inChan2 := make(chan interface{}, 10)

	inChan1 <- 1
	inChan1 <- 3
	inChan1 <- 5
	inChan1 <- 7
	inChan1 <- 9
	close(inChan1)

	inChan2 <- 2
	inChan2 <- 4
	inChan2 <- 6
	inChan2 <- 8
	inChan2 <- 10
	close(inChan2)

	outChan := pipeline.MergeChannels([]chan interface{}{inChan1, inChan2})

	var ints []int
	for e := range outChan {
		ints = append(ints, e.(int))
	}
	sort.Ints(ints)
	fmt.Println(ints)

	// Output: [1 2 3 4 5 6 7 8 9 10]
}
