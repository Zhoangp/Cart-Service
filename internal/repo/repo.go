package repo

import (
	"github.com/Zhoangp/Cart-Service/internal/model"
	"gorm.io/gorm"
)

type cartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) *cartRepo {
	return &cartRepo{db: db}
}
func (rp cartRepo) FindDataWithCondition(condition map[string]any) (*model.Cart, error) {
	var cart model.Cart
	err := rp.db.Model(cart).Where(condition).Preload("Courses").Find(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}
func (rp cartRepo) FindCartWithUser(userId int) (int, error) {
	var res model.Cart
	if err := rp.db.Table(model.Cart{}.TableName()).Where("user_id = ?", userId).First(&res).Error; err != nil {
		return -1, err
	}
	return res.Id, nil
}
func (rp cartRepo) AddItemToCart(item *model.CartCourse) error {
	db := rp.db.Begin()
	if err := db.Table(model.CartCourse{}.TableName()).Create(item).Error; err != nil {
		db.Rollback()

		return err
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return err
	}
	return nil

}
