package core

type NetDirTrace struct {
	Name string `json:"name"`
}

// NetDir 115的网络目录
type NetDir struct {
	Count int            `json:"count"`
	Path  []*NetDirTrace `json:"path"`
	Data  []*NetFile     `json:"data"`
}
