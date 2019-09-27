package tododriver

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

type grpcServer struct {
	*todov1beta1.UnimplementedTodoListServer

	createTodo kitgrpc.Handler
	listTodos  kitgrpc.Handler
	markAsDone kitgrpc.Handler
}

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) todov1beta1.TodoListServer {
	return &grpcServer{
		createTodo: kitgrpc.NewServer(
			endpoints.CreateTodo,
			decodeCreateTodoGRPCRequest,
			encodeCreateTodoGRPCResponse,
			options...,
		),
		listTodos: kitgrpc.NewServer(
			endpoints.ListTodos,
			decodeListTodosGRPCRequest,
			encodeListTodosGRPCResponse,
			options...,
		),
		markAsDone: kitgrpc.NewServer(
			endpoints.MarkAsDone,
			decodeMarkAsDoneGRPCRequest,
			encodeMarkAsDoneGRPCResponse,
			options...,
		),
	}
}

func (s *grpcServer) CreateTodo(
	ctx context.Context,
	req *todov1beta1.CreateTodoRequest,
) (*todov1beta1.CreateTodoResponse, error) {
	_, rep, err := s.createTodo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*todov1beta1.CreateTodoResponse), nil
}

func decodeCreateTodoGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*todov1beta1.CreateTodoRequest)

	return createTodoRequest{
		Text: req.GetText(),
	}, nil
}

func encodeCreateTodoGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(createTodoResponse)

	return &todov1beta1.CreateTodoResponse{
		Id: resp.ID,
	}, nil
}

func (s *grpcServer) ListTodos(
	ctx context.Context,
	req *todov1beta1.ListTodosRequest,
) (*todov1beta1.ListTodosResponse, error) {
	_, rep, err := s.listTodos.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*todov1beta1.ListTodosResponse), nil
}

func decodeListTodosGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func encodeListTodosGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(listTodosResponse)

	grpcResp := &todov1beta1.ListTodosResponse{
		Todos: make([]*todov1beta1.Todo, len(resp.Todos)),
	}

	for i, t := range resp.Todos {
		grpcResp.Todos[i] = &todov1beta1.Todo{
			Id:   t.ID,
			Text: t.Text,
			Done: t.Done,
		}
	}

	return grpcResp, nil
}

func (s *grpcServer) MarkAsDone(
	ctx context.Context,
	req *todov1beta1.MarkAsDoneRequest,
) (*todov1beta1.MarkAsDoneResponse, error) {
	_, rep, err := s.markAsDone.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*todov1beta1.MarkAsDoneResponse), nil
}

func decodeMarkAsDoneGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*todov1beta1.MarkAsDoneRequest)

	return markAsDoneRequest{
		ID: req.GetId(),
	}, nil
}

func encodeMarkAsDoneGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		err := f.Failed()

		return nil, status.Error(getErrorCode(err), err.Error())
	}

	return &todov1beta1.MarkAsDoneResponse{}, nil
}

func getErrorCode(err error) codes.Code {
	code := codes.Internal

	// nolint: gocritic
	switch errors.Cause(err).(type) {
	case todo.NotFoundError:
		code = codes.NotFound
	}

	return code
}
