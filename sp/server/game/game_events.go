package game

const (
	SET_AT       = 0
	RESET_AT     = 1
	SUBMIT_WORD  = 2
	CONFIRM_WORD = 3
)

type SetLetterAtEvent struct {
	PlayerID int
	Row      int
	Column   int
	Letter   string
}

type ResetAtEvent struct {
	PlayerID int
	Row      int
	Column   int
}

type SubmitWordEvent struct {
	PlayerID int
}

type ConfirmWordEvent struct {
	PlayerID int
}
