package database

import (
	"reflect"
	"testing"

	"github.com/pagarme/teleport/action"
	"github.com/pagarme/teleport/config"
)

var event *Event

func setupDml() {
	db = New(config.Database{
		Options: map[string]string{
			"sslmode": "disable",
		},
	})

	columns := []*Column{
		&Column{
			"id",
			1,
			"int4",
			"pg_catalog",
			"0",
			true,
			false,
			nil,
		},
		&Column{
			"content",
			1,
			"text",
			"pg_catalog",
			"0",
			false,
			false,
			nil,
		},
	}

	classes := []*Table{
		&Table{
			"123",
			"r",
			"test_table",
			columns,
			[]*Index{},
			nil,
		},
	}

	db.Schemas = map[string]*Schema{
		"public": &Schema{
			"123",
			"public",
			classes,
			nil,
			nil,
			nil,
			db,
		},
	}

	event = &Event{
		Id:            "1",
		Kind:          "dml",
		Status:        "waiting_batch",
		TriggerTag:    "public.test_table",
		TriggerEvent:  "UPDATE",
		TransactionId: "0",
		Data:          nil,
	}
}

func TestDmlDiff(t *testing.T) {
	setupDml()

	var tests = []struct {
		dml    *Dml
		output []action.Action
	}{
		{
			// Test insert a row
			&Dml{
				nil,
				&map[string]interface{}{
					"id":      5,
					"content": "test",
				},
				event,
				db,
				"target_schema",
			},
			[]action.Action{
				&action.InsertRow{
					"target_schema",
					"test_table",
					"id",
					[]action.Row{
						action.Row{
							"test",
							action.Column{
								"content",
								"text",
								true,
								false,
							},
						},
						action.Row{
							5,
							action.Column{
								"id",
								"int4",
								true,
								false,
							},
						},
					},
					false,
				},
			},
		},
		{
			// Test delete a row
			&Dml{
				&map[string]interface{}{
					"id":      5,
					"content": "test",
				},
				nil,
				event,
				db,
				"target_schema",
			},
			[]action.Action{
				&action.DeleteRow{
					"target_schema",
					"test_table",
					action.Row{
						5,
						action.Column{
							"id",
							"int4",
							true,
							false,
						},
					},
				},
			},
		},
		{
			// Test update a row
			&Dml{
				&map[string]interface{}{
					"id":      5,
					"content": "test",
				},
				&map[string]interface{}{
					"id":      5,
					"content": "test updated",
				},
				event,
				db,
				"target_schema",
			},
			[]action.Action{
				&action.UpdateRow{
					"target_schema",
					"test_table",
					action.Row{
						5,
						action.Column{
							"id",
							"int4",
							true,
							false,
						},
					},
					[]action.Row{
						action.Row{
							"test updated",
							action.Column{
								"content",
								"text",
								true,
								false,
							},
						},
						action.Row{
							5,
							action.Column{
								"id",
								"int4",
								true,
								false,
							},
						},
					},
				},
			},
		},
		{
			// Test insert a row with content not present in table
			&Dml{
				nil,
				&map[string]interface{}{
					"id":           5,
					"content":      "test",
					"other_column": "something is wrong",
				},
				event,
				db,
				"target_schema",
			},
			[]action.Action{
				&action.InsertRow{
					"target_schema",
					"test_table",
					"id",
					[]action.Row{
						action.Row{
							"test",
							action.Column{
								"content",
								"text",
								true,
								false,
							},
						},
						action.Row{
							5,
							action.Column{
								"id",
								"int4",
								true,
								false,
							},
						},
					},
					false,
				},
			},
		},
	}

	for _, test := range tests {
		actions := test.dml.Diff()

		if !reflect.DeepEqual(actions, test.output) {
			t.Errorf(
				"diff %#v => %#v, want %#v",
				test.dml,
				actions[0],
				test.output[0],
			)
		}
	}
}
