package game

type WordMeta struct {
	RowStart    int
	ColumnStart int
	RowEnd      int
	ColumnEnd   int
	Points      int
}

type Word struct {
	WordMeta WordMeta
	PlayerID int
	Content  string
	Points   int
}
