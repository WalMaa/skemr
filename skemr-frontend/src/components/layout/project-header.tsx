import { Link, useParams } from "@tanstack/react-router";
import { SidebarTrigger } from "@/components/ui/sidebar.tsx";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator
} from "@/components/ui/breadcrumb.tsx";
import { useGetProject } from "@/api/project.ts";
import { useGetDatabase } from "@/api/database.ts";

export function ProjectHeader() {

  const { projectId, databaseId } = useParams({ strict: false });

  const { data: project } = useGetProject(projectId || "");
  const { data: database } = useGetDatabase(projectId || "", databaseId || "");

  return (
    <header className="flex h-16 z-100 shrink-0 sticky top-0 bg-background items-center gap-4 border-b px-4">
      <SidebarTrigger/>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink render={
              <Link
                to="/"
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                Projects
              </Link>
            }/>
          </BreadcrumbItem>
          <BreadcrumbSeparator/>
          {
            project && (
              <BreadcrumbItem>
                <BreadcrumbLink render={
                  <Link to={ "/projects/$projectId" } params={ { projectId: project?.id || "" } } className="font-medium">
                    { project?.name }
                  </Link>
                }/>
              </BreadcrumbItem>

            )
          }
          {
            database && (
              <>
                <BreadcrumbSeparator/>
                <BreadcrumbItem>
                  <BreadcrumbLink render={
                    <Link to={ "/projects/$projectId/databases/$databaseId" } params={ { projectId: project?.id || "", databaseId: database?.id || "" } } className="font-medium">
                      { database?.displayName }
                    </Link>
                  }/>
                </BreadcrumbItem>
              </>
            )
          }
        </BreadcrumbList>
      </Breadcrumb>
    </header>
  );
}
