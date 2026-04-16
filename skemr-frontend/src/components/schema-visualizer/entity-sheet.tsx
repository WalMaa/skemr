import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "../ui/sheet";
import { Button } from "../ui/button";
import {
  DotsThreeOutlineVerticalIcon,
  FloppyDiskIcon,
  TrashIcon,
} from "@phosphor-icons/react";
import type { RuleCreationDto, DatabaseRuleType } from "@/types/types";
import { useFieldArray, useForm, Controller } from "react-hook-form";
import { Input } from "../ui/input";
import { useCreateRule, useDeleteRule } from "@/api/rule";
import { useParams } from "@tanstack/react-router";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import type { DataBaseEntityWithRules } from "./database-schema-visualizer";
import React from "react";
import { format } from "date-fns";
import { Badge } from "../ui/badge";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";
import type { Dispatch, SetStateAction } from "react";

interface RuleInput {
  ruleId?: string;
  name: string;
  ruleType: DatabaseRuleType;
}

export default function DataBaseEntitySheet({
  entity,
  type = "generic",
  open,
  onOpenChange,
  hideTrigger = false,
}: {
  entity: DataBaseEntityWithRules;
  type?: "column" | "table" | "schema" | "generic";
  open?: boolean;
  onOpenChange?: Dispatch<SetStateAction<boolean>> | ((open: boolean) => void);
  hideTrigger?: boolean;
}) {
  const { projectId, databaseId } = useParams({
    from: "/(project)/projects/$projectId/databases/$databaseId/",
  });
  const { control, register, getValues, reset } = useForm<{
    rules: RuleInput[];
  }>({
    defaultValues: {
      rules: entity.rules.map((rule) => ({
        ruleId: rule.id,
        name: rule.name,
        ruleType: rule.ruleType,
      })),
    },
  });

  // Reload rules when entity changes
  React.useEffect(() => {
    const rules = entity.rules.map((rule) => ({
      ruleId: rule.id,
      name: rule.name,
      ruleType: rule.ruleType,
    }));
    reset({ rules });
  }, [entity, control, reset]);

  const { fields, append, remove } = useFieldArray({
    control,
    name: "rules",
  });
  const createRuleMutation = useCreateRule();
  const deleteRuleMutation = useDeleteRule();

  const handleSaveRule = async (index: number) => {
    const rule = fields[index];
    const formValues = getValues(`rules.${index}`);

    // Don't save if rule already exists (has a ruleId)
    if (rule.ruleId) {
      return;
    }

    if (!formValues.name || !formValues.ruleType) {
      return;
    }

    try {
      const ruleData: RuleCreationDto = {
        name: formValues.name,
        ruleType: formValues.ruleType,
        databaseEntityId: entity.id,
      };
      await createRuleMutation.mutateAsync({
        projectId,
        databaseId,
        ruleData,
      });
      remove(index);
    } catch (error) {
      console.error("Failed to save rule:", error);
    }
  };

  const handleNewRuleCreation = () => {
    append({ name: "", ruleType: "warn" });
  };

  const handleDeleteRule = async (index: number) => {
    const rule = fields[index];

    // If rule has a ruleId, delete from server
    if (rule.ruleId) {
      try {
        await deleteRuleMutation.mutateAsync({
          projectId,
          databaseId,
          ruleId: rule.ruleId,
        });
        remove(index);
      } catch (error) {
        console.error("Failed to delete rule:", error);
      }
    } else {
      // If no ruleId, just remove from form
      remove(index);
    }
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      {!hideTrigger && (
        <SheetTrigger
          render={
            <Button size={"icon-xs"} variant={"outline"}>
              <DotsThreeOutlineVerticalIcon />
            </Button>
          }
        />
      )}
      <SheetContent>
        <SheetHeader>
          <SheetTitle className={"text-xl"}>
            <span className="capitalize">{type}</span>{" "}
            <samp>{entity.name}</samp>{" "}
          </SheetTitle>
          <SheetDescription>
            Create rules or view details for {type} <samp>{entity.name}</samp>.
          </SheetDescription>
        </SheetHeader>

        <div className="px-4 mb-4">
          <h3 className="text-lg mb-2">Info</h3>
          <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1 text-sm">
            <span className="text-muted-foreground font-medium">Status</span>
            <Badge
              variant={entity.status === "deleted" ? "destructive" : "default"}
              className="capitalize"
            >
              {entity.status}
            </Badge>
            <span className="text-muted-foreground font-medium">
              First seen
            </span>
            <span>{format(entity.firstSeenAt, "PPP p")}</span>
            {entity.status === "deleted" && entity.deletedAt && (
              <>
                <span className="text-muted-foreground font-medium">
                  Deleted at
                </span>
                <span>{format(entity.deletedAt, "PPP p")}</span>
              </>
            )}
            {entity.attributes &&
              Object.entries(entity.attributes).map(([key, value]) => (
                <React.Fragment key={key}>
                  <span className="text-muted-foreground font-medium">
                    {key}
                  </span>
                  <span>{value}</span>
                </React.Fragment>
              ))}
          </div>
        </div>

        <div className="px-4">
          <h3 className="text-lg">Rules</h3>
          <div className="grid grid-cols-[1fr_1fr_auto_0.5fr] my-2">
            <h4 className="text-sm font-semibold">Name</h4>
            <h4 className="text-sm font-semibold">Type</h4>
          </div>
          {/* List rules associated with this entity */}
          {fields.length === 0 && (
            <p className="text-sm text-muted-foreground">No rules added yet.</p>
          )}
          {fields.map((field, index) => (
            <div
              className="grid grid-cols-[1fr_1fr_auto_0.5fr] gap-2 mb-2"
              key={field.id}
            >
              <Input
                {...register(`rules.${index}.name`, {
                  required: "Rule name is required",
                })}
                placeholder="Rule name"
                disabled={!!field.ruleId}
              />
              <Controller
                control={control}
                name={`rules.${index}.ruleType`}
                rules={{ required: "Rule type is required" }}
                render={({ field: controllerField }) => (
                  <Select
                    value={controllerField.value}
                    onValueChange={controllerField.onChange}
                    disabled={!!field.ruleId}
                  >
                    <SelectTrigger className="w-full col-auto">
                      <SelectValue placeholder="Select type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="locked">Locked</SelectItem>
                      <SelectItem value="warn">Warn</SelectItem>
                      <SelectItem value="advisory">Advisory</SelectItem>
                      <SelectItem value="deprecated">Deprecated</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
              <div className="col-span-2 flex items-center justify-end gap-2">
                {!field.ruleId && (
                  <Tooltip>
                    <TooltipTrigger>
                      <Button
                        type="button"
                        size="icon"
                        aria-label="Save rule"
                        onClick={() => handleSaveRule(index)}
                        disabled={createRuleMutation.isPending}
                      >
                        <FloppyDiskIcon />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Save rule</TooltipContent>
                  </Tooltip>
                )}
                <Tooltip>
                  <TooltipTrigger>
                    <Button
                      type="button"
                      size="icon"
                      aria-label="Delete rule"
                      variant="destructive"
                      onClick={() => handleDeleteRule(index)}
                      disabled={deleteRuleMutation.isPending}
                    >
                      <TrashIcon />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Delete rule</TooltipContent>
                </Tooltip>
              </div>
            </div>
          ))}
          <div>
            <Button
              type="button"
              onClick={handleNewRuleCreation}
              className="mt-2"
              variant="outline"
            >
              Add Rule
            </Button>
          </div>
        </div>
        <SheetFooter className="flex justify-end flex-row gap-2 ">
          <SheetClose render={<Button variant={"secondary"}>Close</Button>} />
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
