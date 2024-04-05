package redis

import "github.com/skyrocketOoO/AuthNet/domain"

func vertexToPattern(vertex domain.Vertex) string {
	if vertex == (domain.Vertex{}) {
		return "*"
	}
	res := ""
	if vertex.Ns == "" {
		res += "*"
	} else {
		res += vertex.Ns
	}
	res += "%"
	if vertex.Name == "" {
		res += "*"
	} else {
		res += vertex.Name
	}
	res += "%"
	if vertex.Rel == "" {
		res += "*"
	} else {
		res += vertex.Rel
	}
	return res
}

func vertexToString(v domain.Vertex) string {
	return v.Ns + "%" + v.Name + "%" + v.Rel
}

func edgeToKeyValue(edge domain.Edge) (from string, to string) {
	from = vertexToString(domain.Vertex{
		Ns:   edge.SbjNs,
		Name: edge.SbjName,
		Rel:  edge.SbjRel,
	})
	to = vertexToString(domain.Vertex{
		Ns:   edge.ObjNs,
		Name: edge.ObjName,
		Rel:  edge.ObjRel,
	})
	return
}
