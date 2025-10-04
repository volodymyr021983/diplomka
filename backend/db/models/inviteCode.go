package models

type InvitationCodes struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ServerID  string `gorm:"not null;index"`
	Token     string `gorm:"size:512;not null;unique"`
	Status    string `gorm:"size:50;not null;default:pending"`
	ExpiresAt int64  `gorm:"not null"`
	CreatedAt int64  `gorm:"not null;"`
	UsedAt    *int64
}
