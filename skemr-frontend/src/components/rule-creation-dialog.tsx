import { useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Field, FieldDescription, FieldError } from "@/components/ui/field";
import {
  ArchiveIcon,
  LockIcon,
  PlusIcon,
  WarningIcon,
} from "@phosphor-icons/react";
import { useCreateRule } from "@/api/rule";
import { useGetDatabaseEntities } from "@/api/database-entity";
import { toast } from "sonner";
import { EntityTreeSelector } from "@/components/entity-tree-selector";
import { InfoIcon } from "@phosphor-icons/react/dist/ssr";

const ruleSchema = z.object({
  name: z.string().min(1, "Rule name is required"),
  ruleType: z.enum(["locked", "warn", "advisory", "deprecated"], {
    error: "Rule type is required",
  }),
  databaseEntityId: z.string().min(1, "Database entity is required"),
});

type RuleFormData = z.infer<typeof ruleSchema>;

interface RuleCreationDialogProps {
  projectId: string;
  databaseId: string;
  trigger?: React.ReactElement;
}

export function RuleCreationDialog({
  projectId,
  databaseId,
  trigger,
}: RuleCreationDialogProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [entityPopoverOpen, setEntityPopoverOpen] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    control,
  } = useForm<RuleFormData>({
    resolver: zodResolver(ruleSchema),
    defaultValues: {
      name: "",
      ruleType: undefined,
      databaseEntityId: "",
    },
  });

  const { data: entities, isLoading: entitiesLoading } = useGetDatabaseEntities(
    projectId,
    databaseId,
    "",
  );
  const createRuleMutation = useCreateRule();

  const onSubmit = (data: RuleFormData) => {
    toast.promise(
      createRuleMutation.mutateAsync({
        projectId,
        databaseId,
        ruleData: data,
      }),
      {
        loading: "Creating rule...",
        success: () => {
          setIsOpen(false);
          reset();
          return "Rule created successfully!";
        },
        error: (err) => `Failed to create rule: ${err.message}`,
      },
    );
  };

  const ruleSelectItems = {
    locked: {
      label: "Locked - Prevents modifications",
      icon: LockIcon,
    },
    warn: {
      label: "Warn - Issues warnings",
      icon: WarningIcon,
    },
    advisory: {
      label: "Advisory - Provides guidance",
      icon: InfoIcon,
    },
    deprecated: {
      label: "Deprecated - Marks as deprecated",
      icon: ArchiveIcon,
    },
  };

  const getRuleTypeIcon = (ruleType?: string) => {
    const item = ruleSelectItems[ruleType as keyof typeof ruleSelectItems];
    return item ? <item.icon className="text-muted-foreground" /> : null;
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger
        render={
          trigger || (
            <Button>
              <PlusIcon className="mr-2 h-4 w-4" />
              Create Rule
            </Button>
          )
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New Rule</DialogTitle>
          <DialogDescription>
            Add a new rule to enforce constraints on database entities.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <Field>
            <Label htmlFor="name">
              Rule Name <span className="text-destructive">*</span>
            </Label>
            <Input
              id="name"
              {...register("name")}
              placeholder="e.g., No modifications to users table"
            />
            {errors.name && <FieldError>{errors.name.message}</FieldError>}
          </Field>

          <Field>
            <Label htmlFor="rule-type">
              Rule Type <span className="text-destructive">*</span>
            </Label>
            <Controller
              name="ruleType"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value}
                  onValueChange={(value) => field.onChange(value)}
                >
                  <SelectTrigger id="rule-type">
                    <SelectValue placeholder="Select rule type">
                      {getRuleTypeIcon(field.value)}
                      {ruleSelectItems[
                        field.value as keyof typeof ruleSelectItems
                      ]?.label || "Select rule type"}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    {Object.entries(ruleSelectItems).map(
                      ([value, { label, icon: Icon }]) => (
                        <SelectItem key={value} value={value}>
                          <Icon className="text-muted-foreground" />
                          {label}
                        </SelectItem>
                      ),
                    )}
                  </SelectContent>
                </Select>
              )}
            />
            <FieldDescription>
              Choose the enforcement level for this rule.
            </FieldDescription>
            {errors.ruleType && (
              <FieldError>{errors.ruleType.message}</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="entity">
              Database Entity <span className="text-destructive">*</span>
            </Label>
            <Controller
              name="databaseEntityId"
              control={control}
              render={({ field }) => (
                <EntityTreeSelector
                  id="entity"
                  entities={entities ?? []}
                  value={field.value}
                  onChange={field.onChange}
                  popoverOpen={entityPopoverOpen}
                  onPopoverOpenChange={setEntityPopoverOpen}
                  loading={entitiesLoading}
                />
              )}
            />
            <FieldDescription>
              Select the database entity to apply this rule to.
            </FieldDescription>
            {errors.databaseEntityId && (
              <FieldError>{errors.databaseEntityId.message}</FieldError>
            )}
          </Field>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setIsOpen(false)}>
            Cancel
          </Button>
          <Button
            isLoading={createRuleMutation.isPending}
            disabled={createRuleMutation.isPending}
            onClick={handleSubmit(onSubmit)}
          >
            Create Rule
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
