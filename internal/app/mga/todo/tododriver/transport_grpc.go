package tododriver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) todov1beta1.TodoListServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return todov1beta1.TodoListKitServer{
		CreateTodoHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.AddItem,
			decodeCreateTodoGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeCreateTodoGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		ListTodosHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.ListItems,
			decodeListTodosGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeListTodosGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		MarkAsCompleteHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.MarkAsComplete,
			decodeMarkAsCompleteGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeMarkAsCompleteGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeCreateTodoGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*todov1beta1.CreateTodoRequest)

	return AddItemRequest{
		NewItem: todo.NewItem{Title: req.GetTitle()},
	}, nil
}

func encodeCreateTodoGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(AddItemResponse)

	return &todov1beta1.CreateTodoResponse{
		Id: resp.Item.ID,
	}, nil
}

func decodeListTodosGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func encodeListTodosGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ListItemsResponse)

	grpcResp := &todov1beta1.ListTodosResponse{
		Todos: make([]*todov1beta1.Todo, len(resp.Items)),
	}

	for i, t := range resp.Items {
		grpcResp.Todos[i] = &todov1beta1.Todo{
			Id:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
		}
	}

	return grpcResp, nil
}

func decodeMarkAsCompleteGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*todov1beta1.MarkAsCompleteRequest)

	return MarkAsCompleteRequest{
		Id: req.GetId(),
	}, nil
}

func encodeMarkAsCompleteGRPCResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &todov1beta1.MarkAsCompleteResponse{}, nil
}
