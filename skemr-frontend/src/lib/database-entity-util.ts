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

  // Assign columns to their respective tables
  const tablesWithColumns = tables.map((table) => ({
    ...table,
    items: columns.filter((col) => col.parentId === table.id),
  }));

  // Assign tables to their respective schemas
  const schemasWithTables = schemas.map((schema) => ({
    ...schema,
    items: tablesWithColumns.filter((table) => table.parentId === schema.id),
  }));

  return schemasWithTables;
};

type DatabaseEntityWithItems = DatabaseEntity & {
  items: DatabaseEntity[] | null;
};
