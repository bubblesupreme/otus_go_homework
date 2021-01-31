package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func makeInteruptable(in In, doneCh Out) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-doneCh:
				return
			default:
			}
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			case <-doneCh:
				return
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out In = in
	for _, s := range stages {
		inWrap := makeInteruptable(out, done)
		out = s(inWrap)
	}

	return out
}
