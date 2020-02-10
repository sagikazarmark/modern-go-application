package tododriver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) todov1beta1.TodoListServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return todov1beta1.TodoListKitServer{
		CreateTodoHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.CreateTodo,
			decodeCreateTodoGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeCreateTodoGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		ListTodosHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.ListTodos,
			decodeListTodosGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeListTodosGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		MarkAsDoneHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.MarkAsDone,
			decodeMarkAsDoneGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeMarkAsDoneGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeCreateTodoGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*todov1beta1.CreateTodoRequest)

	return CreateTodoRequest{
		Text: req.GetText(),
	}, nil
}

func encodeCreateTodoGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(CreateTodoResponse)

	return &todov1beta1.CreateTodoResponse{
		Id: resp.Id,
	}, nil
}

func decodeListTodosGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func encodeListTodosGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ListTodosResponse)

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

func decodeMarkAsDoneGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*todov1beta1.MarkAsDoneRequest)

	return MarkAsDoneRequest{
		Id: req.GetId(),
	}, nil
}

func encodeMarkAsDoneGRPCResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &todov1beta1.MarkAsDoneResponse{}, nil
}
