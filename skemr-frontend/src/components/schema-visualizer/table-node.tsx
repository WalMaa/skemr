import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

import DataBaseEntitySheet from "./entity-sheet";
import type { DataBaseEntityWithRules } from "./database-schema-visualizer";
import { cn } from "@/lib/utils";

export default function TableNode({
  data,
}: {
  data: {
    label: string;
    table: DataBaseEntityWithRules;
    type: string;
    columns?: DataBaseEntityWithRules[];
  };
}) {
  const { columns, table } = data;
  return (
    <Card
      className={cn(
        "w-64 shadow-md gap-2 cursor-default",
        data.table.status === "deleted" &&
          "border border-dashed border-destructive/75",
      )}
    >
      <CardHeader className="pb-2">
        <div className="flex items-center gap-1 justify-between">
          <div className="flex relative min-w-0 gap-2 w-full items-center">
            <span className="absolute -top-2.5 text-[9px] text-muted-foreground">
              table
            </span>
            <CardTitle className="min-w-0">
              <h3 className="text-sm w-36 truncate">{data.label}</h3>
            </CardTitle>
            <Badge size={"sm"} variant="outline">
              {table.rules.length} Rule
              {table.rules.length !== 1 ? "s" : ""}
            </Badge>
          </div>
          <div className="mr-1 shrink-0">
            <DataBaseEntitySheet type="table" entity={table} />
          </div>
        </div>
      </CardHeader>
      {columns && columns.length > 0 && (
        <CardContent className="pt-0">
          <div>
            {columns.map((column, index) => {
              const columnRules = column.rules || [];
              return (
                <div
                  className={cn(
                    "flex border-t [[data-column-status=deleted]+&]:border-t-0 hover:bg-accent transition-colors px-1 py-1 justify-between items-center",
                    column.status === "deleted" &&
                      "text-foreground  border-dashed border border-destructive/75 hover:bg-destructive/5",
                  )}
                  data-column-name={column.name}
                  data-column-status={column.status}
                  key={index}
                >
                  <div
                    key={index}
                    className="text-xs flex items-center text-muted-foreground"
                  >
                    {column.name}
                    {columnRules.length > 0 && (
                      <Badge size={"sm"} variant="outline" className="ml-2">
                        {columnRules.length} Rule
                        {columnRules.length > 1 ? "s" : ""}
                      </Badge>
                    )}
                  </div>

                  <DataBaseEntitySheet type="column" entity={column} />
                </div>
              );
            })}
          </div>
        </CardContent>
      )}
    </Card>
  );
}
