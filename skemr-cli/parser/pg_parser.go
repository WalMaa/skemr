package parser

import (
	"log/slog"

	pgquery "github.com/pganalyze/pg_query_go/v6"
)

type StatementAction struct {
	Target   string    // e.g., column name, table name, database name
	Action   SqlAction // The type of action performed (e.g., CREATE, DROP, ALTER)
	Relation string    // e.g., table name for column actions
}

type SqlAction string

const (
	// Database level actions
	SqlActionCreateDatabase SqlAction = "CREATE DATABASE"
	SqlActionRenameDatabase SqlAction = "RENAME DATABASE"
	SqlActionDropDatabase   SqlAction = "DROP DATABASE"
	// -- Name level actions
	SqlActionCreateTable SqlAction = "CREATE TABLE"
	SqlActionRenameTable SqlAction = "RENAME TABLE"
	SqlActionDropTable   SqlAction = "DROP TABLE"

	// ---- Column level actions
	SqlActionModifyDataType SqlAction = "MODIFY DATA TYPE"
	SqlActionRenameColumn   SqlAction = "RENAME COLUMN"
	SqlActionDropColumn     SqlAction = "DROP COLUMN"
	SqlActionAddColumn      SqlAction = "ADD COLUMN"

	SqlActionInsertRow SqlAction = "INSERT ROW"
	// Fallback
	SqlActionUndefined SqlAction = "UNDEFINED"
)

func parseStatement(stmt *pgquery.RawStmt) (StatementAction, error) {
	slog.Debug("Parsing", "statement", stmt.String())
	node := stmt.GetStmt()
	// Check for a DROP DATABASE statement
	if node.GetDropdbStmt() != nil {
		return parseDropDatabase(node)
	}

	if dropTableStmt := node.GetDropStmt(); dropTableStmt != nil {
		return parseDrop(dropTableStmt)
	}

	if alterTablestmt := node.GetAlterTableStmt(); alterTablestmt != nil {
		return parseAlterTable(alterTablestmt)
	}

	//Check for Rename column or table
	if renameStmt := node.GetRenameStmt(); renameStmt != nil {
		return parseRenameStmt(renameStmt)
	}

	if insertStmt := node.GetInsertStmt(); insertStmt != nil {
		return parseInsertStmt(insertStmt)
	}

	if createStmt := node.GetCreateStmt(); createStmt != nil {
		// Handle CREATE TABLE or other create statements if needed
		// For now, we will just return an undefined action
		return StatementAction{
			Action: SqlActionCreateTable,
		}, nil
	}

	if createDbStmt := node.GetCreatedbStmt(); createDbStmt != nil {
		return parseCreateDatabaseStmt(createDbStmt)
	}

	// If the statement is not recognized, return an undefined action
	slog.Warn("Unsupported statement type", "statement", stmt.String())
	return StatementAction{
		Target:   "",
		Action:   SqlActionUndefined,
		Relation: "",
	}, nil
}

/*
ParseSql parses a migration file and returns a structured representation of the SQL.
*/
func ParseSql(sql string) ([]StatementAction, error) {
	tree, err := pgquery.Parse(sql)
	result := make([]StatementAction, 0)
	if err != nil {
		slog.Error("Failed to parse SQL", "error", err, "sql", sql)

	}
	stmts := tree.Stmts

	for _, stmt := range stmts {
		stmt.String()
		statementAction, err := parseStatement(stmt)
		if err != nil {
			slog.Error("Error parsing node", "error", err)
		}
		result = append(result, statementAction)
	}

	return result, nil

}

func parseInsertStmt(insertStmt *pgquery.InsertStmt) (StatementAction, error) {
	relName := insertStmt.Relation.Relname
	target := ""
	action := SqlActionInsertRow

	return StatementAction{
		Target:   target,
		Action:   action,
		Relation: relName,
	}, nil
}

func parseCreateDatabaseStmt(createDbStmt *pgquery.CreatedbStmt) (StatementAction, error) {
	dbName := createDbStmt.Dbname
	action := SqlActionCreateDatabase

	return StatementAction{
		Target:   dbName,
		Action:   action,
		Relation: "",
	}, nil
}

func parseRenameStmt(renameStmt *pgquery.RenameStmt) (StatementAction, error) {
	relName := ""
	target := renameStmt.Subname
	action := SqlActionUndefined

	switch renameStmt.GetRenameType() {
	// If renaming a table
	case pgquery.ObjectType_OBJECT_TABLE:

		action = SqlActionRenameTable
		target = renameStmt.Relation.Relname

	// If renaming a database
	case pgquery.ObjectType_OBJECT_DATABASE:
		action = SqlActionRenameDatabase
	// If renaming a column
	case pgquery.ObjectType_OBJECT_COLUMN:
		action = SqlActionRenameColumn
		relName = renameStmt.Relation.Relname

	}

	return StatementAction{
		Target:   target,
		Action:   action,
		Relation: relName,
	}, nil
}

func parseDropDatabase(node *pgquery.Node) (StatementAction, error) {
	dropDb := node.GetDropdbStmt()
	dbName := dropDb.Dbname

	return StatementAction{
		Target:   dbName,
		Action:   SqlActionDropDatabase,
		Relation: "",
	}, nil
}

func parseDrop(dropStmt *pgquery.DropStmt) (StatementAction, error) {
	relName := ""
	target := ""
	action := SqlActionDropTable

	// If we are dropping a table
	if dropStmt.RemoveType == pgquery.ObjectType_OBJECT_TABLE {
		tableName := dropStmt.GetObjects()[0].GetList().Items[0].GetString_().GetSval()
		relName = tableName
		target = tableName
		action = SqlActionDropTable
	}

	return StatementAction{
		Target:   target,
		Action:   action,
		Relation: relName,
	}, nil
}

func parseAlterTable(stmt *pgquery.AlterTableStmt) (StatementAction, error) {
	relName := stmt.Relation.Relname
	target := ""
	action := SqlActionUndefined

	for _, cmd := range stmt.Cmds {
		target = cmd.GetAlterTableCmd().Name

		// Determine the action based on the subtype of the command
		switch cmd.GetAlterTableCmd().GetSubtype() {

		// If dropping a column
		case pgquery.AlterTableType_AT_DropColumn:
			action = SqlActionDropColumn
		// If modifying a column data type
		case pgquery.AlterTableType_AT_AlterColumnType:

			action = SqlActionModifyDataType
		// If adding a column
		case pgquery.AlterTableType_AT_AddColumn:
			action = SqlActionAddColumn
			target = cmd.GetAlterTableCmd().Def.GetColumnDef().Colname
		}
	}

	return StatementAction{
		Target:   target,
		Action:   action,
		Relation: relName,
	}, nil
}
