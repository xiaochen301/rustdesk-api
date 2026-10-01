package service

import (
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

type GroupService struct {
}

// InfoById 根据用户id取用户信息
func (us *GroupService) InfoById(id uint) *model.Group {
	u := &model.Group{}
	DB.Where("id = ?", id).First(u)
	if u.Id > 0 {
		u.AdminIds = us.AdminIds(u.Id) // XC: 填充群组管理员列表
	}
	return u
}

// XC: AdminIds 群组管理员用户ID列表
func (us *GroupService) AdminIds(groupId uint) []uint {
	var ids []uint
	DB.Model(&model.GroupAdmin{}).Where("group_id = ?", groupId).Order("id asc").Pluck("user_id", &ids)
	return ids
}

// XC: SetAdmins 设置群组管理员（全量替换）
func (us *GroupService) SetAdmins(groupId uint, userIds []uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", groupId).Delete(&model.GroupAdmin{}).Error; err != nil {
			return err
		}
		for _, uid := range userIds {
			if uid == 0 {
				continue
			}
			if err := tx.Create(&model.GroupAdmin{GroupId: groupId, UserId: uid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (us *GroupService) List(page, pageSize uint, where func(tx *gorm.DB)) (res *model.GroupList) {
	res = &model.GroupList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := DB.Model(&model.Group{})
	if where != nil {
		where(tx)
	}
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Find(&res.Groups)
	// XC: 填充群组管理员列表
	for _, g := range res.Groups {
		g.AdminIds = us.AdminIds(g.Id)
	}
	return
}

// Create 创建
func (us *GroupService) Create(u *model.Group) error {
	res := DB.Create(u).Error
	return res
}
func (us *GroupService) Delete(u *model.Group) error {
	return DB.Delete(u).Error
}

// Update 更新
func (us *GroupService) Update(u *model.Group) error {
	return DB.Model(u).Updates(u).Error
}

// DeviceGroupInfoById 根据用户id取用户信息
func (us *GroupService) DeviceGroupInfoById(id uint) *model.DeviceGroup {
	u := &model.DeviceGroup{}
	DB.Where("id = ?", id).First(u)
	return u
}

func (us *GroupService) DeviceGroupList(page, pageSize uint, where func(tx *gorm.DB)) (res *model.DeviceGroupList) {
	res = &model.DeviceGroupList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := DB.Model(&model.DeviceGroup{})
	if where != nil {
		where(tx)
	}
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Find(&res.DeviceGroups)
	return
}

func (us *GroupService) DeviceGroupCreate(u *model.DeviceGroup) error {
	res := DB.Create(u).Error
	return res
}
func (us *GroupService) DeviceGroupDelete(u *model.DeviceGroup) error {
	return DB.Delete(u).Error
}

func (us *GroupService) DeviceGroupUpdate(u *model.DeviceGroup) error {
	return DB.Model(u).Updates(u).Error
}
