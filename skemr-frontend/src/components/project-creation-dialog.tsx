import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Field, FieldError } from "@/components/ui/field";
import { PlusIcon } from "@phosphor-icons/react";
import { useCreateProject } from "@/api/project";
import { toast } from "sonner";

const projectSchema = z.object({
  name: z.string().min(1, "Project name is required"),
});

type ProjectFormData = z.infer<typeof projectSchema>;

interface ProjectCreationDialogProps {
  trigger?: React.ReactElement;
}

export function ProjectCreationDialog({ trigger }: ProjectCreationDialogProps) {
  const [isOpen, setIsOpen] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<ProjectFormData>({
    resolver: zodResolver(projectSchema),
    defaultValues: {
      name: "",
    },
  });

  const createProjectMutation = useCreateProject();

  const onSubmit = (data: ProjectFormData) => {
    toast.promise(createProjectMutation.mutateAsync(data), {
      loading: "Creating project...",
      success: () => {
        setIsOpen(false);
        reset();
        return "Project created successfully!";
      },
      error: (error) => (`Failed to create project: ${error.message} ${JSON.stringify(error.data)}`),
    });
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger
        render={
          trigger || (
            <Button>
              <PlusIcon className="mr-2 h-4 w-4" />
              New Project
            </Button>
          )
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New Project</DialogTitle>
          <DialogDescription>
            Create a new project to organize your databases and rules.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-4 py-4">
            <Field>
              <Label htmlFor="name">
                Project Name <span className="text-destructive">*</span>
              </Label>
              <Input id="name" {...register("name")} placeholder="My Project" />
              {errors.name && <FieldError>{errors.name.message}</FieldError>}
            </Field>
          </div>
          <DialogFooter>
            <Button type="submit" disabled={createProjectMutation.isPending}>
              {createProjectMutation.isPending
                ? "Creating..."
                : "Create Project"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
