package runner

// Group collects actors (functions) and runs them concurrently.
// When one actor (function) returns, all actors are interrupted.
type Group interface {
	// Add an actor (function) to the group.
	Add(execute func() error, interrupt func(error))
}

// Runner is an actor.
type Runner interface {
	Start() error
	Stop(err error)
}

// Register actors in a group.
func Register(group Group, actors ...Runner) {
	for _, actor := range actors {
		group.Add(actor.Start, actor.Stop)
	}
}
