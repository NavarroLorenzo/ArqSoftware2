package domain

import "time"

type Role string

const (
	RoleNormal Role = "normal"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex;size:191" json:"email"`
	Password  string    `json:"-"` // hash
	Role      Role      `gorm:"type:enum('normal','admin');default:'normal'" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
