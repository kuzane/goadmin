package datastore

import (
    "context"
    "errors"

    "gorm.io/gorm"

    "goamin/server/model"
)

func (s *storage) CreateEndpoint(c context.Context, v *model.Endpoint) (err error) {
    if err = s.engine.WithContext(c).Create(v).Error; err != nil {
        return
    }
    return
}

func (s *storage) DeleteEndpoint(c context.Context, v *model.Endpoint) (err error) {
    if err = s.engine.WithContext(c).Select("Roles").Delete(&v).Error; err != nil {
        return
    }
    return
}

func (s *storage) UpdateEndpoint(c context.Context, v *model.Endpoint) (err error) {
    if err = s.engine.WithContext(c).Where("id = ?", v.ID).Updates(v).Error; err != nil {
        return
    }
    return
}

func (s *storage) GetEndpointDetail(c context.Context, req *model.DetailOptions) (data *model.Endpoint, err error) {
    tx := s.engine.WithContext(c).Model(&model.Endpoint{})
    switch req.DescribeBy {
    case model.DESCRIBE_BY_ID:
        if err := tx.Preload("Roles").Where("id = ?", req.Id).First(&data).Error; err != nil {
            return nil, err
        }
    case model.DESCRIBE_BY_NAME:
        if err := tx.Preload("Roles").Where("identity = ?", req.Name).First(&data).Error; err != nil {
            return nil, err
        }
    }
    return data, nil
}
func (s *storage) GetEndpointList(c context.Context, req *model.EndpointListOptions) (data *model.EndpointSet, err error) {
    var (
        count  int64
        slices = make([]*model.Endpoint, 0)
        tx     = s.engine.WithContext(c).Model(&model.Endpoint{})
    )

    if req.All {
        if err := tx.Count(&count).Find(&slices).Error; err != nil {
            return nil, err
        }
        return &model.EndpointSet{Total: count, Items: slices}, nil
    }
    if req.PerPage < 1 {
        req.Page = 1
    }
    if req.Page < 1 {
        req.Page = 1
    }
    if req.Keyword != "" {
        filter := "%" + req.Keyword + "%"
        tx = tx.Where("path LIKE ? OR method LIKE ? OR kind LIKE ? OR remark LIKE ?", filter, filter, filter, filter, filter)
    }
    if req.Path != "" {
        tx = tx.Where("path LIKE ?", "%"+req.Path+"%")
    }
    if req.Method != "" {
        tx = tx.Where("method LIKE ?", "%"+req.Method+"%")
    }
    if req.Kind != "" {
        tx = tx.Where("kind LIKE ?", "%"+req.Kind+"%")
    }
    if req.Remark != "" {
        tx = tx.Where("remark LIKE ?", "%"+req.Remark+"%")
    }

    if err := tx.Count(&count).Offset(req.PerPage * (req.Page - 1)).Limit(req.PerPage).Find(&slices).Error; err != nil {
        return nil, err
    }

    return &model.EndpointSet{Total: count, Items: slices}, nil
}
func (s *storage) GetEndpointCount(c context.Context) (count int64, err error) {
    if err := s.engine.WithContext(c).Model(&model.Endpoint{}).Count(&count).Error; err != nil {
        return 0, err
    }

    return count, nil
}

func (s *storage) SetEndpoint(c context.Context, v *model.Endpoint) error {
    data, err := s.GetEndpointDetail(c, model.NewDescribeRequestByName(v.Identity))
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return s.CreateEndpoint(c, v)
    }
    v.ID = data.ID
    return s.UpdateEndpoint(c, v)
}
