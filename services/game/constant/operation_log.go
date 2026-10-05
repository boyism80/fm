package constant

type OperationLogKind string

const (
	OperationLogGuildCreateUnsettled         OperationLogKind = "guild_create_unsettled"
	OperationLogGuildIncCapacityUnsettled    OperationLogKind = "guild_inc_capacity_unsettled"
	OperationLogAllianceCreateUnsettled      OperationLogKind = "alliance_create_unsettled"
	OperationLogAllianceIncCapacityUnsettled OperationLogKind = "alliance_inc_capacity_unsettled"
)
