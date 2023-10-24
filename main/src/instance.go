package src

import (
	"strconv"
	"time"
)

// Instance 共识实例
type Instance struct {
	blocks                []*Block                     // 共识出的块
	timeout               time.Duration                // 出块时间, ms
	hasExecutedIndex      int                          // 最新被执行的区块下标, 默认为0
	lastBlockOutTimeStamp time.Time                    // 上次出块时间
	id                    int                          // instance id
	maxBlockNumber        int                          // 最大区块数，用于自动停下实验
	simulate              *SimulateEngine              // 模拟执行实例
	acgs                  []ACG                        //模拟执行的所有子块的acg
	record                map[string]map[string][]Unit //统计每笔交易在每个地址后面直接相连的读操作个数
	cascade               map[string]int               // 在每个地址上的级联度
	finish                bool
}

func newInstance(timeout time.Duration, id int) *Instance {
	instance := new(Instance)
	instance.blocks = make([]*Block, 0)
	instance.lastBlockOutTimeStamp = time.Now()
	instance.timeout = timeout * time.Millisecond
	instance.hasExecutedIndex = 0
	instance.id = id
	instance.maxBlockNumber = 30
	instance.finish = false
	return instance
}

// checkTimeout 判断是否应该出块
func (instance *Instance) checkTimeout() bool {
	if time.Since(instance.lastBlockOutTimeStamp) >= instance.timeout {
		return true
	}
	return false
}

// updateLastBlockOutTimeStamp 更新出块的时间戳
func (instance *Instance) updateLastBlockOutTimeStamp() {
	instance.lastBlockOutTimeStamp = time.Now()
}

// blockOut 出块
func (instance *Instance) blockOut() {
	if len(instance.blocks) >= instance.maxBlockNumber {
		return
	}
	txs := globalSmallBank.GenTxSet(config.BlockSize)
	block := newBlock(txs)
	instance.blocks = append(instance.blocks, block)
}
func (instance *Instance) start() {
	// 暂时先设置成，如果生成的区块达到一定数量就停下
	go func() {
		instance.updateLastBlockOutTimeStamp()
		for {
			if len(instance.blocks) >= instance.maxBlockNumber {
				instance.finish = true
				//fmt.Println("Instance " + strconv.Itoa(instance.id) + " finished...")
				break
			}
			if instance.checkTimeout() {
				instance.blockOut()
				instance.updateLastBlockOutTimeStamp()
			}
		}
	}()
}
func (instance *Instance) simulateExecution(number int) int {
	lastIndex := instance.hasExecutedIndex + number
	if instance.hasExecutedIndex == len(instance.blocks) {
		return 0
	}
	if instance.hasExecutedIndex+number > len(instance.blocks) {
		lastIndex = len(instance.blocks)
	}
	instance.simulate = newSimulateEngine(instance.blocks[instance.hasExecutedIndex:lastIndex])
	instance.acgs = instance.simulate.SimulateExecution()
	instance.record, instance.cascade = computeCascade(instance.acgs)
	return lastIndex - instance.hasExecutedIndex
}
func (instance *Instance) abortReadSet(readSet []Unit) {
	repeatCheck := make(map[string]bool)
	if len(readSet) == 0 {
		return
	}
	for _, unit := range readSet {
		_, repeat := repeatCheck[unit.tx.txHash]
		if repeat || unit.tx.abort {
			continue
		}
		repeatCheck[unit.tx.txHash] = true
		unit.tx.abort = true
		CascadeInAddress, haveCascade := instance.record[unit.tx.txHash]
		if haveCascade {
			for _, eachReadSet := range CascadeInAddress {
				instance.abortReadSet(eachReadSet)
			}
		}

	}
}
func (instance *Instance) CascadeAbort(writeAddress *map[string]bool) {
	hasAbort := make(map[string]bool, 0)
	localWriteAddress := make([]string, 0) // 当前acgs所涉及的写集，用于更新writeAddress
	for _, acg := range instance.acgs {
		for address, stateSet := range acg {
			if len(stateSet.WriteSet) != 0 {
				localWriteAddress = append(localWriteAddress, address)
			}
			_, exist := (*writeAddress)[address]
			// 如果读了排序在前的instance写的address，并且是第一个读的acg
			if exist {
				if len(stateSet.ReadSet) != 0 {
					_, has := hasAbort[address]
					// 第一个读的acg
					if !has {
						hasAbort[address] = true
						instance.abortReadSet(stateSet.ReadSet)
					}
				}
			}
		}
	}
	for _, address := range localWriteAddress {
		(*writeAddress)[address] = true
	}
}
func (instance *Instance) getAbortTxs(n int) []*Transaction {
	abortTxs := make([]*Transaction, 0)
	for _, block := range instance.blocks[instance.hasExecutedIndex : instance.hasExecutedIndex+n] {
		for _, tx := range block.txs {
			if tx.abort {
				abortTxs = append(abortTxs, tx)
			}
		}
	}
	return abortTxs
}
func (instance *Instance) execute(n int) []*Transaction {
	//for _, acg := range instance.acgs {
	//	for address, stateSet := range acg {
	//		if len(stateSet.WriteSet) != 0 {
	//			globalSmallBank.Write(address, stateSet.WriteSet[len(stateSet.WriteSet)-1].op.WriteResult)
	//		}
	//	}
	//}
	abortTxs := make([]*Transaction, 0)
	lastIndex := instance.hasExecutedIndex + n
	if instance.hasExecutedIndex == len(instance.blocks) {
		return abortTxs
	}
	if instance.hasExecutedIndex+n > len(instance.blocks) {
		lastIndex = len(instance.blocks)
	}
	for _, block := range instance.blocks[instance.hasExecutedIndex : instance.hasExecutedIndex+n] {
		block.finishTime = time.Since(block.createTime)
		for _, tx := range block.txs {
			if tx.abort {
				abortTxs = append(abortTxs, tx)
				continue
			}
			for _, op := range tx.Ops {
				if op.Type == OpWrite {
					globalSmallBank.Write(op.Key, op.WriteResult)
				}
			}
		}
		block.finish = true
	}
	instance.hasExecutedIndex = lastIndex
	return abortTxs
}
func (instance *Instance) reExecute(txs []*Transaction) {
	for _, tx := range txs {
		for _, op := range tx.Ops {
			if op.Type == OpRead {
				op.ReadResult = globalSmallBank.Read(op.Key)
			} else {
				readResult, _ := strconv.Atoi(globalSmallBank.Read(op.Key))
				amount, _ := strconv.Atoi(op.Val)
				op.WriteResult = strconv.Itoa(readResult + amount)
				globalSmallBank.Write(op.Key, op.WriteResult)
			}
		}
	}
}
