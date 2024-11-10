package model

type ServerConfig struct {
	SKey  string `json:"skey"   gorm:"column:skey;primarykey"`
	Value string `json:"value" gorm:"type:longtext"`
}

func (ServerConfig) TableName() string {
	return "server_configs"
}
