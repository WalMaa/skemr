import { ApiKeyManager } from "@/components/api-key-manager";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/(project)/projects/$projectId/api-keys/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId } = Route.useParams();
  return <ApiKeyManager projectId={projectId} />;
}
