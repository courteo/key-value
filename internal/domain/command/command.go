package command

const (
	UnknownID = iota
	GetID
	SetID
	DeleteID
)

var (
	Unknown = "UNKNOWN"
	Set     = "SET"
	Get     = "GET"
	Del     = "DEL"
)

var namesToId = map[string]int{
	Set: SetID,
	Get: GetID,
	Del: DeleteID,
}

func commandNameToCommandID(command string) int {
	status, found := namesToId[command]
	if !found {
		return UnknownID
	}

	return status
}

const (
	setCommandArgumentsNumber = 2
	getCommandArgumentsNumber = 1
	delCommandArgumentsNumber = 1
)

var argumentsNumber = map[int]int{
	SetID:    setCommandArgumentsNumber,
	GetID:    getCommandArgumentsNumber,
	DeleteID: delCommandArgumentsNumber,
}

func commandArgumentsNumber(commandID int) int {
	return argumentsNumber[commandID]
}
