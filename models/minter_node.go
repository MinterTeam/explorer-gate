package models

import (
	"strconv"
	"time"
)

type MinterNode struct {
	ID         uint    `gorm:"primary_key"`
	Host       string  `json:"host"        gorm:"type:varchar(255)"`
	Port       uint    `json:"port"        gorm:"type:int; default:8841"`
	Version    string  `json:"version"     gorm:"type:varchar(255)"`
	Ping       float32 `json:"ping"        gorm:"type:numeric(7,3)"`
	IsSecure   bool    `json:"is_secure"   gorm:"default:false"`
	IsActive   bool    `json:"is_active"   gorm:"default:false"`
	IsLocal    bool    `json:"is_local"    gorm:"default:false"`
	IsExcluded bool    `json:"is_excluded" gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func (node MinterNode) GetFullLink() string {
	protocol := `http`

	if node.IsSecure {
		protocol += `s://`
	} else {
		protocol += `://`
	}

	return protocol + node.Host + `:` + strconv.Itoa(int(node.Port))
}
