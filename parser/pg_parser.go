package parser

import (
	"fmt"
	"log/slog"

	pgquery "github.com/pganalyze/pg_query_go/v6"
)

type StatementAction struct {
	Target   string
	Action   SqlAction
	Relation string
}

type SqlAction string

const (
	SqlActionDropColumn     SqlAction = "DROP COLUMN"
	SqlActionRenameColumn   SqlAction = "RENAME COLUMN"
	SqlActionModifyDataType SqlAction = "MODIFY DATA TYPE"
	SqlActionDropDatabase   SqlAction = "DROP DATABASE"
	SqlActionUndefined      SqlAction = "UNDEFINED"
	SqlActionInsertRow      SqlAction = "INSERT ROW"
	SqlActionCreateTable    SqlAction = "CREATE TABLE"
	SqlActionRenameTable    SqlAction = "RENAME TABLE"
	SqlActionDropTable      SqlAction = "DROP TABLE"
)

/*
ParseSql parses the SQL statement and returns a structured representation of the SQL.
*/
func ParseSql(sql string) (StatementAction, error) {
	tree, err := pgquery.Parse(sql)
	if err != nil {
		slog.Error("Failed to parse SQL", "error", err, "sql", sql)
		return StatementAction{
			Action: SqlActionUndefined,
		}, err
	}
	stmts := tree.Stmts

	if len(stmts) == 0 {
		return StatementAction{}, nil
	}

	node := stmts[0].GetStmt()

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

	// If the statement is not recognized, return an error
	return StatementAction{
		Target:   "",
		Action:   SqlActionUndefined,
		Relation: "",
	}, fmt.Errorf("unsupported SQL statement: %s", sql)
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

func parseRenameStmt(renameStmt *pgquery.RenameStmt) (StatementAction, error) {
	relName := renameStmt.Relation.Relname
	target := renameStmt.Subname
	action := SqlActionRenameColumn

	typer := renameStmt.GetRenameType()
	fmt.Println(typer)

	switch renameStmt.GetRenameType() {
	// If renaming a table
	case pgquery.ObjectType_OBJECT_TABLE:
		// If renaming a column
		action = SqlActionRenameTable
	case pgquery.ObjectType_OBJECT_COLUMN:

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
		relName = dropStmt.GetObjects()[0].GetList().Items[0].GetString_().GetSval()
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

		case pgquery.AlterTableType_AT_DropColumn:
			action = SqlActionDropColumn
		case pgquery.AlterTableType_AT_AlterColumnType:

			action = SqlActionModifyDataType
		}
	}

	return StatementAction{
		Target:   target,
		Action:   action,
		Relation: relName,
	}, nil
}
