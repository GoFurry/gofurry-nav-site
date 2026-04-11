package models

import "time"

type AdminAccount struct {
	ID                int64      `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	PasswordHash      string     `gorm:"column:password_hash;type:text;not null" json:"password_hash"`
	SessionVersion    int64      `gorm:"column:session_version;type:bigint;not null" json:"session_version"`
	CreatedAt         time.Time  `gorm:"column:created_at;type:timestamp(0) without time zone;not null;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;type:timestamp(0) without time zone;not null;autoUpdateTime" json:"updated_at"`
	PasswordUpdatedAt *time.Time `gorm:"column:password_updated_at;type:timestamp(0) without time zone" json:"password_updated_at"`
}

func (AdminAccount) TableName() string {
	return "gfa_admin_account"
}
