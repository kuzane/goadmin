package store

type Opts struct {
	Driver      string
	Config      string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifeTime int
	MaxIdleTime int
	ShowSQL     bool
}
