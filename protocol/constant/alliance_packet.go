package constant

type AllianceC2SOperation uint8

const (
	AllianceC2SLoadInfo         AllianceC2SOperation = 0x01
	AllianceC2SLeave            AllianceC2SOperation = 0x02
	AllianceC2SCreate           AllianceC2SOperation = 0x05
	AllianceC2SInvite           AllianceC2SOperation = 0x03
	AllianceC2SAcceptInvite     AllianceC2SOperation = 0x04
	AllianceC2SExpel            AllianceC2SOperation = 0x06
	AllianceC2SChangeLeader     AllianceC2SOperation = 0x07
	AllianceC2SChangeRankTitles AllianceC2SOperation = 0x08
	AllianceC2SChangeMemberRank AllianceC2SOperation = 0x09
	AllianceC2SChangeNotice     AllianceC2SOperation = 0x0A
	AllianceC2SDenyInvite       AllianceC2SOperation = 0x16
)

type AllianceS2CSubOpcode uint8

const (
	AllianceS2CChangeMembership   AllianceS2CSubOpcode = 0x01
	AllianceS2CChangeLeader       AllianceS2CSubOpcode = 0x02
	AllianceS2CInvite             AllianceS2CSubOpcode = 0x03
	AllianceS2CChangeGuildMembers AllianceS2CSubOpcode = 0x04
	AllianceS2CChangeMemberRank   AllianceS2CSubOpcode = 0x05
	AllianceS2CShowInfo           AllianceS2CSubOpcode = 0x0C
	AllianceS2CShowGuilds         AllianceS2CSubOpcode = 0x0D
	AllianceS2CMemberOnline       AllianceS2CSubOpcode = 0x0E
	AllianceS2CCreate             AllianceS2CSubOpcode = 0x0F
	AllianceS2CRemoveGuild        AllianceS2CSubOpcode = 0x10
	AllianceS2CAddGuild           AllianceS2CSubOpcode = 0x12
	AllianceS2CUpdateInfo         AllianceS2CSubOpcode = 0x17
	AllianceS2CUpdateMember       AllianceS2CSubOpcode = 0x18
	AllianceS2CUpdateLeader       AllianceS2CSubOpcode = 0x19
	AllianceS2CUpdateMemberRank   AllianceS2CSubOpcode = 0x1B
	AllianceS2CDisband            AllianceS2CSubOpcode = 0x1D
)
