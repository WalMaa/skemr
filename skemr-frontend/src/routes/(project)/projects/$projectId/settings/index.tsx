import { useMemo, useState } from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useDeleteProject, useGetProject } from "@/api/project";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
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
import { Alert } from "@/components/ui/alert";
import { useTheme } from "@/components/theme-provider";
import { MoonIcon, SunIcon, WarningIcon } from "@phosphor-icons/react";
import { toast } from "sonner";

export const Route = createFileRoute(
  "/(project)/projects/$projectId/settings/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId } = Route.useParams();
  const navigate = useNavigate();
  const { theme, setTheme } = useTheme();
  const { data: project, isLoading } = useGetProject(projectId);
  const { mutateAsync: deleteProject, isPending: isDeleting } =
    useDeleteProject();

  const [confirmationText, setConfirmationText] = useState("");

  const projectName = project?.name ?? "";
  const canDelete = useMemo(
    () => projectName.length > 0 && confirmationText === projectName,
    [confirmationText, projectName],
  );
  const isDarkMode = theme === "dark";

  const handleDeleteProject = async () => {
    if (!projectId || !canDelete) {
      return;
    }
    toast.promise(deleteProject(projectId), {
      loading: "Deleting project...",
      success: () => {
        navigate({ to: "/" });
        return "Project deleted successfully!";
      },
      error: "Failed to delete project. Please try again.",
    });
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-sm text-muted-foreground">
          Manage project information, appearance, and destructive actions.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Appearance</CardTitle>
          <CardDescription>Customize how the interface looks.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between gap-4 rounded-lg border p-4">
            <div className="space-y-1">
              <Label htmlFor="dark-mode-switch">Dark mode</Label>
              <p className="text-xs text-muted-foreground">
                Enable dark mode across the app. This preference is saved on
                this device.
              </p>
            </div>
            <div className="flex items-center gap-2">
              <SunIcon className="text-muted-foreground" />
              <Switch
                id="dark-mode-switch"
                checked={isDarkMode}
                onCheckedChange={(checked) =>
                  setTheme(checked ? "dark" : "light")
                }
                aria-label="Toggle dark mode"
              />
              <MoonIcon className="text-muted-foreground" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Project information</CardTitle>
          <CardDescription>Details about this project.</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-sm text-muted-foreground">Loading...</div>
          ) : (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <div className="text-xs uppercase text-muted-foreground">
                  Name
                </div>
                <div className="text-sm font-medium">
                  {project?.name ?? "-"}
                </div>
              </div>
              <div></div>
              <div>
                <div className="text-xs uppercase text-muted-foreground">
                  Created
                </div>
                <div className="text-sm font-medium">
                  {project?.createdAt
                    ? new Date(project.createdAt).toLocaleString()
                    : "-"}
                </div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <Card className="border-destructive/40">
        <CardHeader>
          <CardTitle>Danger zone</CardTitle>
          <CardDescription>
            Delete this project and all related data. This action cannot be
            undone.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <AlertDialog>
            <AlertDialogTrigger
              render={
                <Button variant="destructive" disabled={!projectName}>
                  Delete project
                </Button>
              }
            />

            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Delete project</AlertDialogTitle>

                <AlertDialogDescription>
                  <Alert className="my-2" variant={"destructive"}>
                    <WarningIcon />
                    This action cannot be undone. This will permanently delete
                    the project and all related data.
                  </Alert>
                  To confirm, please type "<code>{projectName}</code>" below:
                </AlertDialogDescription>
              </AlertDialogHeader>
              <div className="space-y-2">
                <Label htmlFor="project-delete-confirmation">
                  Project name
                </Label>
                <Input
                  id="project-delete-confirmation"
                  placeholder={projectName || "Project name"}
                  value={confirmationText}
                  onChange={(event) => setConfirmationText(event.target.value)}
                />
              </div>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleDeleteProject}
                  disabled={!canDelete || isDeleting}
                >
                  {isDeleting ? "Deleting..." : "Delete"}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </CardContent>
      </Card>
    </div>
  );
}
