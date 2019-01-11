package runner

import "testing"

type runCloser struct {
	started bool
	closed  bool
}

func (r *runCloser) Run() error {
	r.started = true

	return nil
}

func (r *runCloser) Close() error {
	r.closed = true

	return nil
}

func TestRunCloserRunner(t *testing.T) {
	rc := &runCloser{}
	runner := NewRunCloserRunner(rc)

	err := runner.Start()
	if err != nil {
		t.Fatal(err)
	}

	if !rc.started {
		t.Error("run closer is expected to be started")
	}

	runner.Stop(nil)

	if !rc.closed {
		t.Error("run closer is expected to be closed")
	}
}
