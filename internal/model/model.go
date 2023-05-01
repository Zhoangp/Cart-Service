package model

import "github.com/Zhoangp/Cart-Service/pkg/common"

type Cart struct {
	common.SQLModel
	UserId   int      `gorm:"column:user_id"`
	CourseId int      `gorm:"column:course_id"`
	Courses  []Course `gorm:"many2many:Cart_Course;"`
}
type CartCourse struct {
	CartId   int `gorm:"column:cart_id"`
	CourseId int `gorm:"column:course_id"`
}
type Course struct {
	common.SQLModel
	Title             string        `json:"title" gorm:"column:title"`
	CourseDescription string        `json:"description" gorm:"column:description"`
	CourseLevel       string        `json:"level"  gorm:"column:level"`
	CourseLanguage    string        `json:"language" gorm:"column:language"`
	CoursePrice       float64       `json:"price" gorm:"column:price"`
	CourseCurrency    string        `json:"currency" gorm:"column:currency"`
	CourseDiscount    float32       `json:"discount" gorm:"column:discount"`
	CourseDuration    string        `json:"duration" gorm:"column:duration"`
	CourseStatus      string        `json:"status" gorm:"column:status"`
	CourseRating      float32       `json:"rating" gorm:"column:rating"`
	InstructorID      int           `json:"instructor_id" gorm:"column:instructor_id"`
	CourseThumbnail   *common.Image `json:"thumbnail" gorm:"column:thumbnail"`
}

func (CartCourse) TableName() string {
	return "cart_courses"
}
func (Course) TableName() string {
	return "Courses"
}

func (Cart) TableName() string {
	return "Cart"
}
