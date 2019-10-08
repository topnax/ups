package model

type LetterSet struct {
	List map[Tile]struct{} //empty structs occupy 0 memory
}

func (s *LetterSet) Has(v Tile) bool {
	_, ok := s.List[v]
	return ok
}

func (s *LetterSet) Add(v Tile) {
	s.List[v] = struct{}{}
}

func (s *LetterSet) Remove(v Tile) {
	delete(s.List, v)
}

func (s *LetterSet) Clear() {
	s.List = make(map[Tile]struct{})
}

func NewSet() *LetterSet {
	s := &LetterSet{}
	s.List = make(map[Tile]struct{})
	return s
}
