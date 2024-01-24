//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var ProjectMembers = newProjectMembersTable("public", "project_members", "")

type projectMembersTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnString
	ProjectID   postgres.ColumnString
	UserID      postgres.ColumnString
	IsOwner     postgres.ColumnBool
	IsActive    postgres.ColumnBool
	Permissions postgres.ColumnString
	Metadata    postgres.ColumnString
	IsDeleted   postgres.ColumnBool
	CreatedAt   postgres.ColumnTimestampz
	UpdatedAt   postgres.ColumnTimestampz
	DeletedAt   postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ProjectMembersTable struct {
	projectMembersTable

	EXCLUDED projectMembersTable
}

// AS creates new ProjectMembersTable with assigned alias
func (a ProjectMembersTable) AS(alias string) *ProjectMembersTable {
	return newProjectMembersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ProjectMembersTable with assigned schema name
func (a ProjectMembersTable) FromSchema(schemaName string) *ProjectMembersTable {
	return newProjectMembersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ProjectMembersTable with assigned table prefix
func (a ProjectMembersTable) WithPrefix(prefix string) *ProjectMembersTable {
	return newProjectMembersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ProjectMembersTable with assigned table suffix
func (a ProjectMembersTable) WithSuffix(suffix string) *ProjectMembersTable {
	return newProjectMembersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newProjectMembersTable(schemaName, tableName, alias string) *ProjectMembersTable {
	return &ProjectMembersTable{
		projectMembersTable: newProjectMembersTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newProjectMembersTableImpl("", "excluded", ""),
	}
}

func newProjectMembersTableImpl(schemaName, tableName, alias string) projectMembersTable {
	var (
		IDColumn          = postgres.StringColumn("id")
		ProjectIDColumn   = postgres.StringColumn("project_id")
		UserIDColumn      = postgres.StringColumn("user_id")
		IsOwnerColumn     = postgres.BoolColumn("is_owner")
		IsActiveColumn    = postgres.BoolColumn("is_active")
		PermissionsColumn = postgres.StringColumn("permissions")
		MetadataColumn    = postgres.StringColumn("metadata")
		IsDeletedColumn   = postgres.BoolColumn("is_deleted")
		CreatedAtColumn   = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn   = postgres.TimestampzColumn("updated_at")
		DeletedAtColumn   = postgres.TimestampzColumn("deleted_at")
		allColumns        = postgres.ColumnList{IDColumn, ProjectIDColumn, UserIDColumn, IsOwnerColumn, IsActiveColumn, PermissionsColumn, MetadataColumn, IsDeletedColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn}
		mutableColumns    = postgres.ColumnList{ProjectIDColumn, UserIDColumn, IsOwnerColumn, IsActiveColumn, PermissionsColumn, MetadataColumn, IsDeletedColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn}
	)

	return projectMembersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		ProjectID:   ProjectIDColumn,
		UserID:      UserIDColumn,
		IsOwner:     IsOwnerColumn,
		IsActive:    IsActiveColumn,
		Permissions: PermissionsColumn,
		Metadata:    MetadataColumn,
		IsDeleted:   IsDeletedColumn,
		CreatedAt:   CreatedAtColumn,
		UpdatedAt:   UpdatedAtColumn,
		DeletedAt:   DeletedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
