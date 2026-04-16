import { PostgresIcon } from "@/assets/icons";
import { type NodeProps, Handle, Position } from "@xyflow/react";
import { Card } from "../ui/card";

export default function DatabaseNode({
  data,
}: NodeProps & { data: { label: string } }) {
  return (
    <Card className="flex flex-col gap-2 items-center justify-center min-w-[180px]">
      <Handle type="source" position={Position.Bottom} />
      <PostgresIcon />
      <div className="text-center font-semibold  text-lg">{data.label}</div>
      <div className="text-xs text-muted-foreground">PostgreSQL</div>
    </Card>
  );
}
