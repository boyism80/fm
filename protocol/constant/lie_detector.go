package constant

type LieDetectorResult uint8

const (
	LieDetectorResultNotHunting LieDetectorResult = 1
	LieDetectorResultTested     LieDetectorResult = 2
	LieDetectorResultBusy       LieDetectorResult = 3
	LieDetectorResultQuestion   LieDetectorResult = 6
	LieDetectorResultCaught     LieDetectorResult = 7
	LieDetectorResultPassed     LieDetectorResult = 9
)
