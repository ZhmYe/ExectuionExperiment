package src

type FabricPP struct {
	rate float64
	txs  []*Transaction
}

func newFabricPP(txs []*Transaction) *FabricPP {
	fabricPp := new(FabricPP)
	fabricPp.rate = 0
	fabricPp.txs = txs
	return fabricPp
}

func (fabricPP *FabricPP) TransactionSort() {
	return
}

func (fabricPP *FabricPP) getAbortRate() float64 {
	abort := 0
	for _, tx := range fabricPP.txs {
		if tx.abort {
			abort += 1
		}
	}
	fabricPP.rate = float64(abort) / float64(len(fabricPP.txs))
	return fabricPP.rate
}
