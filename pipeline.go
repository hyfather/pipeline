// Package pipeline provides a simplistic implementation of pipelines
// as outlined in https://blog.golang.org/pipelines
package pipeline

import (
	"sync"
)

// Pipeline type defines a pipeline to which stages can be added and run.
type Pipeline []StageFn

// StageFn is a lower level type that chains together multiple stages
// using channels.
type StageFn func(inChan <-chan interface{}) (outChan chan interface{})

// ProcessFn types are defined by users of the package and passed in
// to instantiate a meaningful pipeline.
type ProcessFn func(inObj interface{}) (outObj interface{})

// New is a convenience method of create a new instance of the Pipeline type
func New() Pipeline {
	return Pipeline{}
}

// AddStage is a convenience method for adding a stage with fanSize = 1
// See AddStageWithFanOut for more information.
func (p *Pipeline) AddStage(inFunc ProcessFn) {
	*p = append(*p, fanningStageFnFactory(inFunc, 1))
}

// AddStageWithFanOut adds a parallel fan-out ProcessFn to the pipeline. The
// fanSize number indicates how many instances of this stage will read from the
// previous stage and process the data flowing through simultaneously to take
// advantage of parallel CPU usage.
// Most pipelines will have multiple stages, and the order in which AddStageWithFanOut
// AddStage is invoked matters -- the first invocation indicates the first stage
// and so forth.
func (p *Pipeline) AddStageWithFanOut(inFunc ProcessFn, fanSize uint64) {
	*p = append(*p, fanningStageFnFactory(inFunc, fanSize))
}

// Run starts the pipeline with all the stages that have been added. Run is not
// a blocking function and will return immediately with a doneChan. Consumers
// can wait on the doneChan for an indication of when the pipeline has completed
// processing.
// The pipeline runs until its `inChan` channel is open. Once the `inChan` is closed,
// the pipeline stages will sequentially complete from the first stage to the last.
// Once all stages are complete, the last outChan is drained and the doneChan is closed.
func (p *Pipeline) Run(inChan <-chan interface{}) (doneChan chan struct{}) {
	for _, stage := range *p {
		inChan = stage(inChan)
	}

	doneChan = make(chan struct{})
	go func() {
		defer close(doneChan)
		for range inChan {
			// pull objects from inChan so that the gc marks them
		}
	}()
	return
}

// stageFnFactory makes a standard stage function from a given ProcessFn.
// Stage functions accept an inChan and return an outChan, allowing us to
// chain multiple functions into a pipeline.
func stageFnFactory(inFunc ProcessFn) (outFunc StageFn) {
	return func(inChan <-chan interface{}) (outChan chan interface{}) {
		outChan = make(chan interface{})
		go func() {
			defer close(outChan)
			for inObj := range inChan {
				if outObj := inFunc(inObj); outObj != nil {
					outChan <- outObj
				}
			}
		}()
		return
	}
}

// fanningStageFnFactory makes a stage function that fans into multiple
// goroutines increasing the stage throughput depending on the CPU
func fanningStageFnFactory(inFunc ProcessFn, fanSize uint64) (outFunc StageFn) {
	return func(inChan <-chan interface{}) (outChan chan interface{}) {
		var channels []chan interface{}
		for i := uint64(0); i < fanSize; i++ {
			channels = append(channels, stageFnFactory(inFunc)(inChan))
		}
		outChan = fanIn(channels)
		return
	}
}

// fanIn merges an array of channels into one channel
func fanIn(channels []chan interface{}) (outChan chan interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	outChan = make(chan interface{})
	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			defer wg.Done()
			for obj := range ch {
				outChan <- obj
			}
		}(ch)
	}

	go func() {
		defer close(outChan)
		wg.Wait()
	}()
	return
}