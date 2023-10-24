package src

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func TestSmallbank(output bool) *Smallbank {
	smallbank := NewSmallbank(config.path)
	fmt.Println("init smallbank success")
	if !output {
		return smallbank
	}
	return smallbank
}
func TestGenerateTransaction(output bool) []*Transaction {
	txNumber := 400
	txs := globalSmallBank.GenTxSet(txNumber)
	if !output {
		return txs
	}
	for _, tx := range txs {
		//if index%10 == 0 {
		fmt.Print(tx.txHash)
		fmt.Print(" ")
		fmt.Print(tx.txType)
		fmt.Print(" ")
		for _, op := range tx.Ops {
			fmt.Print(op.Key)
			fmt.Print(" ")
		}
		fmt.Println()
		//}
	}
	return txs
}
func TestGetACG(output bool) ([]*Transaction, ACG) {
	txs := TestGenerateTransaction(false)
	fmt.Println("generate tx success")
	fmt.Println("tx number:" + strconv.Itoa(len(txs)))
	acg := getACG(txs)
	if !output {
		return txs, acg
	}
	for address, stateSet := range acg {
		fmt.Print(address)
		fmt.Print(" ")
		fmt.Print(len(stateSet.ReadSet))
		fmt.Print(" ")
		fmt.Print(len(stateSet.WriteSet))
		fmt.Println()
	}
	return txs, acg
}
func TestBuildConflictGraph(output bool) {
	txs := TestGenerateTransaction(false)
	fmt.Println("generate tx success")
	fmt.Println("tx number:" + strconv.Itoa(len(txs)))
	graph := buildConflictGraph(txs)
	for i, _ := range graph {
		for j, _ := range graph[i] {
			fmt.Print(graph[i][j])
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
func TestTarjan(output bool) {
	txs := TestGenerateTransaction(false)
	fmt.Println("generate tx success")
	fmt.Println("tx number:" + strconv.Itoa(len(txs)))
	cg := newCG(txs)
	cg.getSubGraph()
	fmt.Println(len(cg.subGraph))
}
func TestFindCycles(output bool) {
	txs := TestGenerateTransaction(false)
	fmt.Println("generate tx success")
	fmt.Println("tx number:" + strconv.Itoa(len(txs)))
	cg := newCG(txs)
	cg.getSubGraph()
	cg.getAllCycles()
	fmt.Println(len(cg.cycles))

}
func TestZipfian(output bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	z := NewZipfianWithItems(int64(config.OriginKeys), config.ZipfianConstant)
	counter := make(map[int64]int, 0)
	for i := 0; i < 10000; i++ {
		result := z.Next(r)
		_, exist := counter[result]
		if !exist {
			counter[result] = 1
		} else {
			counter[result] += 1
		}
	}
	fmt.Println(counter)
}
func TestFabricPP(output bool) {
	txs := TestGenerateTransaction(false)
	fmt.Println("generate tx success")
	fmt.Println("tx number:" + strconv.Itoa(len(txs)))
	startTime := time.Now()
	f := newFabricPP(txs)
	f.TransactionSort()
	fmt.Println(time.Since(startTime))
	rate := f.getAbortRate()
	fmt.Println(strconv.FormatFloat(rate*100, 'f', 2, 64) + "%")
	//fmt.Println(f.order)
}
func TestNezha(output bool) {
	txs := TestGenerateTransaction(false)
	fmt.Println("generate tx success")
	fmt.Println("tx number:" + strconv.Itoa(len(txs)))
	//TransactionSort(acg)
	nezha := newNeZha(txs)
	nezha.TransactionSort()
	rate := nezha.getAbortRate()
	fmt.Println(strconv.FormatFloat(rate*100, 'f', 2, 64) + "%")
}
func TestFabric(output bool) {
	txs := TestGenerateTransaction(false)
	fmt.Println("generate tx success")
	fmt.Println("tx number:" + strconv.Itoa(len(txs)))
	//TransactionSort(acg)
	fabric := newFabric(txs)
	fabric.TransactionSort()
	rate := fabric.getAbortRate()
	fmt.Println(strconv.FormatFloat(rate*100, 'f', 2, 64) + "%")
}
func TestInstance(output bool) {
	instance := newInstance(time.Duration(50), 0)
	//startTime := time.Now()
	instance.start()
	fmt.Println(len(instance.blocks))
	//fmt.Println(time.Since(startTime))
}
func TestSimulateEngine(output bool) {
	instance := newInstance(time.Duration(50), 0)
	//startTime := time.Now()
	instance.maxBlockNumber = 2
	instance.start()
	simulate := newSimulateEngine(instance.blocks)
	simulate.SimulateExecution()
	abort := 0
	for _, block := range simulate.blocks {
		for _, tx := range block.txs {
			if tx.abort {
				abort++
			}
		}
	}
	fmt.Println(float64(abort) / float64(len(simulate.blocks)*config.BlockSize))
}
func TestInstanceCascadeAbort(output bool) {
	instance := newInstance(time.Duration(50), 0)
	//startTime := time.Now()
	instance.maxBlockNumber = 2
	instance.start()
	instance.simulateExecution(1)
	//simulate := newSimulateEngine(instance.blocks[1:])
	//simulate.SimulateExecution()
	abort := 0
	for _, block := range instance.blocks[:1] {
		for _, tx := range block.txs {
			if tx.abort {
				abort++
			}
		}
	}
	fmt.Println(float64(abort) / float64(len(instance.blocks[:1])*config.BlockSize))
	abort = 0
	firstBlock := instance.blocks[1]
	writeAddress := make(map[string]bool, 0)
	for _, tx := range firstBlock.txs {
		for _, op := range tx.Ops {
			if op.Type == OpWrite {
				writeAddress[op.Key] = true
			}
		}
	}
	instance.CascadeAbort(&writeAddress)
	for _, block := range instance.blocks[:1] {
		for _, tx := range block.txs {
			if tx.abort {
				abort++
			}
		}
	}
	fmt.Println(float64(abort) / float64(len(instance.blocks[:1])*config.BlockSize))
	nezha := newNeZha(instance.blocks[1].txs)
	nezha.TransactionSort()
	for _, tx := range nezha.txs {
		if tx.abort {
			abort++
		}
	}
	fmt.Println(float64(abort) / float64(len(instance.blocks[:2])*config.BlockSize))

}
func TestPeer(output bool) {
	peer := newPeer(4)
	peer.run()

}
func TestBaseline(output bool) {
	peer := newPeer(4)
	peer.runInBaseline()
}
