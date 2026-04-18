package constant

type PartyOperationCode uint8
type PartySubOpcode uint8
type PartyStatusCode byte

const (
	PartyC2SCreate       PartyOperationCode = 1
	PartyC2SLeave        PartyOperationCode = 2
	PartyC2SAcceptInvite PartyOperationCode = 3
	PartyC2SInvite       PartyOperationCode = 4
	PartyC2SExpel        PartyOperationCode = 5
	PartyC2SChangeLeader PartyOperationCode = 6
)

const (
	PartyS2CInvite       PartySubOpcode = 4
	PartyS2CPartyCreated PartySubOpcode = 8
	PartyS2CPartyUpdate  PartySubOpcode = 12
	PartyS2CPartyJoin    PartySubOpcode = 15
	PartyS2CSilentUpdate PartySubOpcode = 7
	PartyS2CLeaderChange PartySubOpcode = 26
	PartyS2CPartyPortal  PartySubOpcode = 34
)

const (
	PartyS2CStatusUnexpected     PartyStatusCode = 1
	PartyS2CStatusBeginnerCreate PartyStatusCode = 10
	PartyS2CStatusUnexpected11   PartyStatusCode = 11
	PartyS2CStatusNotInParty     PartyStatusCode = 13
	PartyS2CStatusUnexpected14   PartyStatusCode = 14
	PartyS2CStatusAlreadyInParty PartyStatusCode = 16
	PartyS2CStatusPartyFull      PartyStatusCode = 17
	PartyS2CStatusUnableFindChar PartyStatusCode = 18
	PartyS2CStatusUnexpected19   PartyStatusCode = 19
	PartyS2CStatusLeaderChange27 PartyStatusCode = 27
	PartyS2CStatusLeaderChange28 PartyStatusCode = 28
	PartyS2CStatusLeaderChange29 PartyStatusCode = 29
)

var partyStatusByOperationAndError = map[PartyOperationCode]map[int32]PartyStatusCode{
	PartyC2SCreate: {
		2: PartyS2CStatusAlreadyInParty,
	},
	PartyC2SLeave: {
		6: PartyS2CStatusNotInParty,
	},
	PartyC2SAcceptInvite: {
		5:  PartyS2CStatusPartyFull,
		11: PartyS2CStatusUnableFindChar,
		4:  PartyS2CStatusUnableFindChar,
		2:  PartyS2CStatusAlreadyInParty,
	},
	PartyC2SInvite: {
		5:  PartyS2CStatusPartyFull,
		12: PartyS2CStatusUnableFindChar,
		3:  PartyS2CStatusUnableFindChar,
		15: PartyS2CStatusAlreadyInParty,
		13: PartyS2CStatusUnableFindChar,
		14: PartyS2CStatusUnableFindChar,
		4:  PartyS2CStatusUnableFindChar,
	},
	PartyC2SExpel: {
		6: PartyS2CStatusNotInParty,
		7: PartyS2CStatusLeaderChange28,
		8: PartyS2CStatusUnableFindChar,
	},
	PartyC2SChangeLeader: {
		7: PartyS2CStatusLeaderChange28,
		8: PartyS2CStatusLeaderChange27,
		9: PartyS2CStatusLeaderChange29,
		4: PartyS2CStatusUnableFindChar,
	},
}

func PartyStatusForInternalError(op PartyOperationCode, errorCode int32) PartyStatusCode {
	if byErr, ok := partyStatusByOperationAndError[op]; ok {
		if status, exists := byErr[errorCode]; exists {
			return status
		}
	}
	return PartyS2CStatusUnexpected
}
