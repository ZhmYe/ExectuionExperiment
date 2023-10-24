package main

import (
	"main/src"
	"math/rand"
	"time"
)

func test() {
	//src.TestSmallbank(true)
	//src.TestGenerateTransaction(true)
	//src.TestGetACG(true)
	//src.TestBuildConflictGraph(true)
	//src.TestTarjan(true)
	//src.TestFindCycles(true)
	//src.TestNezha(true)
	//src.TestFabric(true)
	//src.TestFabricPP(true)
	//src.TestZipfian(true)
	//src.TestInstance(true)
	//src.TestSimulateEngine(true)
	//src.TestInstanceCascadeAbort(true)
	//src.TestPeer(true)
	src.TestBaseline(true)
}
func eval() {
	//src.AbortRate_RunningTime_Evaluation()
	//src.Instance_Not_Miss_Evaluation()
	//src.Instance_Abort_Evaluation()
	//src.Instance_ReExectution_Time_Evaluation()
	//src.Instance_Execution_Time_Evaluation()
	//src.Instance_Waiting_Time_Evalutation()
	//src.Instance_Not_Execute_Block_Number_Evaluation()
	//src.Baseline_Not_Execute_Block_Number_Evaluation()
	src.Baseline_Waiting_Time_Evalutation()
}
func main() {
	rand.Seed(time.Now().UnixNano())
	//test()
	eval()
}
