package http

import (
	"context"
	"github.com/Zhoangp/Cart-Service/pb"
	"github.com/Zhoangp/Cart-Service/pkg/common"
)

type CartUsecase interface {
	GetCart(fakeId string) (*pb.GetCartResponse, error)
	AddToCart(cartId string, courseId string) error
}
type cartHandler struct {
	uc CartUsecase
	pb.UnimplementedCartServiceServer
}

func NewCartHandler(uc CartUsecase) *cartHandler {
	return &cartHandler{uc: uc}
}
func HandleError(err error) *pb.ErrorResponse {
	if errors, ok := err.(*common.AppError); ok {
		return &pb.ErrorResponse{
			Code:    int64(errors.StatusCode),
			Message: errors.Message,
		}
	}
	appErr := common.ErrInternal(err.(error))
	return &pb.ErrorResponse{
		Code:    int64(appErr.StatusCode),
		Message: appErr.Message,
	}
}
func (hdl cartHandler) GetCart(ctx context.Context, request *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	cart, err := hdl.uc.GetCart(request.Id)

	if err != nil {
		return &pb.GetCartResponse{
			Error: HandleError(err),
		}, nil
	}
	return cart, nil
}
func (hdl cartHandler) AddItem(ctx context.Context, request *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	if err := hdl.uc.AddToCart(request.CartId, request.CourseId); err != nil {
		return &pb.AddItemResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.AddItemResponse{}, nil
}
