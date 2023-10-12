package src

// Unit 操作单元，每一笔交易的某一个read/write操作
type Unit struct {
	op Op           // 实际执行的操作
	tx *Transaction // 交易标识
}

func newUnit(op Op, tx *Transaction) *Unit {
	unit := new(Unit)
	unit.op = op
	unit.tx = tx
	return unit
}

// StateSet

// StateSet ACG中的每一行
type StateSet struct {
	ReadSet  []Unit // 读集
	WriteSet []Unit // 写集
}

func newStateSet() *StateSet {
	set := new(StateSet)
	set.ReadSet = make([]Unit, 0)
	set.WriteSet = make([]Unit, 0)
	return set
}
func (stateSet *StateSet) appendToReadSet(unit Unit) {
	stateSet.ReadSet = append(stateSet.ReadSet, unit)
}
func (stateSet *StateSet) appendToWriteSet(unit Unit) {
	stateSet.WriteSet = append(stateSet.WriteSet, unit)
}

// ACG address->StateSet
type ACG = map[string]StateSet

// 构建并发交易所对应的ACG
func getACG(txs []*Transaction) ACG {
	acg := make(ACG)
	for _, tx := range txs {
		for _, op := range tx.Ops {
			_, exist := acg[op.Key]

			// 如果在acg中不存在address,新建一个StateSet
			if !exist {
				acg[op.Key] = *newStateSet()
			}

			unit := newUnit(op, tx) // 新建操作单元
			stateSet := acg[op.Key]

			// 根据读/写操作加入到StateSet的两部分中
			switch unit.op.Type {
			case OpRead:
				stateSet.appendToReadSet(*unit)
			case OpWrite:
				stateSet.appendToWriteSet(*unit)
			}
			acg[op.Key] = stateSet
		}
	}
	return acg
}
