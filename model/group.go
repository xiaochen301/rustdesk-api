package model

const (
	GroupTypeDefault = 1 // 默认
	GroupTypeShare   = 2 // 共享
)

// XC: 群组管理模式
const (
	GroupModeCentralized = 1 // 集中式：群组管理员看全组，普通成员只看自己
	GroupModeEqual       = 2 // 平权式：全员互看（共享通讯录）
)

type Group struct {
	IdModel
	Name string `json:"name" gorm:"default:'';not null;"`
	Type int    `json:"type" gorm:"default:1;not null;"`
	Mode int    `json:"mode" gorm:"default:1;not null;"` // XC: 1=集中式 2=平权式
	TimeModel

	// XC: 非持久化——群组管理员列表（查询后填充）
	AdminIds []uint `json:"admin_ids" gorm:"-"`
}

// XC: 群组管理员（多对多关联）
type GroupAdmin struct {
	IdModel
	GroupId uint `json:"group_id" gorm:"not null;index"`
	UserId  uint `json:"user_id" gorm:"not null;index"`
	TimeModel
}

type GroupList struct {
	Groups []*Group `json:"list"`
	Pagination
}

type DeviceGroup struct {
	IdModel
	Name string `json:"name" gorm:"default:'';not null;"`
	TimeModel
}

type DeviceGroupList struct {
	DeviceGroups []*DeviceGroup `json:"list"`
	Pagination
}
