package models

type Users struct {
	Base
	Name     string `form:"name" json:"name,omitempty" bindind:"required"`
	Password string `form:"password" json:"hashed_password" bindind:"required"`
	Email    string `gorm:"type:varchar(110);unique_index" form:"email" json:"email,omitempty" binding:"required"`
	Location string `form:"location" json:"location,omitempty"`
	Avatar   string `form:"avatar" json:"avatar,omitempty"`
}
