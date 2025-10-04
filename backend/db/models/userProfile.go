package models

type UserProfile struct {
	ID                uint            `gorm:"primaryKey;autoIncrement"`
	UserID            string          `gorm:"uniqueIndex"`
	Username          string          `gorm:"unique"`
	Email             string          `gorm:"unique"`
	Servers           []Servers       `gorm:"foreignKey:OwnerId;references:UserID;constraint:OnDelete:CASCADE;"`
	ServerMemberships []ServerMembers `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE;"`
}
