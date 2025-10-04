package models

type Servers struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	ServerId        string `gorm:"unique"`
	Servername      string
	OwnerId         string
	Channels        []Channels        `gorm:"foreignKey:OwnServerId;references:ServerId;constraint:OnDelete:CASCADE;"`
	ServerMembers   []ServerMembers   `gorm:"foreignKey:ServerID;references:ServerId;constraint:OnDelete:CASCADE;"`
	InvitationCodes []InvitationCodes `gorm:"foreignKey:ServerID;references:ServerId;constraint:OnDelete:CASCADE;"`
}
