package sql

type Edge struct {
	ID      uint   `gorm:"primarykey"`
	ObjNs   string `gorm:"index:obj"`
	ObjName string `gorm:"index:obj"`
	ObjRel  string `gorm:"index:obj"`
	SbjNs   string `gorm:"index:sbj"`
	SbjName string `gorm:"index:sbj"`
	SbjRel  string `gorm:"index:sbj"`
}
