package datastore

import (
    "context"

    "goamin/server/model"
)

func (s storage) ServerConfigDelete(c context.Context, k string) error {
    config := &model.ServerConfig{
        SKey: k,
    }

    return s.engine.WithContext(c).Create(config).Error
}

func (s storage) ServerConfigSet(c context.Context, k, v string) error {
    return s.engine.WithContext(c).Save(&model.ServerConfig{SKey: k, Value: v}).Error
}

func (s storage) ServerConfigGet(c context.Context, k string) (string, error) {
    config := new(model.ServerConfig)
    if err := s.engine.WithContext(c).Where("skey = ?", k).First(&config).Error; err != nil {
        return "", err
    }
    return config.Value, nil
}
