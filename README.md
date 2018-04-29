[![GoDoc](https://godoc.org/github.com/hyfather/pipeline?status.svg)](https://godoc.org/github.com/hyfather/pipeline)
[![Build Status](https://travis-ci.org/hyfather/pipeline.svg?branch=master)](https://travis-ci.org/hyfather/pipeline)
# pipeline

This package provides a simplistic implementation of Go pipelines
as outlined in [Go Concurrency Patterns: Pipelines and cancellation.](https://blog.golang.org/pipelines)

# Docs
GoDoc available [here.](https://godoc.org/github.com/hyfather/pipeline)

# Example Usage

```
import "github.com/hyfather/pipeline"

p := pipeline.New()
p.AddStageWithFanOut(myStage, 10)
p.AddStageWithFanOut(anotherStage, 100)
doneChan := p.Run(inChan)

<- doneChan
```

More comprehensive examples can be found [here.](./examples)
