package src

// Config 全局配置定义
type Config struct {
	OriginKeys      int // 初始Key的数量
	HotKey          float64
	HotKeyRate      float64 // 有HotKeyRate的交易访问HotKey的状态
	path            string
	ZipfianConstant float64
	BlockSize       int // 区块大小
	instanceNumber  int // instance个数
}

var config = Config{OriginKeys: 10000, HotKey: 0.2, HotKeyRate: 1, path: "leveldb", ZipfianConstant: 0.8, BlockSize: 200, instanceNumber: 4}
