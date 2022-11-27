package telegram

type command string

const (
	COMMAND_UNKNOWN command = "unknown"

	COMMAND_START        command = "/start"
	COMMAND_NEW_TARGET   command = "New target"
	COMMAND_LIST_TARGETS command = "List targets"
	COMMAND_BACK         command = "Back"
)

func (c command) IsOk() bool {
	return c != COMMAND_UNKNOWN
}

func (c command) String() string {
	return string(c)
}

func (c command) IsReset() bool {
	return c == COMMAND_START || c == COMMAND_BACK
}

func ParseCommand(text string) command {
	switch text {
	case COMMAND_START.String():
		return COMMAND_START
	case COMMAND_NEW_TARGET.String():
		return COMMAND_NEW_TARGET
	case COMMAND_LIST_TARGETS.String():
		return COMMAND_LIST_TARGETS
	case COMMAND_BACK.String():
		return COMMAND_BACK
	}

	return COMMAND_UNKNOWN
}
