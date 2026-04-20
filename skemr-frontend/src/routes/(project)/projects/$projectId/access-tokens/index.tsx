import { createFileRoute } from "@tanstack/react-router";
import { AccessTokenManager } from "@/components/accesstoken/access-token-manager.tsx";

export const Route = createFileRoute(
  "/(project)/projects/$projectId/access-tokens/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId } = Route.useParams();
  return <AccessTokenManager projectId={projectId} />;
}
