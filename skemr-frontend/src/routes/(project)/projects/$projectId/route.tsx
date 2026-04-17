import { projectDetailQuery} from "@/api/project";
import type { RouterContext } from "../../../__root";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import Cookies from "js-cookie";
import { databaseListQuery, useGetDatabases } from "@/api/database";
import { AppSidebar } from "@/components/layout/app-sidebar.tsx";
import { ProjectHeader } from "@/components/layout/project-header.tsx";

export const Route = createFileRoute("/(project)/projects/$projectId")({
  component: RouteComponent,
  loader: async ({ context, params }) => {
    const { queryClient } = context as RouterContext;
    await queryClient.ensureQueryData(projectDetailQuery(params.projectId!));
    await queryClient.ensureQueryData(databaseListQuery(params.projectId!, ""));
  },
});

function RouteComponent() {
  const defaultOpen = Cookies.get("sidebar_state") === "true";

  const { projectId } = Route.useParams();
  const { data: databases } = useGetDatabases(projectId, "");
  return (
    <SidebarProvider defaultOpen={defaultOpen}>
      <AppSidebar projectId={projectId} databases={databases} />
      <SidebarInset>
        <ProjectHeader />
        <div className=" p-2 md:p-5 overflow-scroll">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
