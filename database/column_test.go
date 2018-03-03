package database

import (
	"reflect"
	"testing"

	"github.com/pagarme/teleport/action"
	"github.com/pagarme/teleport/batcher/ddldiff"
)

var class *Table

var defaultContext ddldiff.Context

func init() {
	schema = &Schema{
		"123",
		"test_schema",
		[]*Table{},
		nil,
		nil,
		nil,
		nil,
	}

	class = &Table{
		"123",
		"r",
		"test_table",
		[]*Column{},
		[]*Index{},
		schema,
	}

	defaultContext = ddldiff.Context{
		Schema: "default_context",
	}
}

func TestColumnDiff(t *testing.T) {
	var tests = []struct {
		pre    *Column
		post   *Column
		output []action.Action
	}{
		{
			// Diff a column creation
			nil,
			&Column{
				"test_col",
				1,
				"text",
				"pg_catalog",
				"0",
				false,
				false,
				class,
			},
			[]action.Action{
				&action.CreateColumn{
					"default_context",
					"test_table",
					action.Column{
						"test_col",
						"text",
						true,
						false,
					},
				},
			},
		},
		{
			// Diff a column update
			&Column{
				"test_col",
				1,
				"text",
				"pg_catalog",
				"0",
				false,
				false,
				class,
			},
			&Column{
				"test_col_2",
				1,
				"int4",
				"pg_catalog",
				"0",
				false,
				false,
				class,
			},
			[]action.Action{
				&action.AlterColumn{
					"default_context",
					"test_table",
					action.Column{
						"test_col",
						"text",
						true,
						false,
					},
					action.Column{
						"test_col_2",
						"int4",
						true,
						false,
					},
				},
			},
		},
	}

	for _, test := range tests {
		// Avoid passing a interface with nil pointer
		// to Diff and breaking comparisons with nil.
		var preObj ddldiff.Diffable
		if test.pre == nil {
			preObj = nil
		} else {
			preObj = test.pre
		}

		actions := test.post.Diff(preObj, defaultContext)

		if !reflect.DeepEqual(actions, test.output) {
			t.Errorf(
				"diff %#v with %#v => %v, want %d",
				test.pre,
				test.post,
				actions,
				test.output,
			)
		}
	}
}

func TestColumnChildren(t *testing.T) {
	attr := &Column{
		"test_col",
		1,
		"text",
		"pg_catalog",
		"0",
		false,
		false,
		class,
	}

	children := attr.Children()

	if len(children) != 0 {
		t.Errorf("attr children => %d, want %d", len(children), 0)
	}
}

func TestColumnDrop(t *testing.T) {
	attr := &Column{
		"test_col",
		1,
		"text",
		"pg_catalog",
		"0",
		false,
		false,
		class,
	}

	actions := attr.Drop(defaultContext)

	if len(actions) != 1 {
		t.Errorf("actions => %d, want %d", len(actions), 1)
	}

	dropAction, ok := actions[0].(*action.DropColumn)

	if !ok {
		t.Errorf("action is not DropColumn")
	}

	if dropAction.SchemaName != defaultContext.Schema {
		t.Errorf("drop action schema name => %s, want %s", dropAction.SchemaName, defaultContext.Schema)
	}

	if dropAction.TableName != class.RelationName {
		t.Errorf("drop action table name => %s, want %s", dropAction.TableName, class.RelationName)
	}

	if dropAction.Column.Name != attr.Name {
		t.Errorf("drop action column name => %s, want %s", dropAction.Column.Name, attr.Name)
	}

	if dropAction.Column.Type != attr.TypeName {
		t.Errorf("drop action column name => %s, want %s", dropAction.Column.Type, attr.TypeName)
	}
}

func TestColumnIsEqual(t *testing.T) {
	pre := &Column{
		"test_col",
		1,
		"text",
		"pg_catalog",
		"0",
		false,
		false,
		class,
	}

	post := &Column{
		"test_col_2",
		1,
		"int4",
		"pg_catalog",
		"0",
		false,
		false,
		class,
	}

	if !post.IsEqual(pre) {
		t.Errorf("expect classes to be equal")
	}

	post.Name = pre.Name
	post.Num = 2

	if post.IsEqual(pre) {
		t.Errorf("expect classes not to be equal")
	}

	preOtherType := &Index{
		"123",
		"test_index_2",
		"create index",
		class,
	}

	if post.IsEqual(preOtherType) {
		t.Errorf("expect two different types not to be equal")
	}
}
