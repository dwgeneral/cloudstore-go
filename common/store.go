package common

// StoreType 存储类型(表示文件存到哪里)
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : 节点本地
	StoreLocal
	// StoreOSS : 阿里OSS
	StoreOSS
	// StoreAll : 所有类型的存储都存一份数据
	StoreAll
)