package game

const (
	SET_AT = 0
	RESET_AT  = 1
	SUBMIT_WORD = 2
	CONFIRM_WORD = 3
)

type GameEvent struct {
	EventType int
}

type SetAtEvent struct {
	Event GameEvent
	PlayerID int
	Row int
	Column int
	Letter string
}

type ResetAtEvent struct {
	Event GameEvent
	PlayerID int
	Row int
	Column int
}

type SubmitWordEvent struct {
	Event GameEvent
	PlayerID int
}

type ConfirmWordEvent struct {
	Event GameEvent
	PlayerID int
}

