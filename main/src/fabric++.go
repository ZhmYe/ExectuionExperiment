package src

type FabricPP struct {
	rate  float64
	txs   []*Transaction
	cg    CG
	order []int
}

func newFabricPP(txs []*Transaction) *FabricPP {
	fabricPp := new(FabricPP)
	fabricPp.rate = 0
	fabricPp.txs = txs
	return fabricPp
}

func (f *FabricPP) TransactionSort() {
	f.cg = *newCG(f.txs)
	f.cg.getSubGraph()
	//fmt.Println(len(f.cg.subGraph))
	f.cg.getAllCycles()
	//fmt.Println(len(f.cg.cycles))
	f.cg.TransactionAbort()
	f.DAGSort()
}

func (f *FabricPP) DAGSort() {
	f.order = TopologicalOrder(f.cg.graph)
	//fmt.Println(len(f.order))
}
func (f *FabricPP) getAbortRate() float64 {
	abort := 0
	for _, tx := range f.txs {
		if tx.abort {
			abort += 1
		}
	}
	f.rate = float64(abort) / float64(len(f.txs))
	return f.rate
}
