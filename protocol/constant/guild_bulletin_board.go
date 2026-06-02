package constant

type GuildBulletinBoardC2SAction uint8

const (
	GuildBulletinBoardC2SWriteThread  GuildBulletinBoardC2SAction = 0
	GuildBulletinBoardC2SDeleteThread GuildBulletinBoardC2SAction = 1
	GuildBulletinBoardC2SListThreads  GuildBulletinBoardC2SAction = 2
	GuildBulletinBoardC2SShowThread   GuildBulletinBoardC2SAction = 3
	GuildBulletinBoardC2SWriteReply   GuildBulletinBoardC2SAction = 4
	GuildBulletinBoardC2SDeleteReply  GuildBulletinBoardC2SAction = 5
)

type GuildBulletinBoardS2CAction uint8

const (
	GuildBulletinBoardS2CThreadList GuildBulletinBoardS2CAction = 6
	GuildBulletinBoardS2CShowThread GuildBulletinBoardS2CAction = 7
)

const (
	GuildBulletinBoardTitleMaxLen         = 25
	GuildBulletinBoardBodyMaxLen          = 600
	GuildBulletinBoardReplyMaxLen         = 25
	GuildBulletinBoardThreadsPerPage      = 10
	GuildBulletinBoardNoticeLocalThreadID = 0
	GuildBulletinBoardIconCashMin         = 0x64
	GuildBulletinBoardIconCashMax         = 0x6a
)
