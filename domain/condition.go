package domain

type SearchCond struct {
	In Compare `json:"in"`
}

type Compare struct {
	Nses  []string `json:"nses"`
	Names []string `json:"names"`
	Rels  []string `json:"rels"`
}

func (c *SearchCond) ShouldStop(vertex Vertex) bool {
	if len(c.In.Nses) == 0 && len(c.In.Names) == 0 && len(c.In.Rels) == 0 {
		// means no specific condition, never stop
		return false
	}
	for _, ns := range c.In.Nses {
		if vertex.Ns == ns {
			return false
		}
	}
	for _, name := range c.In.Names {
		if vertex.Name == name {
			return false
		}
	}
	for _, rel := range c.In.Rels {
		if vertex.Rel == rel {
			return false
		}
	}
	return true
}

type CollectCond struct {
	In Compare `json:"in"`
}

func (c *CollectCond) ShouldCollect(vertex Vertex) bool {
	if len(c.In.Nses) == 0 && len(c.In.Names) == 0 && len(c.In.Rels) == 0 {
		// means no specific conditions, collect all vertexs
		return true
	}
	for _, ns := range c.In.Nses {
		if vertex.Ns == ns {
			return true
		}
	}
	for _, name := range c.In.Names {
		if vertex.Name == name {
			return true
		}
	}
	for _, rel := range c.In.Rels {
		if vertex.Rel == rel {
			return true
		}
	}
	return false
}
