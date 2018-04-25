# pipeline

This package provides a simplistic implementation of Go pipelines
as outlined in [Go Concurrency Patterns: Pipelines and cancellation.](https://blog.golang.org/pipelines)

# Example Usage

```
import "github.com/hyfather/pipeline"

p := pipeline.New()
p.AddStageWithFanOut(myStage, 10)
p.Run(inChan)
```

More comprehensive examples can be found [here](./examples)
