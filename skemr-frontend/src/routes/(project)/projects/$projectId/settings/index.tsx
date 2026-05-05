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
import { cn } from "@/lib/utils";
import { CheckIcon, WarningIcon } from "@phosphor-icons/react";
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
    <div className="space-y-6 max-w-3xl">
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
          <div className="grid gap-4 sm:grid-cols-3">
            <ThemeOption
              label="Light"
              value="light"
              selected={theme === "light"}
              onSelect={setTheme}
            />
            <ThemeOption
              label="Dark"
              value="dark"
              selected={theme === "dark"}
              onSelect={setTheme}
            />
            <ThemeOption
              label="System settings"
              value="system"
              selected={theme === "system"}
              onSelect={setTheme}
            />
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

type ThemeValue = "light" | "dark" | "system";

type ThemeOptionProps = {
  label: string;
  value: ThemeValue;
  selected: boolean;
  onSelect: (theme: ThemeValue) => void;
};

function ThemeOption({ label, value, selected, onSelect }: ThemeOptionProps) {
  return (
    <button
      type="button"
      className="group space-y-2 text-left"
      onClick={() => onSelect(value)}
      aria-pressed={selected}
    >
      <div
        className={cn(
          "relative h-20 overflow-hidden rounded-lg border bg-background transition-colors group-hover:border-primary/70",
          selected && "border-primary ring-2 ring-primary/20",
        )}
      >
        <ThemePreview value={value} />
        {selected ? (
          <span className="absolute bottom-2 right-2 flex size-5 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-sm">
            <CheckIcon className="size-3.5" weight="bold" />
          </span>
        ) : null}
      </div>
      <div className="text-sm font-medium text-foreground">{label}</div>
    </button>
  );
}

function ThemePreview({ value }: { value: ThemeValue }) {
  if (value === "system") {
    return (
      <div className="flex h-full">
        <PreviewPane mode="light" />
        <PreviewPane mode="dark" />
      </div>
    );
  }

  return <PreviewPane mode={value} />;
}

function PreviewPane({ mode }: { mode: "light" | "dark" }) {
  const dark = mode === "dark";

  return (
    <div
      className={cn(
        "relative flex h-full flex-1 items-start px-6 py-5",
        dark ? "bg-[#141414] text-white" : "bg-white text-[#111111]",
      )}
    >
      <div
        className={cn(
          "absolute left-6 right-0 top-6 h-14 rounded-tl-lg border",
          dark
            ? "border-white/10 bg-[#181818]"
            : "border-black/5 bg-[#f3f3f3]",
        )}
      />
      <span className="relative text-2xl leading-none p-2">Aa</span>
    </div>
  );
}
