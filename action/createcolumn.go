package action

import (
	"encoding/gob"
	"fmt"
)

type CreateColumn struct {
	SchemaName string
	TableName  string
	Column     Column
}

// Register type for gob
func init() {
	gob.Register(&CreateColumn{})
}

func (a *CreateColumn) Execute(c *Context) error {
	_, err := c.Tx.Exec(
		fmt.Sprintf(
			"ALTER TABLE \"%s\".\"%s\" ADD COLUMN \"%s\" %s\"%s\"%s;",
			a.SchemaName,
			a.TableName,
			a.Column.Name,
			a.Column.GetTypeSchemaStr(a.SchemaName),
			a.Column.Type,
			a.notNullStatement(),
		),
	)

	return err
}

func (a *CreateColumn) notNullStatement() string {
	if a.Column.NotNull {
		return " NOT NULL"
	}

	return ""
}

func (a *CreateColumn) Filter(targetExpression string) bool {
	return true
}

func (a *CreateColumn) NeedsSeparatedBatch() bool {
	return false
}
