package game

type PlayerSet struct {
	List map[Player]struct{}
}

func (s *PlayerSet) Has(v Player) bool {
	_, ok := s.List[v]
	return ok
}

func (s *PlayerSet) Add(v Player) {
	s.List[v] = struct{}{}
}

func (s *PlayerSet) Remove(v Player) {
	delete(s.List, v)
}

func (s *PlayerSet) Clear() {
	s.List = make(map[Player]struct{})
}

func NewPlayerSet() *PlayerSet {
	s := &PlayerSet{}
	s.List = make(map[Player]struct{})
	return s
}
