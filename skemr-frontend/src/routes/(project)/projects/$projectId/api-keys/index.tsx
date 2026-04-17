import { createFileRoute } from "@tanstack/react-router";
import { ApiKeyManager } from "@/components/accesstoken/api-key-manager.tsx";

export const Route = createFileRoute(
  "/(project)/projects/$projectId/api-keys/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId } = Route.useParams();
  return <ApiKeyManager projectId={projectId} />;
}
