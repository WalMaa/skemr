import { useMemo, useEffect } from "react";
import {
  ReactFlow,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  type Edge,
  type Node,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { useGetDatabaseEntities } from "@/api/database-entity";
import { useGetDatabase } from "@/api/database";
import TableNode from "./table-node";
import DatabaseNode from "./database-node";
import type { DatabaseEntity, Rule } from "@/types/types";
import { useGetRules } from "@/api/rule";
import SchemaGroupNode from "./schema-group-node";
import { useTheme } from "@/components/theme-provider";

const nodeTypes = {
  databaseTable: TableNode,
  database: DatabaseNode,
  schemaGroup: SchemaGroupNode,
};

const GROUP_PADDING = 40;
const GROUP_OFFSET_Y = 250;
const GROUP_OFFSET_X = 40;
const TABLE_WIDTH = 280;
const TABLE_BASE_HEIGHT = 80; // Header height
const TABLE_ROW_HEIGHT = 32; // Height per column row

const SCHEMA_MINIMUM_WIDTH = 250;

const LAYOUT_GROUP_COLUMNS = 4; // Number of tables per row in a group

// Calculate actual table height based on number of columns
function getTableHeight(columnCount: number): number {
  return TABLE_BASE_HEIGHT + columnCount * TABLE_ROW_HEIGHT;
}

export type DataBaseEntityWithRules = DatabaseEntity & { rules: Rule[] };

export function DatabaseSchemaFlow({
  projectId,
  databaseId,
}: {
  projectId: string;
  databaseId: string;
}) {
  const { theme } = useTheme();
  const { data: databaseEntities, isLoading } = useGetDatabaseEntities(
    projectId,
    databaseId,
    "",
  );

  const { data: database } = useGetDatabase(projectId, databaseId);
  const { data: rules } = useGetRules(projectId, databaseId, "");

  // Determine if we're in dark mode
  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches);

  const { nodes: computedNodes, edges: computedEdges } = useMemo(() => {
    if (!databaseEntities || !database) {
      return { nodes: [], edges: [] };
    }
    const nodes: Node[] = [];
    const edges: Edge[] = [];

    // Group by type
    const schemas = databaseEntities.filter((e) => e.type === "schema");
    const tables = databaseEntities.filter((e) => e.type === "table");
    const columns = databaseEntities.filter((e) => e.type === "column");

    let yOffset = 50;

    // Group columns by table
    const columnsByTable = columns.reduce(
      (acc, column) => {
        if (column.parentId) {
          if (!acc[column.parentId]) {
            acc[column.parentId] = []; // Initialize array if it doesn't exist
          }
          // Attach rules to columns
          acc[column.parentId].push({
            ...column,
            rules:
              rules?.filter((rule) => rule.databaseEntity.id === column.id) ||
              [],
          });
        }
        return acc;
      },
      {} as Record<string, DataBaseEntityWithRules[]>,
    );

    // Group tables by schema to calculate group sizes
    const tablesBySchema = tables.reduce(
      (acc, table) => {
        const schemaId = table.parentId || "";
        if (!acc[schemaId]) acc[schemaId] = [];
        acc[schemaId].push(table);
        return acc;
      },
      {} as Record<string, typeof tables>,
    );

    // Calculate total width of all schemas to center database node
    let totalWidth = 0;
    let tempX = 0;
    schemas.forEach((schema) => {
      const schemaTables = tablesBySchema[schema.id] || [];
      const tableCount = schemaTables.length;
      const dynamicWidth = Math.max(
        SCHEMA_MINIMUM_WIDTH,
        Math.min(2, tableCount) * (TABLE_WIDTH + 20) + GROUP_PADDING * 2,
      );
      tempX += GROUP_OFFSET_X;
      totalWidth = tempX + dynamicWidth;
      tempX += dynamicWidth;
    });

    // Add database node centered
    const databaseX = totalWidth > 0 ? totalWidth / 2 - 90 : 300;
    nodes.push({
      id: database.id,
      type: "database",
      position: { x: databaseX, y: yOffset },
      data: { label: database.dbName, type: "database" },
    });

    yOffset += GROUP_OFFSET_Y;

    // Add schema nodes as groups with dynamic height
    let maxGroupHeight = 0;
    let previousGroupX = 0;
    let previousWidth = 0;
    schemas.forEach((schema) => {
      const schemaTables = tablesBySchema[schema.id] || [];
      const tableCount = schemaTables.length;

      // Calculate actual height needed for all tables
      let totalHeight = GROUP_PADDING;
      let currentRowMaxHeight = 0;
      schemaTables.forEach((table, index) => {
        const columns = columnsByTable[table.id] || [];
        const tableHeight = getTableHeight(columns.length);
        const isNewRow = index % LAYOUT_GROUP_COLUMNS === 0 && index > 0;

        if (isNewRow) {
          totalHeight += currentRowMaxHeight + 20;
          currentRowMaxHeight = 0;
        }
        currentRowMaxHeight = Math.max(currentRowMaxHeight, tableHeight);
      });
      totalHeight += currentRowMaxHeight + GROUP_PADDING;

      const dynamicHeight = Math.max(250, totalHeight); // Minimum height 250
      const dynamicWidth = Math.max(
        SCHEMA_MINIMUM_WIDTH,
        Math.min(LAYOUT_GROUP_COLUMNS, tableCount) * (TABLE_WIDTH + 20) +
          GROUP_PADDING * 2,
      );

      maxGroupHeight = Math.max(maxGroupHeight, dynamicHeight);
      const groupX = previousGroupX + previousWidth + GROUP_OFFSET_X;
      previousGroupX = groupX;
      previousWidth = dynamicWidth;

      nodes.push({
        id: schema.id,
        type: "schemaGroup",
        position: { x: groupX, y: yOffset },
        data: {
          schema: {
            ...schema,
            rules:
              rules?.filter((rule) => rule.databaseEntity.id === schema.id) ||
              [],
          },
        },
        style: {
          width: dynamicWidth,
          height: dynamicHeight,
        },
      });
      // Connect schema to database
      edges.push({
        id: `e-${database.id}-${schema.id}`,
        source: database.id,
        selectable: false,
        target: schema.id,
        type: "smoothstep",
      });
    });

    yOffset += maxGroupHeight + 40; // Dynamic offset based on tallest group plus padding

    // Add table nodes inside groups - position relative to each schema
    Object.entries(tablesBySchema).forEach(([schemaId, schemaTables]) => {
      let currentY = GROUP_PADDING;
      let currentRow = 0;
      let maxHeightInRow = 0;

      schemaTables.forEach((table, index) => {
        const columns = columnsByTable[table.id] || [];
        const tableHeight = getTableHeight(columns.length);

        // Check if we're starting a new row
        const col = index % LAYOUT_GROUP_COLUMNS;
        const row = Math.floor(index / LAYOUT_GROUP_COLUMNS);

        if (row !== currentRow) {
          // Move to next row
          currentY += maxHeightInRow + 20; // Add spacing between rows
          maxHeightInRow = 0;
          currentRow = row;
        }

        maxHeightInRow = Math.max(maxHeightInRow, tableHeight);

        const x = col * (TABLE_WIDTH + 20) + GROUP_PADDING; // Add horizontal spacing

        nodes.push({
          id: table.id,
          type: "databaseTable",
          position: { x, y: currentY },
          parentId: schemaId,
          data: {
            table: {
              ...table,
              rules:
                rules?.filter((rule) => rule.databaseEntity.id === table.id) ||
                [],
            } as DataBaseEntityWithRules,
            label: table.name,
            type: table.type,
            columns: columns,
          },
          extent: "parent",
        });
      });
    });

    // Remove column nodes

    return { nodes, edges };
  }, [databaseEntities, database, rules]);

  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  useEffect(() => {
    setNodes(computedNodes);
    setEdges(computedEdges);
  }, [computedNodes, computedEdges, setNodes, setEdges]);

  return (
    <div style={{ width: "100%", height: "600px" }}>
      {isLoading ? (
        <div className="flex items-center justify-center h-full">
          <div className="text-lg text-muted-foreground">
            Loading database schema...
          </div>
        </div>
      ) : (
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          nodeTypes={nodeTypes}
          elementsSelectable={true}
          nodesDraggable={false}
          connectOnClick={false}
          edgesReconnectable={false}
          fitView
        >
          <Controls
            style={{
              backgroundColor: isDark ? "hsl(var(--background))" : undefined,
              borderColor: isDark ? "hsl(var(--border))" : undefined,
            }}
          />
          <Background gap={12} size={1} />
        </ReactFlow>
      )}
    </div>
  );
}
