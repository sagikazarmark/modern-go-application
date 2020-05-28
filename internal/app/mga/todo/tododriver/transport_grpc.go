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
		AddItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.AddItem,
			decodeAddItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeAddItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		ListItemsHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.ListItems,
			decodeListItemsGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeListItemsGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		DeleteItemsHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.DeleteItems,
			decodeDeleteItemsGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeDeleteItemsGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		GetItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.GetItem,
			decodeGetItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeGetItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		UpdateItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.UpdateItem,
			decodeUpdateItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeUpdateItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		DeleteItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.DeleteItem,
			decodeDeleteItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeDeleteItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeAddItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*todov1beta1.AddItemRequest)

	return AddItemRequest{
		NewItem: todo.NewItem{
			Title: req.GetTitle(),
			Order: int(req.GetOrder()),
		},
	}, nil
}

func encodeAddItemGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(AddItemResponse)

	return &todov1beta1.AddItemResponse{
		Item: marshalItemGRPC(resp.Item),
	}, nil
}

func decodeListItemsGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return ListItemsRequest{}, nil
}

func encodeListItemsGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ListItemsResponse)

	items := make([]*todov1beta1.TodoItem, 0, len(resp.Items))

	for _, item := range resp.Items {
		items = append(items, marshalItemGRPC(item))
	}

	return &todov1beta1.ListItemsResponse{
		Items: items,
	}, nil
}

func decodeDeleteItemsGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return DeleteItemsRequest{}, nil
}

func encodeDeleteItemsGRPCResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &todov1beta1.DeleteItemsResponse{}, nil
}

func decodeGetItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*todov1beta1.GetItemRequest)

	return GetItemRequest{
		Id: req.GetId(),
	}, nil
}

func encodeGetItemGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetItemResponse)

	return &todov1beta1.GetItemResponse{
		Item: marshalItemGRPC(resp.Item),
	}, nil
}

func decodeUpdateItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*todov1beta1.UpdateItemRequest)

	var (
		title     *string
		completed *bool
		order     *int
	)

	if req.Title != nil {
		title = &req.Title.Value
	}

	if req.Completed != nil {
		completed = &req.Completed.Value
	}

	if req.Order != nil {
		o := int(req.Order.Value)
		order = &o
	}

	return UpdateItemRequest{
		Id: req.GetId(),
		ItemUpdate: todo.ItemUpdate{
			Title:     title,
			Completed: completed,
			Order:     order,
		},
	}, nil
}

func encodeUpdateItemGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(UpdateItemResponse)

	return &todov1beta1.UpdateItemResponse{
		Item: marshalItemGRPC(resp.Item),
	}, nil
}

func decodeDeleteItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*todov1beta1.DeleteItemRequest)

	return DeleteItemRequest{
		Id: req.GetId(),
	}, nil
}

func encodeDeleteItemGRPCResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &todov1beta1.DeleteItemResponse{}, nil
}

func marshalItemGRPC(item todo.Item) *todov1beta1.TodoItem {
	return &todov1beta1.TodoItem{
		Id:        item.ID,
		Title:     item.Title,
		Completed: item.Completed,
		Order:     int32(item.Order),
	}
}
