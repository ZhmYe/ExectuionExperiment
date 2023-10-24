package src

import (
	"fmt"
	"strconv"
	"time"
)

type algorithmType int

const (
	BasicFabric algorithmType = iota
	FabricPlusPlus
	Nezha
)

// AbortRate_RunningTime_Evaluation 测试fabric、fabric++、nezha的中止率和执行时间实验
func AbortRate_RunningTime_Evaluation() ([3][6][5]time.Duration, [3][6][5]float64) {
	smallbank := TestSmallbank(false)
	blockSize := 200
	timeResult := new([3][6][5]time.Duration)
	rateResult := new([3][6][5]float64)
	for block := 0; block < 6; block++ {
		txNumber := blockSize * (block*2 + 2)
		for hotRate := 0; hotRate < 5; hotRate++ {
			config.ZipfianConstant = float64(hotRate)*0.2 + 0.2
			if config.ZipfianConstant == 1 {
				config.ZipfianConstant = 0.999
			}
			smallbank.UpdateZipfian()
			for a := 0; a <= 0; a++ {
				startTime := time.Now()
				rate := float64(0)
				for i := 0; i < 10; i++ {
					txs := smallbank.GenTxSet(txNumber)
					switch a {
					//case 0:
					//	f := newFabric(txs)
					//	f.TransactionSort()
					//	rate += f.getAbortRate()
					case 0:
						f := newFabricPP(txs)
						f.TransactionSort()
						rate += f.getAbortRate()
						//case 1:
						//	f := newNeZha(txs)
						//	f.TransactionSort()
						//	rate += f.getAbortRate()
					}
				}
				timeResult[a][block][hotRate] = time.Since(startTime) / 10
				rateResult[a][block][hotRate] = rate / 10
			}
		}
	}
	return *timeResult, *rateResult
}

// Instance_Not_Miss_Evaluation 针对单一instance
// 不同试探性执行度（提前执行了多少个区块）、冲突率下的命中率
func Instance_Not_Miss_Evaluation() {
	//取skew=0.6,0.8,0.99，instance产生至多11个区块，第一个区块假设缓慢到达不执行，后续n - 1个区块先执行
	//命中率测试，先取出第一个区块内所有交易的读、写集，然后将后面n - 1个区块的读、写集取出，得到后者读集与前者写集的交集
	//命中率 = |交集| / |后者读集|
	fmt.Println("=========================试探性执行度命中率测试===============")
	for skew := 0.6; skew <= 1; skew += 0.2 {
		fmt.Println("	skew=" + strconv.FormatFloat(skew, 'f', 2, 64))
		if skew == 1 {
			skew = 0.99
		}
		config.ZipfianConstant = skew
		globalSmallBank.UpdateZipfian()
		instance := newInstance(time.Duration(10), 0)
		//instance.maxBlockNumber = 11 // 试探性执行度1 ~ 10
		instance.start()

		firstBlock := instance.blocks[0]
		preWriteSet := make(map[string]bool, 0)
		for _, tx := range firstBlock.txs {
			for _, op := range tx.Ops {
				if op.Type == OpWrite {
					preWriteSet[op.Key] = true
				}
			}
		}
		// 试探性执行度1 ~ 10
		for n := 1; n <= 10; n++ {
			notMiss := 0
			latterReadSet := make(map[string]bool, 0)
			// 依次遍历试探性执行的区块
			for i := 1; i <= n; i++ {
				block := instance.blocks[i]
				for _, tx := range block.txs {
					for _, op := range tx.Ops {
						if op.Type == OpRead {
							latterReadSet[op.Key] = true
							_, exist := preWriteSet[op.Key]
							if exist {
								notMiss++
							}
						}
					}
				}
			}
			fmt.Println("		试探性执行度: " + strconv.Itoa(n) + " , 命中率: " + strconv.FormatFloat(float64(notMiss)/float64(len(latterReadSet))*100, 'f', 2, 64) + "%")
		}
	}
}

// Instance_Abort_Evaluation 针对单一instance
// 不同试探性执行度、冲突率下的中止率
// 为了方便写，我们生成n个区块，前n-1个区块先执行，最后一个区块当做实验的第一个区块
func Instance_Abort_Evaluation() {
	fmt.Println("=========================试探性执行度中止率测试===============")
	for skew := 0.6; skew <= 1; skew += 0.2 {
		fmt.Println("	skew=" + strconv.FormatFloat(skew, 'f', 2, 64))
		if skew == 1 {
			skew = 0.99
		}
		config.ZipfianConstant = skew
		globalSmallBank.UpdateZipfian()
		// 试探性执行度1 ~ 10
		for n := 1; n <= 10; n++ {
			instance := newInstance(time.Duration(10), 0)
			//instance.maxBlockNumber = 11 // 试探性执行度1 ~ 10
			instance.start()
			instance.simulateExecution(n)
			abort := 0
			firstBlock := instance.blocks[n]
			writeAddress := make(map[string]bool, 0)
			for _, tx := range firstBlock.txs {
				for _, op := range tx.Ops {
					if op.Type == OpWrite {
						writeAddress[op.Key] = true
					}
				}
			}
			instance.CascadeAbort(&writeAddress)
			for _, block := range instance.blocks[:n] {
				for _, tx := range block.txs {
					if tx.abort {
						abort++
					}
				}
			}
			//nezha := newNeZha(instance.blocks[n].txs)
			//nezha.TransactionSort()
			//for _, tx := range nezha.txs {
			//	if tx.abort {
			//		abort++
			//	}
			//}
			fmt.Print(abort)
			fmt.Print(" ")
			//fmt.Println()
			fmt.Println(float64(abort) / float64((n)*config.BlockSize))
		}
	}
}
func Instance_ReExectution_Time_Evaluation() {
	fmt.Println("=========================试探性执行度运行时间测试===============")
	// 先测试下第一个块并发预执行+abort + 级联abort + 重执行时间
	// 然后测试所有块串行并发预执行abort + 重执行时间
	// 其实重执行时间应该差不多，主要就是看级联abort时间是否比剩下块串行并发预执行abort的时间长
	for skew := 0.6; skew <= 1; skew += 0.2 {
		fmt.Println("	skew=" + strconv.FormatFloat(skew, 'f', 2, 64))
		if skew == 1 {
			skew = 0.99
		}
		config.ZipfianConstant = skew
		globalSmallBank.UpdateZipfian()
		// 试探性执行度1 ~ 10
		for n := 1; n <= 10; n++ {
			instance := newInstance(time.Duration(10), 0)
			//instance.maxBlockNumber = 11 // 试探性执行度1 ~ 10
			instance.start()
			instance.simulateExecution(n)
			//  到这里前n个块模拟执行完成
			// 模拟执行第1个块
			s := newSimulateEngine(instance.blocks[n : n+1])
			s.SimulateExecution()
			instance.acgs = append(instance.acgs, s.acgs[0])
			firstBlock := instance.blocks[n]
			writeAddress := make(map[string]bool, 0)
			for _, tx := range firstBlock.txs {
				if tx.abort {
					continue
				}
				for _, op := range tx.Ops {
					if op.Type == OpWrite {
						writeAddress[op.Key] = true
					}
				}
			}
			// 级联abort
			instance.CascadeAbort(&writeAddress)
			abortTxs := instance.getAbortTxs(n + 1)
			// 执行未abort的交易,只需对每个acg的每个address的最后一个写进行操作即可
			instance.execute(n + 1)
			startTime := time.Now()
			for i := 0; i < 100; i++ {
				instance.reExecute(abortTxs)
			}
			fmt.Print("重执行时间")
			fmt.Println(time.Since(startTime) / 100)
			fmt.Println()
			// 重执行
		}
	}
}
func Instance_Execution_Time_Evaluation() {
	// 跑的时候把maxblocknumber改为2000
	for skew := 0.6; skew <= 1; skew += 0.2 {
		fmt.Println("	skew=" + strconv.FormatFloat(skew, 'f', 2, 64))
		if skew == 1 {
			skew = 0.99
		}
		config.ZipfianConstant = skew
		globalSmallBank.UpdateZipfian()
		// 试探性执行度1 ~ 10
		for n := 2; n <= 11; n++ {
			instance := newInstance(time.Duration(10), 0)
			//instance.maxBlockNumber = 11 // 试探性执行度1 ~ 10
			instance.start()
			startTime := time.Now()
			for i := 0; i < 100; i++ {
				instance.simulateExecution(n)
				abortTxs := instance.execute(n)
				instance.reExecute(abortTxs)
			}
			fmt.Println(time.Since(startTime) / 100)
		}
	}
}

// Instance_Waiting_Time_Evalutation 测试新方案不同速度的instance不同区块高度的区块的等待时间
func Instance_Waiting_Time_Evalutation() {
	peer := newPeer(4)
	peer.run()
	for _, instance := range peer.instances {
		for i, block := range instance.blocks {
			if i >= 20 {
				continue
			}
			fmt.Print("区块高度: ")
			fmt.Print(i)
			fmt.Print(" 等待时间: ")
			fmt.Println(block.finishTime)
		}
	}
}

// Baseline_Waiting_Time_Evalutation 测试baseline不同速度的instance不同区块高度的区块的等待时间
func Baseline_Waiting_Time_Evalutation() {
	peer := newPeer(4)
	peer.runInBaseline()
	for _, instance := range peer.instances {
		for i, block := range instance.blocks {
			if i >= 24 {
				continue
			}
			fmt.Print("区块高度: ")
			fmt.Print(i)
			fmt.Print(" 等待时间: ")
			fmt.Println(block.finishTime)
		}
	}
}

// Instance_Not_Execute_Block_Number_Evaluation 测试新方案不同速度的instance随着时间未执行的区块个数
func Instance_Not_Execute_Block_Number_Evaluation() {
	peer := newPeer(4)
	peer.run()
}

// Baseline_Not_Execute_Block_Number_Evaluation 测试baseline不同速度的instance随着时间未执行的区块个数
func Baseline_Not_Execute_Block_Number_Evaluation() {
	peer := newPeer(4)
	peer.runInBaseline()
}
