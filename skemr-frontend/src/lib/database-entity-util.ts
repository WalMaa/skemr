import type { DatabaseEntity } from "@/types/types";

/**
 * Transform database entities to a hierarchical structure
 * schema
 * --tables
 * ----columns
 *
 * @param entities
 */
export const entitiesToTree = (
    entities: DatabaseEntity[],
): DatabaseEntityWithItems[] => {
    const schemas = entities.filter((e) => e.type === "schema");
    const tables = entities.filter((e) => e.type === "table");
    const columns = entities.filter((e) => e.type === "column");
    const constraints = entities.filter((e) => e.type === "constraint");
    const indexes = entities.filter((e) => e.type === "index");

    // Assign columns, indexes and constraints to their respective tables
    const tablesWithItems = tables.map((table) => ({
        ...table,
        items: () => {
            const tableColumns = columns.filter(
                (col) => col.parentId === table.id,
            );
            const tableIndexes = indexes.filter(
                (idx) => idx.parentId === table.id,
            );
            const tableConstraints = constraints.filter(
                (con) => con.parentId === table.id,
            );
            return [...tableColumns, ...tableIndexes, ...tableConstraints];
        },
    }));

    // Assign tables to their respective schemas
    const schemasWithTables = schemas.map((schema) => ({
        ...schema,
        items: tablesWithItems.filter((table) => table.parentId === schema.id),
    }));

    return schemasWithTables;
};

type DatabaseEntityWithItems = DatabaseEntity & {
    items: DatabaseEntity[] | null;
};
