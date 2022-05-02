package pipes

func Router[T any, N comparable](size int, matches []N, compare func(T) N, in <-chan T) ([]ChanPull[T], ChanPull[T]) {
	orElse := make(chan T, size)
	outs := make([]ChanPull[T], len(matches))
	routes := make(map[N]chan<- T, len(matches))
	for i, match := range matches {
		out := make(chan T, size)
		outs[i] = out
		routes[match] = out
	}
	go routerWorker(compare, in, routes, orElse)
	return outs, orElse
}

func routerWorker[T any, N comparable](compare func(T) N, in <-chan T, routes map[N]chan<- T, orElse chan<- T) {
	defer func() {
		for _, out := range routes {
			close(out)
		}
		close(orElse)
	}()
	for t := range in {
		if route, exists := routes[compare(t)]; exists {
			route <- t
		} else {
			orElse <- t
		}
	}
}

func RouterWithSink[T any, N comparable](size int, matches []N, compare func(T) N, sink func(T), in <-chan T) []ChanPull[T] {
	outs := make([]ChanPull[T], len(matches))
	routes := make(map[N]chan<- T, len(matches))
	for i, match := range matches {
		out := make(chan T, size)
		outs[i] = out
		routes[match] = out
	}
	go routerWithSinkWorker(compare, in, routes, sink)
	return outs
}

func routerWithSinkWorker[T any, N comparable](compare func(T) N, in <-chan T, routes map[N]chan<- T, sink func(T)) {
	defer func() {
		for _, out := range routes {
			close(out)
		}
	}()
	for t := range in {
		if route, exists := routes[compare(t)]; exists {
			route <- t
		} else {
			sink(t)
		}
	}
}
