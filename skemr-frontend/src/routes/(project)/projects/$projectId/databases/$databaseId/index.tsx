import {
  databaseDetailQuery,
  useGetDatabase,
  useSyncDatabaseSchema,
  useDeleteDatabase,
} from "@/api/database";
import { useDeleteRule } from "@/api/rule";
import type { RouterContext } from "@/routes/__root";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
  TableCaption,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  ArrowsCounterClockwiseIcon,
  DatabaseIcon,
  TrashIcon,
  PencilIcon,
  DotsThreeOutlineVerticalIcon,
  InfoIcon,
} from "@phosphor-icons/react";
import { useGetRules } from "@/api/rule";
import type { Rule } from "@/types/types";
import { DatabaseSchemaFlow } from "@/components/schema-visualizer/database-schema-visualizer";
import type { DataBaseEntityWithRules } from "@/components/schema-visualizer/database-schema-visualizer";
import DataBaseEntitySheet from "@/components/schema-visualizer/entity-sheet";
import { DatabaseUpdateDialog } from "@/components/database-update-dialog";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { formatRelative } from "date-fns";
import { useMemo, useState } from "react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { RuleCreationDialog } from "@/components/rule-creation-dialog";

/**
 *
 */
export const Route = createFileRoute(
  "/(project)/projects/$projectId/databases/$databaseId/",
)({
  component: RouteComponent,
  loader: async ({ context, params }) => {
    const { queryClient } = context as RouterContext;
    await queryClient.ensureQueryData(
      databaseDetailQuery(params.projectId!, params.databaseId!),
    );
  },
});

function RouteComponent() {
  const { projectId, databaseId } = Route.useParams();
  const navigate = useNavigate();
  const { data: database } = useGetDatabase(projectId, databaseId);
  const { data: rules, isLoading: rulesLoading } = useGetRules(
    projectId,
    databaseId,
    "",
  );
  const syncSchema = useSyncDatabaseSchema();
  const deleteDatabase = useDeleteDatabase();
  const deleteRule = useDeleteRule();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [ruleToDelete, setRuleToDelete] = useState<Rule | null>(null);
  const [entitySheetEntity, setEntitySheetEntity] =
    useState<DataBaseEntityWithRules | null>(null);

  const rulesByEntityId = useMemo(
    () =>
      (rules ?? []).reduce<Record<string, Rule[]>>((acc, rule) => {
        const entityId = rule.databaseEntity.id;
        if (!acc[entityId]) {
          acc[entityId] = [];
        }
        acc[entityId].push(rule);
        return acc;
      }, {}),
    [rules],
  );

  const toEntitySheetType = (
    type: string,
  ): "column" | "table" | "schema" | "generic" => {
    if (type === "column" || type === "table" || type === "schema") {
      return type;
    }
    return "generic";
  };

  const handleSyncSchema = () => {
    toast.promise(syncSchema.mutateAsync({ projectId, databaseId }), {
      loading: "Syncing database schema...",
      success: "Database schema synced requested successfully!",
      error: (err) => `Error syncing database schema: ${err.message}`,
    });
  };

  const handleDelete = () => {
    toast.promise(
      deleteDatabase.mutateAsync({ projectId, databaseId }).then(() => {
        navigate({ to: "/projects/$projectId", params: { projectId } });
      }),
      {
        loading: "Deleting database...",
        success: "Database deleted successfully!",
        error: (err) => `Error deleting database: ${err.message}`,
      },
    );
  };

  const handleDeleteRule = () => {
    if (!ruleToDelete) {
      return;
    }

    toast.promise(
      deleteRule
        .mutateAsync({ projectId, databaseId, ruleId: ruleToDelete.id })
        .then(() => {
          setRuleToDelete(null);
        }),
      {
        loading: "Deleting rule...",
        success: `Rule "${ruleToDelete.name}" deleted successfully!`,
        error: (err) => `Error deleting rule: ${err.message}`,
      },
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold flex items-center gap-2">
          <DatabaseIcon />
          {database?.displayName}
        </h1>
        <div className="flex gap-2">
          {database && (
            <DatabaseUpdateDialog
              projectId={projectId}
              database={database}
              trigger={
                <Button variant="outline">
                  Edit
                  <PencilIcon className="ml-2" />
                </Button>
              }
            />
          )}
          <Button onClick={handleSyncSchema}>
            Sync Schema
            <ArrowsCounterClockwiseIcon className="ml-2" />
          </Button>
          <Button
            variant="destructive"
            onClick={() => setShowDeleteDialog(true)}
          >
            Delete
            <TrashIcon className="ml-2" />
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Database Details</CardTitle>
            <CardDescription>
              Connection and configuration information
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="flex justify-between">
              <span className="font-medium">Database Name:</span>
              <span>{database?.dbName}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">Type:</span>
              <Badge variant="outline">{database?.databaseType}</Badge>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">Host:</span>
              <span>{database?.host}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">Port:</span>
              <span>{database?.port}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">Username:</span>
              <span>{database?.username}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">Last synced at:</span>
              <span>
                {database?.lastSyncedAt &&
                  formatRelative(database?.lastSyncedAt, new Date())}
              </span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex items-center justify-between">
            <div>
              <CardTitle>Rules</CardTitle>
              <CardDescription>Database rules and constraints</CardDescription>
            </div>
            <RuleCreationDialog projectId={projectId} databaseId={databaseId} />
          </CardHeader>
          <CardContent>
            {rulesLoading ? (
              <div className="text-sm text-muted-foreground">
                Loading rules...
              </div>
            ) : rules && rules.length > 0 ? (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Entity Type</TableHead>
                    <TableHead>Entity</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rules.map((rule: Rule) => (
                    <TableRow key={rule.id}>
                      <TableCell>{rule.name}</TableCell>
                      <TableCell>
                        <Badge variant="secondary">{rule.ruleType}</Badge>
                      </TableCell>
                      <TableCell>{rule.databaseEntity.type}</TableCell>
                      <TableCell>{rule.databaseEntity.name}</TableCell>
                      <TableCell className="text-right">
                        <DropdownMenu>
                          <DropdownMenuTrigger
                            render={
                              <Button
                                variant="ghost"
                                size="icon"
                                aria-label="Rule actions"
                              >
                                <DotsThreeOutlineVerticalIcon />
                              </Button>
                            }
                          />
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={() => {
                                const entity = rule.databaseEntity;
                                setEntitySheetEntity({
                                  ...entity,
                                  rules: rulesByEntityId[entity.id] ?? [],
                                });
                              }}
                            >
                              <InfoIcon />
                              Open entity
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              variant="destructive"
                              onClick={() => setRuleToDelete(rule)}
                            >
                              <TrashIcon />
                              Delete rule
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
                <TableCaption>{rules.length} rules</TableCaption>
              </Table>
            ) : (
              <div className="text-sm text-muted-foreground">
                No rules found
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {database?.lastSyncError && (
        <Alert variant={"destructive"}>
          <AlertTitle>Error syncing your database</AlertTitle>
          <AlertDescription className="flex flex-col">
            <div>
              <span className=" font-bold">Attempts: </span>
              {database.failedConnectionAttempts}.
            </div>
            <div>
              <span className=" font-bold">Error message: </span>
              {database.lastSyncError}
            </div>
            <div>
              <span className=" font-bold">Last attempted at: </span>
              {database.lastSyncedAt &&
                formatRelative(database.lastSyncedAt, new Date())}
            </div>
          </AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Database Schema</CardTitle>
          <CardDescription>
            Visual representation of the database structure
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DatabaseSchemaFlow projectId={projectId} databaseId={databaseId} />
        </CardContent>
      </Card>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Database</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete "{database?.displayName}"? This
              action cannot be undone and will permanently remove all database
              entities, rules, and schema information.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog
        open={!!ruleToDelete}
        onOpenChange={(open) => {
          if (!open) {
            setRuleToDelete(null);
          }
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Rule</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete the rule "{ruleToDelete?.name}"?
              This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteRule.isPending}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteRule}
              disabled={deleteRule.isPending}
              isLoading={deleteRule.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete Rule
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {entitySheetEntity && (
        <DataBaseEntitySheet
          entity={entitySheetEntity}
          type={toEntitySheetType(entitySheetEntity.type)}
          open={!!entitySheetEntity}
          onOpenChange={(open) => {
            if (!open) {
              setEntitySheetEntity(null);
            }
          }}
          hideTrigger
        />
      )}
    </div>
  );
}
