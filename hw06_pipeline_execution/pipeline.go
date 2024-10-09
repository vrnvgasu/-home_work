package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	results := in

	for _, stage := range stages {
		results = stage(executeData(results, done))
	}

	return results
}

func executeData(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for result := range in {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case out <- result:
			}
		}
	}()

	return out
}
