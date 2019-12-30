package game

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
