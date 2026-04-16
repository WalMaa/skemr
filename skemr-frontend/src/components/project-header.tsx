import { Link } from "@tanstack/react-router";
import { SidebarTrigger } from "./ui/sidebar";

interface Project {
  id: string;
  name: string;
  description?: string;
}

interface ProjectHeaderProps {
  project?: Project;
}

export function ProjectHeader({ project }: ProjectHeaderProps) {
  return (
    <header className="flex h-16 shrink-0 items-center gap-4 border-b px-4">
      <SidebarTrigger />
      <nav className="flex items-center gap-2 text-sm">
        <Link
          to="/"
          className="text-muted-foreground hover:text-foreground transition-colors"
        >
          Projects
        </Link>
        <span className="text-muted-foreground">/</span>
        <span className="font-medium">{project?.name || "Loading..."}</span>
      </nav>
    </header>
  );
}
