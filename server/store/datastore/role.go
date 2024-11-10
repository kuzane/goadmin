package datastore

import (
    "context"

    "github.com/mikespook/gorbac"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "goamin/server"
    "goamin/server/model"
)

func (s *storage) CreateRole(c context.Context, v *model.Role) (err error) {
    if err = s.engine.WithContext(c).Create(v).Error; err != nil {
        return
    }
    if server.Config.Server.RBAC, err = s.SetRBAC(c); err != nil {
        return
    }
    return nil
}
func (s *storage) DeleteRole(c context.Context, v *model.Role) (err error) {
    err = s.engine.WithContext(c).Model(&v).Transaction(func(tx *gorm.DB) error {
        err = tx.Select(clause.Associations).Delete(&v).Error
        if err != nil {
            return err
        }

        return nil
    })
    if err != nil {
        return err
    }

    //  需要先实现role变更之后再变更权限
    if server.Config.Server.RBAC, err = s.SetRBAC(c); err != nil {
        return err
    }
    return nil
}
func (s *storage) UpdateRole(c context.Context, v *model.Role) (err error) {
    err = s.engine.WithContext(c).Model(&v).Transaction(func(tx *gorm.DB) error {
        err = tx.Association("Users").Replace(v.Users)
        if err != nil {
            return err
        }
        err = tx.Association("Endpoints").Replace(v.Endpoints)
        if err != nil {
            return err
        }
        // updates  默认更新非零值，使用select选择需要更新的字段，零值也会更新
        err = tx.Where("id = ?", v.ID).Select("rolename", "nickname", "description", "parents", "users", "endpoints", "status").Updates(v).Error
        if err != nil {
            return err
        }

        return nil
    })
    if err != nil {
        return err
    }

    //  需要先实现role变更之后再变更权限
    if server.Config.Server.RBAC, err = s.SetRBAC(c); err != nil {
        return err
    }

    return nil
}
func (s *storage) GetRoleDetail(c context.Context, req *model.DetailOptions) (data *model.Role, err error) {
    tx := s.engine.WithContext(c).Model(&model.Role{})
    switch req.DescribeBy {
    case model.DESCRIBE_BY_ID:
        if err := tx.Preload("Endpoints").Preload("Users").Where("id = ?", req.Id).First(&data).Error; err != nil {
            return nil, err
        }
    case model.DESCRIBE_BY_NAME:
        if err := tx.Preload("Endpoints").Preload("Users").Where("rolename = ?", req.Name).First(&data).Error; err != nil {
            return nil, err
        }
    }
    return data, nil
}
func (s *storage) GetRoleList(c context.Context, req *model.RoleListOptions) (data *model.RoleSet, err error) {
    var (
        count  int64
        slices = make([]*model.Role, 0)
        tx     = s.engine.WithContext(c).Model(&model.Role{})
    )

    if req.All {
        // 因为RBAC设置时需要以来Endpoints,而且在RBAC设置时对role进行全表查询的,所以此处需要联表查询,这样可以保证权限的变更
        if err := tx.Count(&count).Preload("Endpoints").Find(&slices).Error; err != nil {
            return nil, err
        }
        return &model.RoleSet{Total: count, Items: slices}, nil
    }

    if req.PerPage < 1 {
        req.Page = 1
    }
    if req.Page < 1 {
        req.Page = 1
    }
    if req.Keyword != "" {
        filter := "%" + req.Keyword + "%"
        tx = tx.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ?", filter, filter, filter, filter, filter)
    }
    if req.Rolename != "" {
        tx = tx.Where("username LIKE ?", "%"+req.Rolename+"%")
    }
    if req.Nickname != "" {
        tx = tx.Where("nickname LIKE ?", "%"+req.Nickname+"%")
    }

    if err := tx.Count(&count).Preload("Users").Offset(req.PerPage * (req.Page - 1)).Limit(req.PerPage).Find(&slices).Error; err != nil {
        return nil, err
    }

    return &model.RoleSet{Total: count, Items: slices}, nil
}
func (s *storage) GetRoleCount(c context.Context) (count int64, err error) {
    if err := s.engine.WithContext(c).Model(&model.Role{}).Count(&count).Error; err != nil {
        return 0, err
    }

    return count, nil
}

func (s *storage) SetRBAC(c context.Context) (*gorbac.RBAC, error) {
    rbac := gorbac.New()
    // 如果角色之间有继承关系参考: https://mikespook.com/2017/04/how-to-persist-gorbac-instance/
    data, err := s.GetRoleList(c, &model.RoleListOptions{ListOptions: &model.ListOptions{All: true}})

    if err != nil {
        return nil, err
    }
    for _, role := range data.Items {
        if role.Status {
            r := gorbac.NewStdRole(role.Rolename)
            // 角色授权接口权限
            for _, endpoint := range role.Endpoints {
                p := endpoint.Method + ":" + endpoint.Path
                r.Assign(gorbac.NewStdPermission(p))
            }
            rbac.Add(r)
        }
    }

    // 分配角色之间的继承关系
    for _, role := range data.Items {
        if role.Status {
            rbac.SetParents(role.Rolename, role.Parents)
        }
    }

    //get, parents, err := rbac.Get("test")
    //log.Info().Msgf("-----------setupDatabaseRBAC---------%v,%v,%v", get, parents, err)
    return rbac, nil
}
