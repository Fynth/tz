package models

type Subscription struct {
	ID          int    `gorm:"column:id;primaryKey"`
	ServiceName string `gorm:"column:service_name" json:"service_name"`
	Price       int    `gorm:"column:price" json:"price"`
	UserId      string `gorm:"column:user_id" json:"user_id"`
	StartDate   string `gorm:"column:start_date;omitempty" json:"start_date"`
}

func (s *Subscription) TableName() string {
	return "subscriptions"
}
