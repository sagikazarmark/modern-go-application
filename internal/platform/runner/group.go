// Package runner provides a set of actors in the form defined by https://github.com/oklog/run.
package runner

// Group collects actors (functions) and runs them concurrently.
// When one actor (function) returns, all actors are interrupted.
type Group interface {
	// Add an actor (function) to the group.
	Add(execute func() error, interrupt func(error))
}

// Runner is an actor that can be added to a run group.
type Runner interface {
	// Start is the main logic for the actor.
	Start() error

	// Stop is the interrupt function for the actor.
	Stop(err error)
}

// Register actors in a group.
func Register(group Group, actors ...Runner) {
	for _, actor := range actors {
		group.Add(actor.Start, actor.Stop)
	}
}
