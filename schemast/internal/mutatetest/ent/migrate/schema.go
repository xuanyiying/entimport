// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// WithFieldsColumns holds the columns for the "with_fields" table.
	WithFieldsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "existing", Type: field.TypeString},
	}
	// WithFieldsTable holds the schema information for the "with_fields" table.
	WithFieldsTable = &schema.Table{
		Name:       "with_fields",
		Columns:    WithFieldsColumns,
		PrimaryKey: []*schema.Column{WithFieldsColumns[0]},
	}
	// WithModifiedFieldsColumns holds the columns for the "with_modified_fields" table.
	WithModifiedFieldsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Size: 10},
		{Name: "with_modified_field_owner", Type: field.TypeInt, Nullable: true},
	}
	// WithModifiedFieldsTable holds the schema information for the "with_modified_fields" table.
	WithModifiedFieldsTable = &schema.Table{
		Name:       "with_modified_fields",
		Columns:    WithModifiedFieldsColumns,
		PrimaryKey: []*schema.Column{WithModifiedFieldsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "with_modified_fields_users_owner",
				Columns:    []*schema.Column{WithModifiedFieldsColumns[2]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// WithNilFieldsColumns holds the columns for the "with_nil_fields" table.
	WithNilFieldsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
	}
	// WithNilFieldsTable holds the schema information for the "with_nil_fields" table.
	WithNilFieldsTable = &schema.Table{
		Name:       "with_nil_fields",
		Columns:    WithNilFieldsColumns,
		PrimaryKey: []*schema.Column{WithNilFieldsColumns[0]},
	}
	// WithoutFieldsColumns holds the columns for the "without_fields" table.
	WithoutFieldsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
	}
	// WithoutFieldsTable holds the schema information for the "without_fields" table.
	WithoutFieldsTable = &schema.Table{
		Name:       "without_fields",
		Columns:    WithoutFieldsColumns,
		PrimaryKey: []*schema.Column{WithoutFieldsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		UsersTable,
		WithFieldsTable,
		WithModifiedFieldsTable,
		WithNilFieldsTable,
		WithoutFieldsTable,
	}
)

func init() {
	WithModifiedFieldsTable.ForeignKeys[0].RefTable = UsersTable
}
