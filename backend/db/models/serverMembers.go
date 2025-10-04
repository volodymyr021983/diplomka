package models

type ServerMembers struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	UserID   string `gorm:"not null;index"`
	ServerID string `gorm:"not null;index"`
	UserRole string `gorm:"not null;default:owner"`
}
