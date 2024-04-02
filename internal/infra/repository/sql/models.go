package sql

type Edge struct {
	ID         uint   `gorm:"primarykey"`
	AllColumns string `gorm:"index"`
	ObjNs      string `gorm:"index:object"`
	ObjName    string `gorm:"index:object"`
	ObjRel     string `gorm:"index:object"`
	SbjNs      string `gorm:"index:subject"`
	SbjName    string `gorm:"index:subject"`
	SbjRel     string `gorm:"index:subject"`
}
