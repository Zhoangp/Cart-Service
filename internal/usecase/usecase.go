package usecase

import (
	"errors"
	"github.com/Zhoangp/Cart-Service/config"
	"github.com/Zhoangp/Cart-Service/internal/model"
	"github.com/Zhoangp/Cart-Service/pb"
	"github.com/Zhoangp/Cart-Service/pkg/common"
	"github.com/Zhoangp/Cart-Service/pkg/utils"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type CartRepo interface {
	FindDataWithCondition(condition map[string]any) (*model.Cart, error)
	AddItemToCart(item *model.CartCourse) error
	FindCartWithUser(userId int) (int, error)
}
type cartUsecase struct {
	config *config.Config
	repo   CartRepo
	h      *utils.Hasher
}

func NewCartUseCase(cf *config.Config, repo CartRepo, h *utils.Hasher) *cartUsecase {
	return &cartUsecase{cf, repo, h}
}
func (uc cartUsecase) GetCart(fakeId string) (*pb.GetCartResponse, error) {
	id, err := uc.h.Decode(fakeId)

	if err != nil {
		return nil, common.ErrInternal(err)
	}

	cart, err := uc.repo.FindDataWithCondition(map[string]any{"user_id": id})
	if err != nil {
		return nil, err
	}

	res := pb.GetCartResponse{
		Cart: &pb.Cart{
			Id: uc.h.Encode(cart.Id),
		},
	}
	var totalPrice float64
	for _, item := range cart.Courses {
		totalPrice += item.CoursePrice
		res.Cart.Courses = append(res.Cart.Courses, &pb.Course{
			Id:          uc.h.Encode(item.Id),
			Title:       item.Title,
			Description: item.CourseDescription,
			Level:       item.CourseLevel,
			Price:       strconv.FormatFloat(item.CoursePrice, 'f', -1, 64),
			Discount:    item.CourseDiscount,
			Currency:    item.CourseCurrency,
			Duration:    item.CourseDuration,
			Status:      item.CourseStatus,
			Rating:      item.CourseRating,
			Thumbnail: &pb.Image{
				Url:    item.CourseThumbnail.Url,
				Width:  item.CourseThumbnail.Width,
				Height: item.CourseThumbnail.Height,
			},
		})
	}
	res.Cart.Currency = cart.Courses[0].CourseCurrency
	res.Cart.TotalPrice = strconv.FormatFloat(math.Round(totalPrice*100)/100, 'f', -1, 64)
	return &res, nil
}

func (uc cartUsecase) AddToCart(cartId string, courseId string) error {
	newCartId, err := uc.h.Decode(cartId)

	if err != nil {
		return common.ErrInternal(err)
	}
	newCourseId, err := uc.h.Decode(courseId)
	if err != nil {
		return common.ErrInternal(err)
	}
	item := &model.CartCourse{
		CartId:   newCartId,
		CourseId: newCourseId,
	}
	if err = uc.repo.AddItemToCart(item); err != nil {
		if err == gorm.ErrDuplicatedKey {
			return common.NewCustomError(errors.New("course has already been added"), "Course has already been added!")
		}
	}
	return nil
}
