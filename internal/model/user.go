package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"type:varchar(100);not null"` // 密码不返回给前端
	Nickname  string         `json:"nickname" gorm:"type:varchar(50)"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Role      string         `json:"role" gorm:"type:varchar(20);default:'user'"` // 角色：admin, user
	Status    int            `json:"status" gorm:"default:1"`                     // 状态：0-禁用，1-启用
	LastLogin *time.Time     `json:"last_login"`                                  // 最后登录时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 可以在这里进行密码加密等操作
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// 可以在这里进行密码加密等操作
	return nil
}

// UserProfile 用户资料
type UserProfile struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	RealName  string         `json:"real_name" gorm:"type:varchar(50)"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	Address   string         `json:"address" gorm:"type:varchar(255)"`
	Bio       string         `json:"bio" gorm:"type:text"`
	Birthday  *time.Time     `json:"birthday"`
	Gender    string         `json:"gender" gorm:"type:varchar(10)"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}
