import { useState } from "react";
import { useForm } from "react-hook-form";
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
import { PlusIcon } from "@phosphor-icons/react";
import { useCreateDatabase } from "@/api/database";
import { toast } from "sonner";
import { logger } from "@/lib/logger";

const SSL_MODE_OPTIONS = [
  "disable",
  "allow",
  "prefer",
  "require",
  "verify-ca",
  "verify-full",
] as const;

const databaseSchema = z.object({
  displayName: z.string().min(1, "Display name is required"),
  dbName: z.string().min(1, "Database name is required"),
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
  host: z.string().min(1, "Host is required"),
  port: z.number().int().min(1).max(65535, "Port must be between 1 and 65535"),
  databaseType: z.literal("postgres"),
  sslMode: z.enum(SSL_MODE_OPTIONS),
});

type DatabaseFormData = z.infer<typeof databaseSchema>;

const isSslMode = (value: string): value is DatabaseFormData["sslMode"] =>
  SSL_MODE_OPTIONS.includes(value as DatabaseFormData["sslMode"]);

interface DatabaseCreationDialogProps {
  projectId: string;
  trigger?: React.ReactElement;
}

export function DatabaseCreationDialog({
  projectId,
  trigger,
}: DatabaseCreationDialogProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [connectionString, setConnectionString] = useState("");

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
    getValues,
  } = useForm<DatabaseFormData>({
    resolver: zodResolver(databaseSchema),
    defaultValues: {
      displayName: "",
      dbName: "",
      username: "",
      password: "",
      host: "",
      port: 5432,
      databaseType: "postgres",
      sslMode: "prefer",
    },
  });

  const createDatabaseMutation = useCreateDatabase();

  const parseConnectionString = (connStr: string) => {
    try {
      const parsedUrl = new URL(connStr.trim());

      if (
        parsedUrl.protocol !== "postgres:" &&
        parsedUrl.protocol !== "postgresql:"
      ) {
        toast.error(
          "Invalid protocol. Expected postgres:// or postgresql://",
        );
        return;
      }

      const dbName = decodeURIComponent(parsedUrl.pathname.replace(/^\/+/, ""));
      if (!parsedUrl.hostname || !dbName) {
        toast.error(
          "Invalid connection string format. Expected: postgresql://[username:password@]host:port/database[?sslmode=prefer]",
        );
        return;
      }

      const parsedPort = parsedUrl.port
        ? Number.parseInt(parsedUrl.port, 10)
        : 5432;
      if (Number.isNaN(parsedPort) || parsedPort < 1 || parsedPort > 65535) {
        toast.error("Port must be between 1 and 65535");
        return;
      }

      const sslModeParam = (
        parsedUrl.searchParams.get("sslmode") ??
        parsedUrl.searchParams.get("sslMode")
      )
        ?.trim()
        .toLowerCase();

      let sslMode: DatabaseFormData["sslMode"] = "prefer";
      if (sslModeParam) {
        if (!isSslMode(sslModeParam)) {
          toast.error(
            "Unsupported sslmode. Use one of: disable, allow, prefer, require, verify-ca, verify-full",
          );
          return;
        }

        sslMode = sslModeParam;
      }

      setValue("host", parsedUrl.hostname);
      setValue("port", parsedPort);
      setValue("dbName", dbName);
      setValue("displayName", `${dbName} (${parsedUrl.hostname})`);
      setValue("sslMode", sslMode);

      const username = decodeURIComponent(parsedUrl.username);
      const password = decodeURIComponent(parsedUrl.password);

      if (username || password) {
        setValue("username", username);
        setValue("password", password);
        toast.success("Connection string parsed successfully!");
      } else {
        setValue("username", "");
        setValue("password", "");
        toast.success(
          "Connection string parsed successfully! Please enter username and password.",
        );
      }
    } catch (error) {
      logger.error(error, "Error parsing connection string:");
      toast.error("Failed to parse connection string");
    }
  };

  const handleConnectionStringChange = (value: string) => {
    setConnectionString(value);
    if (value.trim()) {
      parseConnectionString(value);
    }
  };

  const onSubmit = (data: DatabaseFormData) => {
    toast.promise(
      createDatabaseMutation.mutateAsync({
        projectId,
        databaseData: data,
      }),
      {
        loading: "Creating database...",
        success: () => {
          setIsOpen(false);
          reset();
          setConnectionString("");
          return "Database created successfully!";
        },
        error: "Failed to create database. Please try again.",
      },
    );
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger
        render={
          trigger || (
            <Button>
              <PlusIcon className="mr-2 h-4 w-4" />
              Add Database
            </Button>
          )
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add New Database</DialogTitle>
          <DialogDescription>
            Connect a new database to this project for schema management and
            rule enforcement.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <Field>
            <Label htmlFor="connection-string">
              Connection String (Optional)
            </Label>
            <Input
              id="connection-string"
              value={connectionString}
              onChange={(e) => handleConnectionStringChange(e.target.value)}
              placeholder="postgresql://username:password@host:port/database?sslmode=prefer"
            />
            <FieldDescription>
              Paste a connection string to auto-fill the form fields below.
              Format: postgresql://[username:password@]host:port/database?sslmode=prefer
            </FieldDescription>
          </Field>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">
                Or fill manually
              </span>
            </div>
          </div>

          <Field>
            <Label htmlFor="display-name">
              Display Name <span className="text-destructive">*</span>
            </Label>
            <Input
              id="display-name"
              {...register("displayName")}
              placeholder="My Production Database"
            />
            {errors.displayName && (
              <FieldError>{errors.displayName.message}</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="db-name">
              Database Name <span className="text-destructive">*</span>
            </Label>
            <Input
              id="db-name"
              {...register("dbName")}
              placeholder="myapp_production"
            />
            {errors.dbName && <FieldError>{errors.dbName.message}</FieldError>}
          </Field>

          <Field>
            <Label htmlFor="database-type">Database Type</Label>
            <Select
              disabled
              value={getValues("databaseType")}
              onValueChange={(value) =>
                value && setValue("databaseType", value as "postgres")
              }
            >
              <SelectTrigger id="database-type">
                <SelectValue placeholder="Select database type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="postgres">PostgreSQL</SelectItem>
              </SelectContent>
            </Select>
            <FieldDescription>
              Currently, only PostgreSQL databases are supported.
            </FieldDescription>
            {errors.databaseType && (
              <FieldError>{errors.databaseType.message}</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="ssl-mode">SSL Mode</Label>
            <Select
              value={getValues("sslMode")}
              onValueChange={(value) =>
                value && setValue("sslMode", value as DatabaseFormData["sslMode"])
              }
            >
              <SelectTrigger id="ssl-mode">
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
            <FieldDescription>
              Controls TLS behavior for PostgreSQL connections.
            </FieldDescription>
            {errors.sslMode && <FieldError>{errors.sslMode.message}</FieldError>}
          </Field>

          <Field>
            <Label htmlFor="host">
              Host <span className="text-destructive">*</span>
            </Label>
            <Input
              id="host"
              {...register("host")}
              placeholder="localhost or db.example.com"
            />
            {errors.host && <FieldError>{errors.host.message}</FieldError>}
          </Field>

          <Field>
            <Label htmlFor="port">
              Port <span className="text-destructive">*</span>
            </Label>
            <Input
              id="port"
              type="number"
              {...register("port", { valueAsNumber: true })}
              placeholder="5432"
            />
            {errors.port && <FieldError>{errors.port.message}</FieldError>}
          </Field>

          <Field>
            <Label htmlFor="username">
              Username <span className="text-destructive">*</span>
            </Label>
            <Input
              id="username"
              {...register("username")}
              placeholder="postgres"
            />
            {errors.username && (
              <FieldError>{errors.username.message}</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="password">
              Password <span className="text-destructive">*</span>
            </Label>
            <Input
              id="password"
              type="password"
              {...register("password")}
              placeholder="Enter database password"
            />
            {errors.password && (
              <FieldError>{errors.password.message}</FieldError>
            )}
          </Field>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setIsOpen(false)}>
            Cancel
          </Button>
          <Button onClick={handleSubmit(onSubmit)}>Create Database</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
