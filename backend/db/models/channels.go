package models

type Channels struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	ChannelId   string `gorm:"unique"`
	OwnServerId string
	Channelname string
	ChannelType string
}
