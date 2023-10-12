package src

import (
	"fmt"
	"strconv"
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
	txNumber := 200 * 1
	smallbank := TestSmallbank(false)
	txs := smallbank.GenTxSet(txNumber)
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
