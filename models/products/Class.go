package products

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Active    int8      `gorm:"default:1;not null;index" json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy int32     `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy int32     `json:"updatedBy"`
}

func (handle *Class) AllClasss(db *gorm.DB) (*[]Class, error) {
	var err error
	classs := []Class{}
	err = db.Debug().Model(&Class{}).Where("active = ?", 1).Limit(100).Find(&classs).Error
	if err != nil {
		return &[]Class{}, err
	}

	return &classs, nil
}

func (handle *Class) PrepareClass() {
	handle.ID = 0
	handle.Name = html.EscapeString(strings.TrimSpace(handle.Name))

	handle.CreatedAt = time.Now()
	handle.UpdatedAt = time.Now()
}

func (handle *Class) ValidateClass() error {

	if handle.Name == "" {
		return errors.New("name required")
	}
	return nil
}

func (handle *Class) CreateClass(db *gorm.DB) (*Class, error) {
	var err = db.Debug().Model(&Class{}).Create(&handle).Error
	if err != nil {
		return &Class{}, err
	}
	return handle, nil
}

func (handle *Class) ClassByName(db *gorm.DB, className string) (*Class, error) {
	var err = db.Debug().Model(&Class{}).Where("name = ?", className).Take(&handle).Error
	if err != nil {
		return &Class{}, err
	}
	return handle, nil
}

func (handle *Class) UpdateClass(db *gorm.DB, id uint32) (*Class, error) {
	var err = db.Debug().Model(&Class{}).Where("id = ?", id).Updates(Class{Name: handle.Name, Active: handle.Active}).Error
	if err != nil {
		return &Class{}, err
	}
	return handle, nil
}

func (handle *Class) DeleteClass(db *gorm.DB, id uint32) (int64, error) {

	db = db.Debug().Model(&Class{}).Where("id = ?", id).Take(&Class{}).Delete(&Class{})

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("class not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
