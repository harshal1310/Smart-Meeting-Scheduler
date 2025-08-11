package model

type User struct {
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserCode string  `gorm:"unique;not null" json:"userCode"` // For custom user IDs like "user1"
	Name     string  `json:"name"`
	Events   []Event `gorm:"foreignKey:UserID;references:UserCode"`
}
