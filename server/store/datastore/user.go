package datastore

import (
	"context"

	"gorm.io/gorm"

	"goamin/server/model"
)

func (s *storage) CreateUser(c context.Context, v *model.User) error {
	return s.engine.WithContext(c).Create(v).Error
}

func (s *storage) DeleteUser(c context.Context, v *model.User) error {
	return s.engine.WithContext(c).Select("Roles").Delete(&v).Error
}

func (s *storage) UpdateUser(c context.Context, v *model.User) (err error) {
	return s.engine.WithContext(c).Model(&v).Transaction(func(tx *gorm.DB) error {
		err = tx.Association("Roles").Replace(v.Roles)
		if err != nil {
			return err
		}
		err = tx.Where("id = ?", v.ID).Select("password", "nickname", "email",
			"phone", "avatar", "status", "description", "access_token", "refresh_token", "roles").Updates(v).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *storage) GetUserDetail(c context.Context, req *model.DetailUserRequest) (data *model.User, err error) {
	tx := s.engine.WithContext(c).Model(&model.User{})
	switch req.DescribeBy {
	case model.DESCRIBEUSER_BY_ID:
		if err := tx.Preload("Roles").Where("id = ?", req.Id).First(&data).Error; err != nil {
			return nil, err
		}
	case model.DESCRIBEUSER_BY_NAME:
		if err := tx.Preload("Roles").Where("username = ?", req.Name).First(&data).Error; err != nil {
			return nil, err
		}
	case model.DESCRIBEUSER_BY_AK:
		if err := tx.Preload("Roles").Where("access_token = ?", req.AK).First(&data).Error; err != nil {
			return nil, err
		}
	case model.DESCRIBEUSER_BY_EMAIL:
		if err := tx.Preload("Roles").Where("email = ?", req.Email).First(&data).Error; err != nil {
			return nil, err
		}
	case model.DESCRIBEUSER_BY_PHONE:
		if err := tx.Preload("Roles").Where("phone = ?", req.Phone).First(&data).Error; err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (s *storage) GetUserList(c context.Context, req *model.UserListOptions) (data *model.UserSet, err error) {
	var (
		count  int64
		slices = make([]*model.User, 0)
		tx     = s.engine.WithContext(c).Model(&model.User{})
	)

	if req.All {
		if err := tx.Preload("Roles").Count(&count).Find(&slices).Error; err != nil {
			return nil, err
		}
		return &model.UserSet{Total: count, Items: slices}, nil
	}
	if req.PerPage < 1 {
		req.Page = 1
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Keyword != "" {
		filter := "%" + req.Keyword + "%"
		tx = tx.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ?", filter, filter, filter, filter)
	}
	if req.Username != "" {
		tx = tx.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Nickname != "" {
		tx = tx.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Email != "" {
		tx = tx.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Phone != "" {
		tx = tx.Where("phone LIKE ?", "%"+req.Phone+"%")
	}

	if err := tx.Preload("Roles").Count(&count).Offset(req.PerPage * (req.Page - 1)).Limit(req.PerPage).Find(&slices).Error; err != nil {
		return nil, err
	}

	return &model.UserSet{Total: count, Items: slices}, nil
}

func (s *storage) GetUserCount(c context.Context) (count int64, err error) {
	if err := s.engine.WithContext(c).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
