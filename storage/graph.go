package storage

type Graph struct {
	Keys []Pair
	Map  *BitMap
}

func (g *Graph) FindIndex(key []byte) (index uint64) {
	for _, p := range g.Keys {
		if p.MatchKey(key) {
			index = p.Index
			return
		}
	}
	return
}

func (g *Graph) FindKey(index uint64) (key []byte) {
	for _, p := range g.Keys {
		if p.MatchIndex(index) {
			key = p.Key
			return
		}
	}
	return
}

func (g *Graph) AddKey(key []byte) (index uint64) {
	if index = g.FindIndex(key); index == 0 {
		index = uint64(len(g.Keys)) + 1
		g.Keys = append(g.Keys, Pair{key, index})
	}
	return
}

func (g *Graph) FinalizeMap() {
	g.Map = NewBitMap(uint64(len(g.Keys)), uint64(len(g.Keys)))
}
