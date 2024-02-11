package models

import "time"

const (
	AccountKey = "account/"
)

type Account struct {
	Id          string    `json:"id" gorm:"default:uuid_generate_v4()" db:"id"`
	Type        string    `json:"type" db:"type"`
	IsDisabled  bool      `json:"isDisabled" db:"is_disabled"`
	CreatedTime time.Time `json:"createdTime" gorm:"default:now()" db:"created_time"`
}
