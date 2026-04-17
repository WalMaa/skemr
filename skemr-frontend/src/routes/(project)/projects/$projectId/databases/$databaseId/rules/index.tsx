import { createFileRoute } from "@tanstack/react-router";
import { useGetRules, useDeleteRule } from "@/api/rule";
import { useGetDatabase } from "@/api/database";
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
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Spinner } from "@/components/ui/spinner";
import {
  Empty,
  EmptyHeader,
  EmptyTitle,
  EmptyDescription,
} from "@/components/ui/empty";
import { TrashIcon } from "@phosphor-icons/react";
import type { Rule } from "@/types/types";
import { toast } from "sonner";
import { format } from "date-fns";
import { RuleCreationDialog } from "@/components/rules/rule-creation-dialog";

export const Route = createFileRoute(
  "/(project)/projects/$projectId/databases/$databaseId/rules/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId, databaseId } = Route.useParams();
  const { data: database } = useGetDatabase(projectId, databaseId);
  const { data: rules, isLoading } = useGetRules(projectId, databaseId, "");
  const deleteRule = useDeleteRule();

  const handleDeleteRule = (ruleId: string, ruleName: string) => {
    toast.promise(deleteRule.mutateAsync({ projectId, databaseId, ruleId }), {
      loading: "Deleting rule...",
      success: `Rule "${ruleName}" deleted successfully!`,
      error: (err) => `Error deleting rule: ${err.message}`,
    });
  };

  const getRuleBadgeVariant = (ruleType: string) => {
    switch (ruleType) {
      case "locked":
        return "destructive";
      case "warn":
        return "default";
      case "advisory":
        return "secondary";
      case "deprecated":
        return "outline";
      default:
        return "secondary";
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Database Rules</h1>
          <p className="text-muted-foreground mt-2">
            Manage rules and constraints for {database?.displayName}
          </p>
        </div>
        <RuleCreationDialog projectId={projectId} databaseId={databaseId} />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Rules</CardTitle>
          <CardDescription>
            {rules?.length || 0} rule{rules?.length !== 1 ? "s" : ""} configured
          </CardDescription>
        </CardHeader>
        <CardContent>
          {!rules || rules.length === 0 ? (
            <Empty>
              <EmptyHeader>
                <EmptyTitle>No rules found</EmptyTitle>
                <EmptyDescription>
                  Create your first rule to add constraints to your database
                  entities.
                </EmptyDescription>
              </EmptyHeader>
            </Empty>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Entity Type</TableHead>
                  <TableHead>Entity Name</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rules.map((rule: Rule) => (
                  <TableRow key={rule.id}>
                    <TableCell className="font-medium">{rule.name}</TableCell>
                    <TableCell>
                      <Badge variant={getRuleBadgeVariant(rule.ruleType)}>
                        {rule.ruleType}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">
                        {rule.databaseEntity.type}
                      </Badge>
                    </TableCell>
                    <TableCell>{rule.databaseEntity.name}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {format(rule.createdAt, "PPP p")}
                    </TableCell>
                    <TableCell className="text-right">
                      <AlertDialog>
                        <AlertDialogTrigger
                          render={
                            <Button variant="destructive" size="sm">
                              <TrashIcon className="h-4 w-4" />
                            </Button>
                          }
                        ></AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Delete Rule</AlertDialogTitle>
                            <AlertDialogDescription>
                              Are you sure you want to delete the rule "
                              {rule.name}"? This action cannot be undone.
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction
                            isLoading={deleteRule.isPending}
                            disabled={deleteRule.isPending}
                              onClick={() =>
                                handleDeleteRule(rule.id, rule.name)
                              }
                            >
                              Delete
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
