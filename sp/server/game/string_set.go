package game

type WordMetaSet struct {
	List map[WordMeta]struct{}
}

func (s *WordMetaSet) Has(v WordMeta) bool {
	_, ok := s.List[v]
	return ok
}

func (s *WordMetaSet) Add(v WordMeta) {
	s.List[v] = struct{}{}
}

func (s *WordMetaSet) Remove(v WordMeta) {
	delete(s.List, v)
}

func (s *WordMetaSet) Clear() {
	s.List = make(map[WordMeta]struct{})
}

func NewWordMetaSet() *WordMetaSet {
	s := &WordMetaSet{}
	s.List = make(map[WordMeta]struct{})
	return s
}
