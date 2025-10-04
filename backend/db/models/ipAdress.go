package models

type IpAddress struct {
	ID              uint `gorm:"primaryKey;autoIncrement"`
	IpAddress       string
	EmailVerifyTime *int64
	ResetPwdTime    *int64
}
