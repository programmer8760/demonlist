package domain

const (
	PermChangePlacements uint16 = 1 << iota
	PermChangeRecordStatus
	PermEditLevels
	PermAddLevels
	PermDeleteLevels

	PermAll uint16 = 1<<iota - 1
)

func HasPermissions(perms uint16, required uint16) bool {
	return perms&required == required
}
