import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
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
import { PlusIcon, TrashIcon } from "@phosphor-icons/react";
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyTitle,
} from "@/components/ui/empty";
import { Spinner } from "@/components/ui/spinner";
import {
  useCreateAccessToken,
  useDeleteAccessToken,
  useGetAccessTokens,
} from "@/api/access-tokens.ts";
import { toast } from "sonner";
import { AccessTokenCreationDialog } from "./access-token-creation-dialog.tsx";
import CopyButton from "@/components/ui/copy-button.tsx";

interface AccessTokenManagerProps {
  projectId: string;
}

export function AccessTokenManager({ projectId }: AccessTokenManagerProps) {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [createdToken, setCreatedToken] = useState<string | null>(null);

  const { data: apiKeys, isPending } = useGetAccessTokens(projectId);
  const deleteApiKeyMutation = useDeleteAccessToken();
  const createApiKeyMutation = useCreateAccessToken();

  const handleCreateApiKey = async (name: string, expiresAt?: Date) => {
    toast.promise(
      createApiKeyMutation.mutateAsync({
        projectId,
        dto: {
          name,
          expiresAt: expiresAt ? new Date(expiresAt).toISOString() : null,
        },
      }),
      {
        loading: "Creating Access token...",
        success: (res) => {
          setCreatedToken(res.token);
          return "Access token created successfully!";
        },
        error: (res) =>
          "Error: " +
          (res.data?.message ?? "Failed to create Access token. Please try again."),
      },
    );
  };

  const handleDeleteAccessToken = async (accessTokenId: string) => {
    toast.promise(deleteApiKeyMutation.mutateAsync({ projectId, accessTokenId: accessTokenId }), {
      loading: "Deleting Access token...",
      success: "Access token deleted successfully!",
      error: "Failed to delete Access token. Please try again.",
    });
  };

  return (
    <div className="container">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold">CI/CD access tokens</h1>
          <p className="text-muted-foreground">
            Manage access tokens for CI/CD container access to this project
          </p>
        </div>
        <AccessTokenCreationDialog
          open={isCreateDialogOpen}
          onOpenChange={setIsCreateDialogOpen}
          onCreate={handleCreateApiKey}
        />
      </div>

      {createdToken && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Access token Created Successfully!</CardTitle>
            <CardDescription>
              Please copy your new access token now. For security reasons, this is
              the only time it will be displayed.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className=" p-4 rounded-md border flex items-center justify-between gap-4">
              <pre className="text-sm font-mono break-all">
                {createdToken
                  ? createdToken
                  : "Error: Access token token not available. Please check if the key was created successfully."}
              </pre>
              <CopyButton text={createdToken || ""} />{" "}
            </div>
          </CardContent>
        </Card>
      )}

      <div className="grid gap-4">
        {isPending ? (
          <div className="flex items-center justify-center py-12">
            <Spinner className="h-8 w-8" />
          </div>
        ) : apiKeys?.length === 0 ? (
          <Empty>
            <EmptyHeader>
              <EmptyTitle>No CI/CD Access tokens found</EmptyTitle>
              <EmptyDescription>
                Create your first access token to allow CI/CD containers to access
                this project's resources.
              </EmptyDescription>
            </EmptyHeader>
            <EmptyContent>
              <Button onClick={() => setIsCreateDialogOpen(true)}>
                <PlusIcon className="mr-2 h-4 w-4" />
                Create Your First CI/CD access token
              </Button>
            </EmptyContent>
          </Empty>
        ) : (
          apiKeys?.map((apiKey) => (
            <Card key={apiKey.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="flex items-center gap-2">
                      {apiKey.name}
                    </CardTitle>
                    <CardDescription>
                      Created {new Date(apiKey.createdAt).toLocaleDateString()}{" "}
                      • Expires:{" "}
                      {apiKey.expiresAt
                        ? new Date(apiKey.expiresAt).toLocaleDateString()
                        : "Never"}
                    </CardDescription>
                  </div>
                  <div className="flex items-center gap-2">
                    <AlertDialog>
                      <AlertDialogTrigger
                        render={
                          <Button variant="outline" size="sm">
                            <TrashIcon className="h-4 w-4" />
                          </Button>
                        }
                      ></AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Delete access token</AlertDialogTitle>
                          <AlertDialogDescription>
                            Are you sure you want to delete this Access token? This
                            action cannot be undone.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction
                            onClick={() => handleDeleteAccessToken(apiKey.id)}
                          >
                            Delete
                          </AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <Label className="text-sm font-medium">Status</Label>
                  <div className="flex items-center gap-2 mt-1">
                    {apiKey.expiresAt ? (
                      <span className="text-xs text-muted-foreground">
                        {new Date(apiKey.expiresAt) < new Date()
                          ? "Expired"
                          : "Valid"}
                      </span>
                    ) : (
                      <span className="text-xs text-muted-foreground">
                        Valid
                      </span>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  );
}
