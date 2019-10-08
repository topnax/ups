package model

type LetterSet struct {
	List map[Letter]struct{} //empty structs occupy 0 memory
}

func (s *LetterSet) Has(v Letter) bool {
	_, ok := s.List[v]
	return ok
}

func (s *LetterSet) Add(v Letter) {
	s.List[v] = struct{}{}
}

func (s *LetterSet) Remove(v Letter) {
	delete(s.List, v)
}

func (s *LetterSet) Clear() {
	s.List = make(map[Letter]struct{})
}

func NewSet() *LetterSet {
	s := &LetterSet{}
	s.List = make(map[Letter]struct{})
	return s
}
