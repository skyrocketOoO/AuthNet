package domain

type Action string

const (
	CreateOperation           Action = "create"
	DeleteOperation           Action = "delete"
	CreateIfNotExistOperation Action = "createIfNotExists"
)

type Operation struct {
	Type Action `json:"action"`
	Edge Edge   `json:"edge"`
}
