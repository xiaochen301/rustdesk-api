package admin

import "github.com/lejianwen/rustdesk-api/v2/model"

type GroupForm struct {
	Id       uint   `json:"id"`
	Name     string `json:"name" validate:"required"`
	Type     int    `json:"type"`
	Mode     int    `json:"mode"`      // XC: 1=集中式 2=平权式
	AdminIds []uint `json:"admin_ids"` // XC: 群组管理员用户ID列表
}

func (gf *GroupForm) FromGroup(group *model.Group) *GroupForm {
	gf.Id = group.Id
	gf.Name = group.Name
	gf.Type = group.Type
	gf.Mode = group.Mode
	gf.AdminIds = group.AdminIds
	return gf
}

func (gf *GroupForm) ToGroup() *model.Group {
	group := &model.Group{}
	group.Id = gf.Id
	group.Name = gf.Name
	group.Type = gf.Type
	group.Mode = gf.Mode
	if group.Mode == 0 {
		group.Mode = model.GroupModeCentralized // XC: 默认集中式
	}
	return group
}

type DeviceGroupForm struct {
	Id   uint   `json:"id"`
	Name string `json:"name" validate:"required"`
}

func (gf *DeviceGroupForm) ToDeviceGroup() *model.DeviceGroup {
	group := &model.DeviceGroup{}
	group.Id = gf.Id
	group.Name = gf.Name
	return group
}
