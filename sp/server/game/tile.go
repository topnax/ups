package game

type Tile struct {
	Row         int    `json:"row"`
	Column      int    `json:"column"`
	Set         bool   `json:"set"`
	Highlighted bool   `json:"highlighted"`
	Type        int    `json:"type"`
	Letter      Letter `json:"letter"`
}

func (tile Tile) getWordMultiplicand() int {
	wordMultiplicand := 1
	if tile.Type == MULTIPLY_WORD_2 {
		wordMultiplicand = 2
	} else if tile.Type == MULTIPLY_WORD_3 {
		wordMultiplicand = 3
	}
	return wordMultiplicand
}

func (tile Tile) getTilePoints() int {
	tilePoints := tile.Letter.Points
	tileMultiplicand := 1
	if tile.Type == MULTIPLY_LETTER_2 {
		tileMultiplicand = 2
	} else if tile.Type == MULTIPLY_LETTER_3 {
		tileMultiplicand = 3
	}
	tilePoints *= tileMultiplicand
	return tilePoints
}
