package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tg_pay_gate/internal/utils/functions"
	_type "tg_pay_gate/internal/utils/type"
)

type User struct {
	ID         uuid.UUID            `gorm:"type:uuid;primary_key" json:"id"`
	Status     _type.UserStatusType `gorm:"index;default:1" json:"status"`
	TgID       int64                `gorm:"unique;index;not null" json:"tg_id"`
	CreateTime int64                `gorm:"autoCreateTime" json:"create_time"`

	//Orders []Order `gorm:"constraint:OnDelete:CASCADE;"` // 用户删除订单也删除
}

func (u *User) TableName() string {
	return "user"

}

func (u *User) DefaultOrder() string {
	return "create_time DESC"
}
func (u *User) ToDict() map[string]interface{} {
	return functions.StructToMap(u)
}
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	//u.Secret = functions.GenerateRandomString(24)

	return
}
