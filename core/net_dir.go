package core

type NetDirTrace struct {
	Name string `json:"name"`
}

// 115的网络目录
type NetDir struct {
	Count int            `json:"count"`
	Path  []*NetDirTrace `json:"path"`
	Data  []*NetFile     `json:"data"`
}
