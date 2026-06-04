package constant

const AllianceInviteTargetNotOnlineMessage = "대상 길드의 길드장이 접속중이며, 현재 채널에 있어야 합니다."

const AllianceInviteGuildNotFoundMessage = "해당 길드가 존재하지 않습니다."

const AllianceInviteDeniedMessageSuffix = " 길드가 길드 연합 초대를 거절하였습니다."

const AllianceNoticeChangedMessagePrefix = "길드 연합 공지사항 : "

const AllianceLeaderChangedMessageSuffix = " 님이 새로운 길드 연합장이 되었습니다."

const AllianceCreateMesoCost int32 = 5_000_000

const AllianceIncreaseCapacityMesoCost int32 = 5_000_000

const AllianceCapacityMax uint32 = 5

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

type AllianceIncreaseCapacityResult int

const (
	AllianceIncreaseCapacityResultOK               AllianceIncreaseCapacityResult = 0
	AllianceIncreaseCapacityResultNotInAlliance    AllianceIncreaseCapacityResult = 1
	AllianceIncreaseCapacityResultNotLeader        AllianceIncreaseCapacityResult = 2
	AllianceIncreaseCapacityResultNotGuildMaster   AllianceIncreaseCapacityResult = 3
	AllianceIncreaseCapacityResultInsufficientMeso AllianceIncreaseCapacityResult = 4
	AllianceIncreaseCapacityResultCapacityMax      AllianceIncreaseCapacityResult = 5
	AllianceIncreaseCapacityResultFailed           AllianceIncreaseCapacityResult = 6
	AllianceIncreaseCapacityResultSendFailed       AllianceIncreaseCapacityResult = -1
)

type AllianceDisbandResult int

const (
	AllianceDisbandResultOK             AllianceDisbandResult = 0
	AllianceDisbandResultNotInAlliance  AllianceDisbandResult = 1
	AllianceDisbandResultNotLeader      AllianceDisbandResult = 2
	AllianceDisbandResultNotGuildMaster AllianceDisbandResult = 3
	AllianceDisbandResultFailed         AllianceDisbandResult = 4
	AllianceDisbandResultSendFailed     AllianceDisbandResult = -1
)

type AllianceLeaveResult int

const (
	AllianceLeaveResultOK             AllianceLeaveResult = 0
	AllianceLeaveResultNotInAlliance  AllianceLeaveResult = 1
	AllianceLeaveResultNotGuildMaster AllianceLeaveResult = 2
	AllianceLeaveResultRankForbidden  AllianceLeaveResult = 3
	AllianceLeaveResultFailed         AllianceLeaveResult = 4
	AllianceLeaveResultSendFailed     AllianceLeaveResult = -1
)
