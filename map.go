package pipes

type MapFunc[T any, N any] func(T) N
type TapFunc[T any] MapFunc[T, T]

func Map[T any, N any](size int, mp MapFunc[T, N], in <-chan T) <-chan N {
	out := make(chan N, size)
	go mapWorker(mp, in, out)
	return out
}

func mapWorker[T any, N any](mp MapFunc[T, N], in <-chan T, out chan<- N) {
	defer close(out)
	for t := range in {
		out <- mp(t)
	}
}

type MapWithErrorFunc[T any, N any] func(T) (N, error)

func MapWithError[T any, N any](size int, mp MapWithErrorFunc[T, N], in <-chan T) (<-chan N, <-chan error) {
	out, err := make(chan N, size), make(chan error, size)
	go mapWithErrorWorker(mp, in, out, err)
	return out, err
}

func mapWithErrorWorker[T any, N any](mp MapWithErrorFunc[T, N], in <-chan T, out chan<- N, err chan<- error) {
	defer func() { close(out); close(err) }()
	for t := range in {
		if n, er := mp(t); er != nil {
			err <- er
		} else {
			out <- n
		}
	}
}

func MapWithErrorSink[T any, N any](size int, mp MapWithErrorFunc[T, N], sink SinkFunc[error], in <-chan T) <-chan N {
	out := make(chan N, size)
	go mapWithErrorSinkWorker(mp, sink, in, out)
	return out
}

func mapWithErrorSinkWorker[T any, N any](mp MapWithErrorFunc[T, N], sink SinkFunc[error], in <-chan T, out chan<- N) {
	defer close(out)
	for t := range in {
		if n, er := mp(t); er != nil {
			sink(er)
		} else {
			out <- n
		}
	}
}
