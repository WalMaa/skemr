import { useState, useEffect } from "react";
import { Controller, useForm } from "react-hook-form";
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
import { PencilIcon } from "@phosphor-icons/react";
import { useUpdateDatabase } from "@/api/database";
import { toast } from "sonner";
import type { Database } from "@/types/types";

const SSL_MODE_OPTIONS = [
  "disable",
  "allow",
  "prefer",
  "require",
  "verify-ca",
  "verify-full",
] as const;

const databaseUpdateSchema = z.object({
  displayName: z.string().min(1, "Display name is required"),
  dbName: z.string().min(1, "Database name is required"),
  username: z.string().min(1, "Username is required"),
  password: z.string().optional(),
  host: z.string().min(1, "Host is required"),
  port: z.number().int().min(1).max(65535, "Port must be between 1 and 65535"),
  databaseType: z.literal("postgres"),
  sslMode: z.enum(SSL_MODE_OPTIONS),
});

type DatabaseUpdateFormData = z.infer<typeof databaseUpdateSchema>;

interface DatabaseUpdateDialogProps {
  projectId: string;
  database: Database;
  trigger?: React.ReactElement;
}

export function DatabaseUpdateDialog({
  projectId,
  database,
  trigger,
}: DatabaseUpdateDialogProps) {
  const [isOpen, setIsOpen] = useState(false);

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
    setValue,
  } = useForm<DatabaseUpdateFormData>({
    resolver: zodResolver(databaseUpdateSchema),
    defaultValues: {
      displayName: database.displayName,
      dbName: database.dbName,
      username: database.username,
      password: "",
      host: database.host,
      port: database.port,
      databaseType: database.databaseType,
      sslMode: database.sslMode ?? "prefer",
    },
  });

  const updateDatabaseMutation = useUpdateDatabase();

  // Reset form when database prop changes or dialog opens
  useEffect(() => {
    if (isOpen) {
      reset({
        displayName: database.displayName,
        dbName: database.dbName,
        username: database.username,
        password: "",
        host: database.host,
        port: database.port,
        databaseType: database.databaseType,
        sslMode: database.sslMode ?? "prefer",
      });
    }
  }, [isOpen, database, reset]);

  const onSubmit = (data: DatabaseUpdateFormData) => {
    // Only include password if it was provided
    const updateData = {
      ...data,
      password: data.password || undefined,
    };

    toast.promise(
      updateDatabaseMutation.mutateAsync({
        projectId,
        databaseId: database.id,
        updateData,
      }),
      {
        loading: "Updating database...",
        success: () => {
          setIsOpen(false);
          return "Database updated successfully!";
        },
        error: "Failed to update database. Please try again.",
      },
    );
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger
        render={
          trigger || (
            <Button variant="outline">
              Edit Database
              <PencilIcon className="ml-2" />
            </Button>
          )
        }
      />
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Update Database</DialogTitle>
          <DialogDescription>
            Update the database connection details.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Field>
            <Label htmlFor="displayName">Display Name</Label>
            <Input id="displayName" {...register("displayName")} />
            <FieldDescription>
              A friendly name to identify this database
            </FieldDescription>
            {errors.displayName && (
              <FieldError>{errors.displayName.message}</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="databaseType">Database Type</Label>
            <Select
              value="postgres"
              onValueChange={(value) =>
                setValue("databaseType", value as "postgres")
              }
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="postgres">PostgreSQL</SelectItem>
              </SelectContent>
            </Select>
            {errors.databaseType && (
              <FieldError>{errors.databaseType.message}</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="sslMode">SSL Mode</Label>
            <Controller
              name="sslMode"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value}
                  onValueChange={(value) =>
                    value && field.onChange(value as DatabaseUpdateFormData["sslMode"])
                  }
                >
                  <SelectTrigger id="sslMode">
                    <SelectValue placeholder="Select SSL mode" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="disable">disable</SelectItem>
                    <SelectItem value="allow">allow</SelectItem>
                    <SelectItem value="prefer">prefer</SelectItem>
                    <SelectItem value="require">require</SelectItem>
                    <SelectItem value="verify-ca">verify-ca</SelectItem>
                    <SelectItem value="verify-full">verify-full</SelectItem>
                  </SelectContent>
                </Select>
              )}
            />
            <FieldDescription>
              Controls TLS behavior for PostgreSQL connections.
            </FieldDescription>
            {errors.sslMode && <FieldError>{errors.sslMode.message}</FieldError>}
          </Field>

          <div className="grid grid-cols-2 gap-4">
            <Field>
              <Label htmlFor="host">Host</Label>
              <Input id="host" {...register("host")} />
              {errors.host && <FieldError>{errors.host.message}</FieldError>}
            </Field>

            <Field>
              <Label htmlFor="port">Port</Label>
              <Input
                id="port"
                type="number"
                {...register("port", { valueAsNumber: true })}
              />
              {errors.port && <FieldError>{errors.port.message}</FieldError>}
            </Field>
          </div>

          <Field>
            <Label htmlFor="dbName">Database Name</Label>
            <Input id="dbName" {...register("dbName")} />
            <FieldDescription>The name of the database on the server</FieldDescription>
            {errors.dbName && <FieldError>{errors.dbName.message}</FieldError>}
          </Field>

          <Field>
            <Label htmlFor="username">Username</Label>
            <Input id="username" {...register("username")} />
            {errors.username && (
              <FieldError>{errors.username.message}</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              {...register("password")}
              placeholder="Leave empty to keep current password"
            />
            <FieldDescription>
              Leave empty to keep the current password
            </FieldDescription>
            {errors.password && (
              <FieldError>{errors.password.message}</FieldError>
            )}
          </Field>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={updateDatabaseMutation.isPending}>
              {updateDatabaseMutation.isPending ? "Updating..." : "Update Database"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
