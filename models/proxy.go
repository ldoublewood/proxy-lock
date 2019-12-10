package models

import "time"

type Proxy struct {
	Model
	Host 		    string	`gorm:"type:varchar(15)" json:"host"`
	Port            int		`gorm:"type:int;not null" json:"port"`
	Protocol        string	`gorm:"type:varchar(5);not null;default:'https'" json:"protocol,omitempty"`
	Owner           string	`gorm:"type:varchar(20);not null;" json:"owner,omitempty"`
	Status          string	`gorm:"type:varchar(10);not null;default:'free'" json:"status,omitempty"`
	DurationMinute  int64   `gorm:"type:integer;not null" json:"DurationMinute"`
	ExpiredAt       time.Time `json:"expired_at,omitempty"`

}
