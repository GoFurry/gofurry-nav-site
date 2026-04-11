package audit

import "time"

type AdminAuditLog struct {
	ID             int64     `gorm:"column:id;type:bigserial;primaryKey;autoIncrement" json:"id"`
	Action         string    `gorm:"column:action;type:varchar(64);not null" json:"action"`
	Resource       string    `gorm:"column:resource;type:varchar(128);not null" json:"resource"`
	TargetID       string    `gorm:"column:target_id;type:varchar(128)" json:"target_id"`
	Operator       string    `gorm:"column:operator;type:varchar(64);not null" json:"operator"`
	SessionVersion int64     `gorm:"column:session_version;type:bigint;not null;default:0" json:"session_version"`
	RequestID      string    `gorm:"column:request_id;type:varchar(128)" json:"request_id"`
	IPAddress      string    `gorm:"column:ip_address;type:varchar(64)" json:"ip_address"`
	UserAgent      string    `gorm:"column:user_agent;type:text" json:"user_agent"`
	BeforeData     string    `gorm:"column:before_data;type:text" json:"before_data"`
	AfterData      string    `gorm:"column:after_data;type:text" json:"after_data"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp(0) without time zone;not null;autoCreateTime" json:"created_at"`
}

func (AdminAuditLog) TableName() string {
	return "gfa_admin_audit_log"
}

type Meta struct {
	Operator       string
	SessionVersion int64
	RequestID      string
	IPAddress      string
	UserAgent      string
}
