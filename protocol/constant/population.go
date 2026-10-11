package constant

type PopulationMode uint8

const (
	PopulationModeDown PopulationMode = 0
	PopulationModeUp   PopulationMode = 1
)

// PopulationResult values are the client's population result codes.
type PopulationResult uint8

const (
	PopulationResultGiven            PopulationResult = 0
	PopulationResultWrongName        PopulationResult = 1
	PopulationResultLevelTooLow      PopulationResult = 2
	PopulationResultAlreadyToday     PopulationResult = 3
	PopulationResultAlreadyThisMonth PopulationResult = 4
	PopulationResultReceived         PopulationResult = 5
	PopulationResultFailed           PopulationResult = 6
)
