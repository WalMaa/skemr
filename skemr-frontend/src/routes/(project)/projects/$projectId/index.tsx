import { useGetProject } from "@/api/project";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useGetDatabases } from "@/api/database";
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
import { Button } from "@/components/ui/button";
import {
  Empty,
  EmptyHeader,
  EmptyTitle,
  EmptyDescription,
  EmptyContent,
} from "@/components/ui/empty";
import { Badge } from "@/components/ui/badge";
import { DatabaseIcon, ShieldCheckIcon } from "@phosphor-icons/react";
import { DatabaseCreationDialog } from "@/components/database-creation-dialog";
import { format } from "date-fns";
import { Spinner } from "@/components/ui/spinner";

export const Route = createFileRoute("/(project)/projects/$projectId/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId } = Route.useParams();
  const { data: project } = useGetProject(projectId);
  const { data: databases, isLoading } = useGetDatabases(projectId, "");

  if (!project) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">{project?.name}</h1>
      </div>

      <div className="grid grid-cols-1 gap-4">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Project Summary</CardTitle>
                <CardDescription>
                  Created: {format(project.createdAt, "PPP")}
                </CardDescription>
              </div>
              <DatabaseCreationDialog projectId={projectId} />
            </div>
          </CardHeader>
          <CardContent>
            <h3 className="text-sm font-medium mb-2">Databases</h3>

            {isLoading ? (
              <div className="text-sm text-muted-foreground">Loading...</div>
            ) : !databases || databases.length === 0 ? (
              <Empty>
                <EmptyHeader>
                  <EmptyTitle>No databases found</EmptyTitle>
                  <EmptyDescription>
                    This project does not have any databases yet.
                  </EmptyDescription>
                </EmptyHeader>
                <EmptyContent>
                  <DatabaseCreationDialog
                    projectId={projectId}
                    trigger={
                      <Button>
                        <DatabaseIcon className="mr-2 h-4 w-4" />
                        Add Your First Database
                      </Button>
                    }
                  />
                </EmptyContent>
              </Empty>
            ) : (
              <Table>
                <TableHeader>
                  <tr>
                    <TableHead>Name</TableHead>
                    <TableHead>DB Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Host</TableHead>
                    <TableHead>Port</TableHead>
                    <TableHead>Username</TableHead>
                    <TableHead>Actions</TableHead>
                  </tr>
                </TableHeader>
                <TableBody>
                  {databases.map((db) => (
                    <TableRow key={db.id}>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <DatabaseIcon />
                          <Link
                            to="/projects/$projectId/databases/$databaseId"
                            params={{
                              projectId: projectId,
                              databaseId: db.id,
                            }}
                            className="font-medium text-primary"
                          >
                            {db.displayName}
                          </Link>
                        </div>
                      </TableCell>
                      <TableCell>{db.dbName}</TableCell>
                      <TableCell>
                        <Badge variant="outline">{db.databaseType}</Badge>
                      </TableCell>
                      <TableCell>{db.host}</TableCell>
                      <TableCell>{db.port}</TableCell>
                      <TableCell>{db.username}</TableCell>
                      <TableCell>
                        <div className="flex gap-2 items-center">
                          <Link
                            to="/projects/$projectId/databases/$databaseId"
                            params={{ projectId, databaseId: db.id }}
                            className="text-sm"
                          >
                            <Button size="sm" variant="ghost">
                              Open
                            </Button>
                          </Link>
                          <Button
                            size="sm"
                            variant="outline"
                            render={
                              <Link
                                to="/projects/$projectId/databases/$databaseId/rules"
                                params={{ projectId, databaseId: db.id }}
                                className="text-sm"
                              >
                                <ShieldCheckIcon className="mr-1" /> Rules
                              </Link>
                            }
                          ></Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
                <TableCaption>{databases.length} databases</TableCaption>
              </Table>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
