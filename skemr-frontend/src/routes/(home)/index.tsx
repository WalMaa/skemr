import { useGetProjects } from "@/api/project";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription, EmptyContent } from "@/components/ui/empty";
import { ProjectCreationDialog } from "@/components/project-creation-dialog";
import { createFileRoute, Link } from "@tanstack/react-router";
import { FolderOpen } from "@phosphor-icons/react";

export const Route = createFileRoute("/(home)/")({
  component: Index,
});

function Index() {
  const { data: projects, isPending, error } = useGetProjects();

  if (isPending) return <div className="p-8">Loading...</div>;

  if (error)
    return (
      <div className="p-8 text-destructive">
        An error has occurred: {error.message}
      </div>
    );

  return (
    <div className="p-8">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Projects</h1>
          <p className="text-muted-foreground mt-2">
            Manage and view all your projects
          </p>
        </div>
        {projects && projects.length > 0 && <ProjectCreationDialog />}
      </div>

      {projects && projects.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map((project) => (
            <Card
              key={project.id}
              className="hover:shadow-lg transition-shadow"
            >
              <CardHeader>
                <CardTitle>{project.name}</CardTitle>
              </CardHeader>
              <CardContent className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  render={
                    <Link
                      to="/projects/$projectId"
                      params={{ projectId: String(project.id) }}
                    >
                      View
                    </Link>
                  }
                />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <Empty>
          <EmptyHeader>
            <FolderOpen size={48} weight="duotone" className="mb-2 text-muted-foreground" />
            <EmptyTitle>No projects yet</EmptyTitle>
            <EmptyDescription>
              Get started by creating your first project
            </EmptyDescription>
          </EmptyHeader>
          <EmptyContent>
            <ProjectCreationDialog />
          </EmptyContent>
        </Empty>
      )}
    </div>
  );
}
