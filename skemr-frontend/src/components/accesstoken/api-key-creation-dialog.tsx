import React, { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { CalendarIcon, PlusIcon } from "@phosphor-icons/react";
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
import { format } from "date-fns";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover.tsx";
import { Field, FieldError, FieldLabel } from "@/components/ui/field.tsx";
import { Calendar } from "@/components/ui/calendar.tsx";

interface ApiKeyCreationDialogProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  onCreate: (name: string, expiresAt?: Date) => void;
  trigger?: React.ReactElement;
}

interface ApiKeyCreationFormData {
  name: string;
  expiresAt?: Date;
}

const apiKeyCreationSchema = z.object({
  name: z.string().min(2, "Key name must be at least 2 characters"),
  expiresAt: z.date().optional(),
});

export function ApiKeyCreationDialog({
  open,
  onOpenChange,
  onCreate,
  trigger,
}: ApiKeyCreationDialogProps) {
  const [isDatePickerOpen, setIsDatePickerOpen] = useState(false);
  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<ApiKeyCreationFormData>({
    resolver: zodResolver(apiKeyCreationSchema),
    mode: "onBlur",
    defaultValues: {
      name: "",
      expiresAt: undefined,
    },
  });

  const handleCreate = (data: ApiKeyCreationFormData) => {
    onCreate(data.name, data.expiresAt);

    onOpenChange?.(false);
    reset();
    setIsDatePickerOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger
        render={
          trigger || (
            <Button>
              <PlusIcon className="mr-2 h-4 w-4" />
              Create CI/CD API Key
            </Button>
          )
        }
      ></DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New CI/CD API Key</DialogTitle>
          <DialogDescription>
            Create a new API key for CI/CD containers to access this
            project&apos;s resources.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleCreate)}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <FieldLabel htmlFor="key-name">Key Name</FieldLabel>
              <Input
                id="key-name"
                {...register("name")}
                placeholder="Enter API key name (e.g., GitLab CI, GitHub Actions)"
              />
              {errors.name && <FieldError>{errors.name.message}</FieldError>}
            </div>
            <Field>
              <FieldLabel htmlFor="expiration-date">
                Expiration Date (optional)
              </FieldLabel>
              <Controller
                name="expiresAt"
                control={control}
                render={({ field }) => (
                  <Popover
                    open={isDatePickerOpen}
                    onOpenChange={setIsDatePickerOpen}
                  >
                    <PopoverTrigger
                      render={
                        <Button
                          variant="outline"
                          data-empty={!field.value}
                          className=" justify-between text-left font-normal data-[empty=true]:text-muted-foreground"
                        >
                          {field.value ? (
                            format(field.value, "PPP")
                          ) : (
                            <span>Pick a date</span>
                          )}
                          <CalendarIcon />
                        </Button>
                      }
                    ></PopoverTrigger>
                    <PopoverContent
                      className="w-auto overflow-hidden p-0"
                      align="end"
                      alignOffset={-8}
                      sideOffset={10}
                    >
                      <Calendar
                        disabled={{ before: new Date() }}
                        mode="single"
                        selected={field.value}
                        onSelect={(val) => {
                          field.onChange(val);
                          setIsDatePickerOpen(false);
                        }}
                      />
                    </PopoverContent>
                  </Popover>
                )}
              />
            </Field>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange?.(false)}
            >
              Cancel
            </Button>
            <Button type="submit">Create CI/CD Key</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
