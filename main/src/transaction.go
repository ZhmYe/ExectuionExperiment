package src

import "time"

// OpType 操作类型
type OpType int

const (
	OpRead  OpType = iota // 读操作
	OpWrite               // 写操作
)

// Op 操作
type Op struct {
	Type        OpType // 操作类型 读/写
	Key         string // 操作的key
	Val         string // 操作的value
	ReadResult  string // 最终读到的结果
	WriteResult string // 最终要写的结果
}

// TxType 交易类型, smallbank
type TxType int

const (
	transactSavings TxType = iota
	depositChecking
	sendPayment
	writeCheck
	amalgamate
	query
)

// Transaction 交易
type Transaction struct {
	txType   TxType // 交易类型
	Ops      []*Op  // 交易中包含的操作
	abort    bool   // 是否abort
	sequence int    // sorting时的序列号
	txHash   string // 交易哈希
}

// Block 区块
type Block struct {
	txs        []*Transaction // 交易
	createTime time.Time      // 被创建的时间，用于衡量等待时间
	finish     bool
	finishTime time.Duration
}

func newBlock(txs []*Transaction) *Block {
	block := new(Block)
	block.txs = txs
	block.createTime = time.Now()
	block.finish = false
	return block
}
