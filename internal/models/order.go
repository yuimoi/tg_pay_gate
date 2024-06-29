package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"tg_pay_gate/internal/utils/functions"
	_type "tg_pay_gate/internal/utils/type"
)

// 使用指针可以方便的置空，使用原则：必须要判断是否为空的情况
type Order struct {
	ID         uuid.UUID             `gorm:"type:uuid;primary_key;not null" json:"id"`
	Status     _type.OrderStatusType `gorm:"default:0;not null" json:"status"` // 0待支付 100已支付 -1超时 -2强行关闭
	CreateTime int64                 `gorm:"index;autoCreateTime;not null" json:"create_time"`
	//EndTime    int64                 `gorm:"index" json:"end_time"` // 结束时间,订单完成,则标记为支付时间

	Price decimal.Decimal `gorm:"not null" json:"price"`

	TgID int64 `gorm:"index;not null" json:"tg_id"` // 不能给unique，一个tg_id会创建多个订单

	//UserID uuid.UUID `json:"user_id"`
	//User   User      `gorm:"foreignKey:UserID"`
}

func (*Order) TableName() string {
	return "order"
}
func (*Order) DefaultOrder() string {
	return "create_time DESC"
}
func (o *Order) ToDict() map[string]interface{} {
	return functions.StructToMap(o)
}

func (o Order) QueryLimitUser(query *gorm.DB, currentUser *User) *gorm.DB {
	query = query.Where("user_id=?", currentUser.ID)
	return query
}
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
func NewOrder(price decimal.Decimal, tgID int64) *Order {
	order := &Order{
		ID:    uuid.New(),
		Price: price,
		TgID:  tgID,
	}
	return order
}
