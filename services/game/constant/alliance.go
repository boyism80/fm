package constant

const AllianceCreateMesoCost int32 = 5_000_000

type AllianceCreateResult int

const (
	AllianceCreateResultOK                  AllianceCreateResult = 0
	AllianceCreateResultInvalidRequirements AllianceCreateResult = 1
	AllianceCreateResultInsufficientMeso    AllianceCreateResult = 2
	AllianceCreateResultInvalidName         AllianceCreateResult = 3
	AllianceCreateResultNameTaken           AllianceCreateResult = 4
	AllianceCreateResultFailed              AllianceCreateResult = 5
	AllianceCreateResultSendFailed          AllianceCreateResult = -1
)
