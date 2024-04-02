package cmd

import errors "github.com/rotisserie/eris"

type DatabaseEnum string

const (
	databaseEnumPg     DatabaseEnum = "pg"
	databaseEnumSqlite DatabaseEnum = "sqlite"
	databaseEnumMongo  DatabaseEnum = "mongo"
	databaseEnumRedis  DatabaseEnum = "redis"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *DatabaseEnum) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *DatabaseEnum) Set(v string) error {
	switch v {
	case string(databaseEnumPg), string(databaseEnumSqlite),
		string(databaseEnumMongo), string(databaseEnumRedis):
		*e = DatabaseEnum(v)
		return nil
	default:
		return errors.New(`must be one of "pg", "sqlite", "mongo", "redis"`)
	}
}

// Type is only used in help text
func (e *DatabaseEnum) Type() string {
	return "DatabaseEnum"
}
