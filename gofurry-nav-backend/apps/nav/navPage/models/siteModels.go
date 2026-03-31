package models

import (
	"time"

	cm "github.com/GoFurry/gofurry-nav-backend/common/models"
)

const TableNameGfnSite = "gfn_site"

// TableName GfnSite's table name
func (*GfnSite) TableName() string {
	return TableNameGfnSite
}

// GfnSite mapped from table <gfn_site>
type GfnSite struct {
	ID         int64        `gorm:"column:id;type:bigint;primaryKey;comment:站点表id" json:"id"`                                           // 站点表id
	Name       string       `gorm:"column:name;type:character varying(255);not null;comment:站点名称" json:"name"`                          // 站点名称
	NameEn     string       `gorm:"column:name_en;type:character varying(255);not null;comment:站点名称-英文" json:"nameEn"`                  // 站点名称-英文
	Domain     string       `gorm:"column:domain;type:json;not null;comment:站点域名" json:"domain"`                                        // 站点域名
	Info       string       `gorm:"column:info;type:text;not null;comment:站点描述" json:"info"`                                            // 站点描述
	InfoEn     string       `gorm:"column:info_en;type:text;not null;comment:站点描述-英文" json:"infoEn"`                                    // 站点描述-英文
	CreateTime cm.LocalTime `gorm:"column:create_time;type:int;type:unsigned;not null;autoCreateTime;comment:创建时间" json:"createTime"`   // 创建时间
	UpdateTime cm.LocalTime `gorm:"column:update_time;type:int;type:unsigned;not null;autoUpdateTime;comment:修改时间" json:"updateTime"`   // 修改时间
	Country    *string      `gorm:"column:country;type:character varying(20);comment:站点所属国家" json:"country"`                            // 站点所属国家
	Nsfw       string       `gorm:"column:nsfw;type:character varying(4);default:''::character varying;comment:是否NSFW 1 0" json:"nsfw"` // 是否NSFW 1 0
	Welfare    string       `gorm:"column:welfare;type:character varying(4);comment:是否公益项目 1 0" json:"welfare"`                         // 是否公益项目 1 0
	Icon       *string      `gorm:"column:icon;type:character varying(255);comment:站点图标" json:"icon"`                                   // 站点图标
	Deleted    bool         `gorm:"column:deleted;type:boolean;comment:软删除" json:"deleted"`
}

const TableNameGfnSiteGroup = "gfn_site_group"

// TableName GfnSiteGroup's table name
func (*GfnSiteGroup) TableName() string {
	return TableNameGfnSiteGroup
}

// GfnSiteGroup mapped from table <gfn_site_group>
type GfnSiteGroup struct {
	ID         int64     `gorm:"column:id;type:bigint;primaryKey;comment:分组表id" json:"id"`                                         // 分组表id
	Name       string    `gorm:"column:name;type:character varying(255);not null;comment:分组名称" json:"name"`                        // 分组名称
	NameEn     string    `gorm:"column:name_en;type:character varying(255);not null;comment:分组名称-英文" json:"nameEn"`                // 分组名称-英文
	Info       string    `gorm:"column:info;type:character varying(255);not null;comment:分组简介" json:"info"`                        // 分组简介
	InfoEn     string    `gorm:"column:info_en;type:character varying(255);not null;comment:分组简介-英文" json:"infoEn"`                // 分组简介-英文
	Priority   int64     `gorm:"column:priority;type:bigint;not null;comment:分组优先级" json:"priority"`                               // 分组优先级
	CreateTime time.Time `gorm:"column:create_time;type:int;type:unsigned;not null;autoCreateTime;comment:创建时间" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;type:int;type:unsigned;not null;autoUpdateTime;comment:修改时间" json:"updateTime"` // 修改时间
}

const TableNameGfnSiteGroupMap = "gfn_site_group_map"

// TableName GfnSiteGroupMap's table name
func (*GfnSiteGroupMap) TableName() string {
	return TableNameGfnSiteGroupMap
}

// GfnSiteGroupMap mapped from table <gfn_site_group_map>
type GfnSiteGroupMap struct {
	ID         int64     `gorm:"column:id;type:bigint;primaryKey;comment:分组映射表id" json:"id"`                                       // 分组映射表id
	SiteID     int64     `gorm:"column:site_id;type:bigint;not null;comment:站点id" json:"siteId,string"`                            // 站点id
	GroupID    int64     `gorm:"column:group_id;type:bigint;not null;comment:分组id" json:"groupId,string"`                          // 分组id
	CreateTime time.Time `gorm:"column:create_time;type:int;type:unsigned;not null;autoCreateTime;comment:创建时间" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;type:int;type:unsigned;not null;autoUpdateTime;comment:修改时间" json:"updateTime"` // 修改时间
}
