package domain

type Edge struct {
	ObjNs   string `json:"obj_ns" bson:"obj_ns"`
	ObjName string `json:"obj_name" bson:"obj_name"`
	ObjRel  string `json:"obj_rel" bson:"obj_rel"`
	SbjNs   string `json:"sbj_ns" bson:"sbj_ns"`
	SbjName string `json:"sbj_name" bson:"sbj_name"`
	SbjRel  string `json:"sbj_rel" bson:"sbj_rel"`
}
type Vertex struct {
	Ns   string `json:"ns"`
	Name string `json:"name"`
	Rel  string `json:"rel"`
}

type TreeNode struct {
	Ns       string      `json:"ns"`
	Name     string      `json:"name"`
	Rel      string      `json:"rel"`
	Children []*TreeNode `json:"children"`
}

type PageOptions struct {
	LastID   uint
	PageSize int
}
