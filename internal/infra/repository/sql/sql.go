package sql

import (
	"context"
	"strings"

	errors "github.com/rotisserie/eris"
	"github.com/skyrocketOoO/AuthNet/domain"
	"gorm.io/gorm"
)

type SqlRepository struct {
	db *gorm.DB
}

func NewSqlRepository(db *gorm.DB) (*SqlRepository, error) {
	return &SqlRepository{
		db: db,
	}, nil
}

func (r *SqlRepository) Ping(c context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return errors.New(err.Error())
	}
	if err := db.PingContext(c); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *SqlRepository) Get(c context.Context, edge domain.Edge,
	queryMode bool) ([]domain.Edge, error) {
	if queryMode {
		var sqlEdges []Edge
		if edge == (domain.Edge{}) {
			if err := r.db.WithContext(c).Find(&sqlEdges).Error; err != nil {
				return nil, errors.New(err.Error())
			}
		} else {
			if err := r.db.WithContext(c).Where(&edge).Find(&sqlEdges).
				Error; err != nil {
				return nil, errors.New(err.Error())
			}
		}
		edges := make([]domain.Edge, len(sqlEdges))
		for i, sqlEdge := range sqlEdges {
			edges[i] = convertToRel(sqlEdge)
		}
		return edges, nil
	} else {
		var sqlEdge Edge
		if err := r.db.WithContext(c).Where("all_columns = ?",
			concatAttr(edge)).Take(&sqlEdge).Error; err != nil {
			return nil, errors.New(err.Error())
		}
		return []domain.Edge{convertToRel(sqlEdge)}, nil
	}
}

func (r *SqlRepository) Create(c context.Context, edge domain.Edge) error {
	sqlRel := convertToSqlModel(edge)
	if err := r.db.WithContext(c).Create(&sqlRel).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (r *SqlRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	if queryMode {
		if err := r.db.WithContext(c).Where(&edge).Delete(&Edge{}).
			Error; err != nil {
			return errors.New(err.Error())
		}
	} else {
		if _, err := r.Get(c, edge, false); err != nil {
			return err
		}
		if err := r.db.Where("all_columns = ?", concatAttr(edge)).
			Delete(&Edge{}).Error; err != nil {
			return errors.New(err.Error())
		}
	}
	return nil
}

// func (r *SqlRepository) BatchOperation(c context.Context, operations []domain.Operation) error {
// 	tx := r.db.Begin()
// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	for _, operation := range operations {
// 		switch operation.Type {
// 		case domain.CreateOperation:
// 			if err := r.Create(operation.Edge); err != nil {
// 				tx.Rollback()
// 				return err
// 			}
// 		case domain.DeleteOperation:
// 			if err := r.Delete(operation.Edge); err != nil {
// 				tx.Rollback()
// 				return err
// 			}
// 		case domain.CreateIfNotExistOperation:
// 			if err := r.Create(operation.Edge); err != nil {
// 				if err != gorm.ErrDuplicatedKey {
// 					tx.Rollback()
// 					return err
// 				}
// 			}
// 		default:
// 			tx.Rollback()
// 			return errors.New("invalid operation type")
// 		}
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *SqlRepository) GetAllNs(c context.Context) ([]string, error) {
// 	sqlQuery := `
// 		SELECT DISTINCT namespace
// 		FROM (
// 			SELECT object_namespace AS namespace FROM edges
// 			UNION
// 			SELECT subject_namespace AS namespace FROM edges
// 		) AS namespaces
// 	`
// 	var nss []string
// 	if err := r.db.Raw(sqlQuery).Scan(&nss).Error; err != nil {
// 		return nil, err
// 	}

// 	return nss, nil
// }

func (r *SqlRepository) ClearAll(c context.Context) error {
	query := "DELETE FROM edges"
	if err := r.db.WithContext(c).Exec(query).Error; err != nil {
		return err
	}
	return nil
}

func convertToSqlModel(edge domain.Edge) Edge {
	return Edge{
		ObjNs:      edge.ObjNs,
		ObjName:    edge.ObjName,
		ObjRel:     edge.ObjRel,
		SbjNs:      edge.SbjNs,
		SbjName:    edge.SbjName,
		SbjRel:     edge.SbjRel,
		AllColumns: concatAttr(edge),
	}
}

func concatAttr(edge domain.Edge) string {
	return strings.Join(
		[]string{
			edge.ObjNs,
			edge.ObjName,
			edge.ObjRel,
			edge.SbjNs,
			edge.SbjName,
			edge.SbjRel,
		},
		"%",
	)
}

func convertToRel(edge Edge) domain.Edge {
	return domain.Edge{
		ObjNs:   edge.ObjNs,
		ObjName: edge.ObjName,
		ObjRel:  edge.ObjRel,
		SbjNs:   edge.SbjNs,
		SbjName: edge.SbjName,
		SbjRel:  edge.SbjRel,
	}
}
