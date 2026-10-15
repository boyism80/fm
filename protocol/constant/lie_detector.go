package constant

type LieDetectorResult uint8

const (
	LieDetectorResultNotFound   LieDetectorResult = 0
	LieDetectorResultNotHunting LieDetectorResult = 1
	LieDetectorResultTested     LieDetectorResult = 2
	LieDetectorResultBusy       LieDetectorResult = 3
	LieDetectorResultStarted    LieDetectorResult = 5
	LieDetectorResultQuestion   LieDetectorResult = 6
	LieDetectorResultCaught     LieDetectorResult = 7
	LieDetectorResultPassed     LieDetectorResult = 9
	LieDetectorResultReported   LieDetectorResult = 11
)
