package datastore

import (
	"context"

	"goamin/server/model"
)

func (s *storage) CreateUserlog(c context.Context, v *model.Userlog) error {
	return s.engine.WithContext(c).Create(v).Error
}

func (s *storage) DeleteUserlog(c context.Context, v *model.Userlog) error {
	return s.engine.WithContext(c).Delete(&v).Error
}

func (s *storage) GetUserlogDetail(c context.Context, req *model.DetailOptions) (data *model.Userlog, err error) {
	tx := s.engine.WithContext(c).Model(&model.Userlog{})
	switch req.DescribeBy {
	case model.DESCRIBE_BY_ID:
		if err := tx.Where("id = ?", req.Id).First(&data).Error; err != nil {
			return nil, err
		}
	case model.DESCRIBE_BY_NAME:
		if err := tx.Where("username = ?", req.Name).First(&data).Error; err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (s *storage) GetUserlogList(c context.Context, req *model.UserlogListOptions) (data *model.UserlogSet, err error) {
	var (
		count  int64
		slices = make([]*model.Userlog, 0)
		tx     = s.engine.WithContext(c).Model(&model.Userlog{})
	)

	if req.All {
		if err := tx.Count(&count).Find(&slices).Error; err != nil {
			return nil, err
		}
		return &model.UserlogSet{Total: count, Items: slices}, nil
	}
	if req.PerPage < 1 {
		req.Page = 1
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Keyword != "" {
		filter := "%" + req.Keyword + "%"
		tx = tx.Where("username LIKE ? OR ip_addr LIKE ? OR start_at LIKE ? OR path LIKE ? OR status LIKE ?", filter, filter, filter, filter, filter)
	}

	if err := tx.Count(&count).Offset(req.PerPage * (req.Page - 1)).Limit(req.PerPage).Order("created_at DESC").Find(&slices).Error; err != nil {
		return nil, err
	}

	return &model.UserlogSet{Total: count, Items: slices}, nil
}

func (s *storage) EmptyUserlog(c context.Context) error {
	return s.engine.WithContext(c).Where("1 = 1").Delete(&model.Userlog{}).Error
}
