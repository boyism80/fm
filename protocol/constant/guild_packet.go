package constant

type GuildOperationCode uint8

const (
	GuildC2SCreate           GuildOperationCode = 0x02
	GuildC2SInvite           GuildOperationCode = 0x05
	GuildC2SAcceptInvite     GuildOperationCode = 0x06
	GuildC2SLeave            GuildOperationCode = 0x07
	GuildC2SExpel            GuildOperationCode = 0x08
	GuildC2SChangeRankTitles GuildOperationCode = 0x0D
	GuildC2SChangeMemberRank GuildOperationCode = 0x0E
	GuildC2SChangeEmblem     GuildOperationCode = 0x0F
	GuildC2SChangeNotice     GuildOperationCode = 0x10
	GuildC2SPurchaseSkill    GuildOperationCode = 0x1D
	GuildC2SActivateSkill    GuildOperationCode = 0x1E
	GuildC2SChangeLeader     GuildOperationCode = 0x1F
)

type GuildSubOpcode uint8

const (
	GuildS2CInvite                 GuildSubOpcode = 0x05
	GuildS2CShowInfo               GuildSubOpcode = 0x1A
	GuildS2CNewMember              GuildSubOpcode = 0x27
	GuildS2CMemberLeft             GuildSubOpcode = 0x2C
	GuildS2CMemberExpelled         GuildSubOpcode = 0x2F
	GuildS2CDisband                GuildSubOpcode = 0x32
	GuildS2CCapacityChange         GuildSubOpcode = 0x3A
	GuildS2CMemberLevelClassUpdate GuildSubOpcode = 0x3C
	GuildS2CMemberOnline           GuildSubOpcode = 0x3D
	GuildS2CRankTitleChange        GuildSubOpcode = 0x3E
	GuildS2CChangeRank             GuildSubOpcode = 0x40
	GuildS2CEmblemChange           GuildSubOpcode = 0x42
	GuildS2CNotice                 GuildSubOpcode = 0x44
	GuildS2CUpdateGP               GuildSubOpcode = 0x48
	GuildS2CShowRanks              GuildSubOpcode = 0x49
)

type GuildResponseCode uint8

const (
	GuildResponseCreateDialog   GuildResponseCode = 0x01
	GuildResponseEmblemDialog   GuildResponseCode = 0x11
	GuildResponseAlreadyInGuild GuildResponseCode = 0x28
	GuildResponseNotInChannel   GuildResponseCode = 0x2A
	GuildResponseNotInGuild     GuildResponseCode = 0x2D
)
