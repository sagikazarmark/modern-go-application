package greetingdriver

import (
	"context"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

type greeterStub struct {
	resp *greeting.HelloResponse
	err  error
}

func (s *greeterStub) SayHello(ctx context.Context, req greeting.HelloRequest) (*greeting.HelloResponse, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.resp, nil
}
