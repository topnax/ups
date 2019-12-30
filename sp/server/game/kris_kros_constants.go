package game

// tile types
const (
	BASIC             = 0
	MULTIPLY_WORD_2   = 1
	MULTIPLY_WORD_3   = 2
	MULTIPLY_LETTER_2 = 3
	MULTIPLY_LETTER_3 = 4
)

// points of letters and their occurrence
func GetLetterPointsTable() map[string][2]int {
	return map[string][2]int{
		"a": {1, 6},
		"á": {2, 2},
		"b": {2, 2},
		"c": {3, 2},
		"č": {4, 2},

		"d": {2, 2},
		"ď": {8, 1},
		"e": {1, 5},
		"é": {5, 1},
		"ě": {5, 2},

		"f":  {8, 1},
		"g":  {8, 1},
		"h":  {3, 2},
		"ch": {4, 2},
		"i":  {1, 4},

		"í": {2, 2},
		"j": {2, 2},
		"k": {1, 4},
		"l": {1, 4},
		"m": {2, 3},

		"n": {1, 3},
		"ň": {6, 1},
		"o": {1, 6},
		"ó": {10, 1},
		"p": {1, 3},

		"r": {1, 4},
		"ř": {4, 2},
		"s": {1, 5},
		"š": {3, 2},
		"t": {1, 4},

		"ť": {6, 1},
		"u": {2, 3},
		"ů": {5, 1},
		"ú": {6, 1},
		"v": {1, 3},

		"x": {10, 2},
		"y": {1, 3},
		"ý": {4, 2},
		"z": {3, 2},
		"ž": {4, 2},
	}
}

// returns the kriskros desk represented as types of each tile
func GetDeskTileTypes() [15][15]int {
	return [15][15]int{
		{MULTIPLY_WORD_3, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_WORD_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, MULTIPLY_WORD_3},
		{BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC},
		{BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC},
		{MULTIPLY_LETTER_2, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, MULTIPLY_LETTER_2},
		{BASIC, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, BASIC},
		{BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC},
		{BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC},

		{MULTIPLY_WORD_3, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, BASIC, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, MULTIPLY_WORD_3},

		{BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC},
		{BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC},
		{BASIC, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, BASIC},
		{MULTIPLY_LETTER_2, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, MULTIPLY_LETTER_2},
		{BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC},
		{BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_3, BASIC, BASIC, BASIC, MULTIPLY_WORD_2, BASIC},
		{MULTIPLY_WORD_3, BASIC, BASIC, MULTIPLY_WORD_2, BASIC, BASIC, BASIC, MULTIPLY_WORD_3, BASIC, BASIC, BASIC, MULTIPLY_LETTER_2, BASIC, BASIC, MULTIPLY_WORD_3},
	}
}
