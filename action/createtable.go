package action

import (
	"encoding/gob"
	// "strings"
)

type CreateTable struct {
	SchemaName string
	TableName  string
	PrimaryKey Column
}

// Register type for gob
func init() {
	gob.Register(&CreateTable{})
}

func (a *CreateTable) Execute(c *Context) error {
	// cmd := fmt.Sprintf(
	// 	"CREATE TABLE IF NOT EXISTS \"%s\".\"%s\" (\"%s\" %s\"%s\"%s PRIMARY KEY);",
	// 	a.SchemaName,
	// 	a.TableName,
	// 	a.PrimaryKey.Name,
	// 	a.PrimaryKey.GetTypeSchemaStr(a.SchemaName),
	// 	a.PrimaryKey.Type,
	// 	a.PrimaryKey.DefaultStatement(),
	// )
	// _, err := c.Tx.Exec(cmd)
	// log.Println(cmd)

	// return err
	return nil
}

func (a *CreateTable) Filter(targetExpression string) bool {
	return IsInTargetExpression(&targetExpression, &a.SchemaName, &a.TableName)
}

func (a *CreateTable) NeedsSeparatedBatch() bool {
	return false
}
