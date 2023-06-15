package http

import (
	"context"
	"github.com/Zhoangp/Cart-Service/config"
	"github.com/Zhoangp/Cart-Service/pb"
	course2 "github.com/Zhoangp/Cart-Service/pb/course"
	error1 "github.com/Zhoangp/Cart-Service/pb/error"
	"github.com/Zhoangp/Cart-Service/pkg/client"
	"github.com/Zhoangp/Cart-Service/pkg/common"
)

type CartUsecase interface {
	GetCart(fakeId string) (*pb.GetCartResponse, error)
	AddToCart(cartId string, courseId string) error
	RemoveItem(cartId string, courseId string) error
	ResetCart(cartId string) error
	NewCart(userId string) error
}
type cartHandler struct {
	uc CartUsecase
	pb.UnimplementedCartServiceServer
	cf *config.Config
}

func NewCartHandler(uc CartUsecase, cf *config.Config) *cartHandler {
	return &cartHandler{uc: uc, cf: cf}
}
func HandleError(err error) *error1.ErrorResponse {
	if errors, ok := err.(*common.AppError); ok {
		return &error1.ErrorResponse{
			Code:    int64(errors.StatusCode),
			Message: errors.Message,
		}
	}
	appErr := common.ErrInternal(err.(error))
	return &error1.ErrorResponse{
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
func (hdl cartHandler) AddItem(ctx context.Context, request *pb.CartItemRequest) (*pb.CartItemResponse, error) {
	courseService, err := client.InitCourseServiceClient(hdl.cf)
	if err != nil {
		return &pb.CartItemResponse{
			Error: HandleError(err),
		}, nil
	}
	course, err := courseService.GetCourse(ctx, &course2.GetCourseRequest{
		Id: request.CourseId,
	})
	if course.Course.Instructor.UserId == request.CartId {
		return &pb.CartItemResponse{
			Error: HandleError(common.NewCustomError(err, "Your are an instructor of this course!")),
		}, nil
	}
	enrollments, err := courseService.GetEnrollments(ctx, &course2.GetEnrollmentsRequest{
		UserId: request.CartId,
	})
	for _, item := range enrollments.Enrollments {
		if request.CourseId == item.CourseId {
			return &pb.CartItemResponse{
				Error: HandleError(common.NewCustomError(err, "You have been in this course!")),
			}, nil
		}
	}
	if err := hdl.uc.AddToCart(request.CartId, request.CourseId); err != nil {
		return &pb.CartItemResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.CartItemResponse{}, nil
}
func (hdl cartHandler) RemoveItem(ctx context.Context, request *pb.CartItemRequest) (*pb.CartItemResponse, error) {
	if err := hdl.uc.RemoveItem(request.CartId, request.CourseId); err != nil {
		return &pb.CartItemResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.CartItemResponse{}, nil
}
func (hdl cartHandler) ResetCart(ctx context.Context, request *pb.ResetCartRequest) (*pb.ResetCartResponse, error) {
	if err := hdl.uc.ResetCart(request.CartId); err != nil {
		return &pb.ResetCartResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.ResetCartResponse{}, nil
}
func (hdl cartHandler) CreateCart(ctx context.Context, request *pb.CreateCartRequest) (*pb.CreateCartResponse, error) {

	if err := hdl.uc.NewCart(request.UserId); err != nil {
		return &pb.CreateCartResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.CreateCartResponse{}, nil
}
