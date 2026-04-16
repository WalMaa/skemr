import { type NodeProps, Handle, Position } from "@xyflow/react";
import { Card } from "../ui/card";
import { Badge } from "../ui/badge";
import { cn } from "@/lib/utils";
import DataBaseEntitySheet from "./entity-sheet";
import type { DataBaseEntityWithRules } from "./database-schema-visualizer";

export default function SchemaGroupNode({
  data,
}: NodeProps & { data: { schema: DataBaseEntityWithRules } }) {
  const { schema } = data;
  return (
    <>
      <Handle
        className="z-10"
        type="target"
        isConnectable={false}
        position={Position.Top}
      />
      <Card
        className={cn(
          "w-full h-full border relative bg-card/50 p-0",
          schema.status === "deleted" && "border-destructive/75 border-dashed ",
        )}
      >
        <div className="flex items-center gap-2 justify-between p-2">
          <div className="flex items-center gap-1">
            <Badge>{schema.name}</Badge>
            <Badge className="text-[9px]"  variant="outline">
              {schema.rules.length} Rule
              {schema.rules.length !== 1 ? "s" : ""}
            </Badge>
          </div>
          <DataBaseEntitySheet type="column" entity={schema} />
        </div>
      </Card>
    </>
  );
}
