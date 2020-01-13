package tododriver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
)

type grpcServer struct {
	*todov1beta1.UnimplementedTodoListServer

	errorEncoder kitxgrpc.EncodeErrorResponseFunc

	createTodo kitgrpc.Handler
	listTodos  kitgrpc.Handler
	markAsDone kitgrpc.Handler
}

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) todov1beta1.TodoListServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return &grpcServer{
		errorEncoder: errorEncoder,

		createTodo: kitgrpc.NewServer(
			endpoints.CreateTodo,
			decodeCreateTodoGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeCreateTodoGRPCResponse, errorEncoder),
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
			kitxgrpc.ErrorResponseEncoder(encodeMarkAsDoneGRPCResponse, errorEncoder),
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
		return nil, s.errorEncoder(ctx, err)
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
		return nil, s.errorEncoder(ctx, err)
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
		return nil, s.errorEncoder(ctx, err)
	}
	return rep.(*todov1beta1.MarkAsDoneResponse), nil
}

func decodeMarkAsDoneGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*todov1beta1.MarkAsDoneRequest)

	return markAsDoneRequest{
		ID: req.GetId(),
	}, nil
}

func encodeMarkAsDoneGRPCResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &todov1beta1.MarkAsDoneResponse{}, nil
}
