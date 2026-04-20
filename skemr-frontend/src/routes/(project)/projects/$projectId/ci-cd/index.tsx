import { useGetDatabases } from "@/api/database";
import { useGetProject } from "@/api/project";
import { GithubIcon, GitlabIcon } from "@/assets/icons";
import CopyButton from "@/components/ui/copy-button";
import { Field, FieldGroup, FieldTitle } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { toast } from "sonner";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Light as SyntaxHighlighter } from "react-syntax-highlighter";
import atelierSavannaDark from "react-syntax-highlighter/dist/esm/styles/hljs/atelier-savanna-dark";
import yaml from "react-syntax-highlighter/dist/esm/languages/hljs/yaml";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";
import { useCreateAccessToken } from "@/api/access-tokens.ts";
import { PlusIcon } from "@phosphor-icons/react";
import { AccessTokenCreationDialog } from "@/components/accesstoken/access-token-creation-dialog.tsx";

SyntaxHighlighter.registerLanguage("yaml", yaml);

export const Route = createFileRoute("/(project)/projects/$projectId/ci-cd/")({
  component: RouteComponent,
});

const vcs = [
  {
    name: "Gitlab CI",
    icon: <GitlabIcon />,
  },
  {
    name: "Github Actions",
    icon: <GithubIcon />,
    comingSoon: true,
  },
];
  
function RouteComponent() {
  const createApiKeyMutation = useCreateAccessToken();
  const [selected, setSelected] = useState("Gitlab CI");
  const [selectedDatabase, setSelectedDatabase] = useState<string | null>(null);
  const [migrationFilesDir, setMigrationFilesDir] = useState<string | null>(
    null,
  );
  const [token, setToken] = useState<string | null>(null);
  const [selfHosted, setSelfHosted] = useState(false);
  const [hostUri, setHostUri] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const selectedOption = vcs.find((v) => v.name === selected);

  const { projectId } = Route.useParams();

  const { data: project } = useGetProject(projectId);

  const { data: databases } = useGetDatabases(projectId, "");

  const gitlabDockerYaml = `run-skemr:
    image: 
      name: walmaa/skemr-cli:latest
      pull_policy: always
      entrypoint: [""]
    stage: test
    script:
      - skemr-cli validate --projectId ${project?.id} --databaseId ${selectedDatabase ?? "<fill-your-database-id>"} --migrationFilesDir ${migrationFilesDir ?? "<fill-path-to-migrations>"} --token ${token ?? "<fill-your-access-token>"}${selfHosted ? ` --host ${hostUri ?? "<fill-your-host-uri>"}` : ""}
  `;
  const handlePlatformChange = (value: string | null) => {
    if (!value) {
      return null;
    }

    setSelected(value);
  };

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
          setToken(res.token);
          return "Access token created successfully!";
        },
        error: (res) =>
          "Error: " +
          (res.data?.message ?? "Failed to create Access token. Please try again."),
      },
    );
    setIsCreateDialogOpen(false);
  };

  return (
    <div className="space-y-4 container max-w-3xl">
      <h1 className="text-3xl font-bold">CI/CD Integration</h1>

      <Select value={selected} onValueChange={handlePlatformChange}>
        <SelectTrigger>
          <div className="flex items-center gap-2">
            {selectedOption?.icon}
            <SelectValue />
          </div>
        </SelectTrigger>
        <SelectContent
          className={"w-60"}
          alignItemWithTrigger={false}
          align="start"
        >
          {vcs.map((option) => (
            <SelectItem
              key={option.name}
              value={option.name}
              disabled={option.comingSoon}
            >
              <div className="flex items-center gap-2">
                {option.icon}
                {option.name}
                {option.comingSoon && (
                  <span className="ml-auto text-xs text-muted-foreground">
                    Coming Soon
                  </span>
                )}
              </div>
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <FieldGroup>
        <Field>
          <FieldTitle>Database</FieldTitle>
          <Select
            value={selectedDatabase ?? ""}
            onValueChange={(v) => setSelectedDatabase(v || null)}
          >
            <SelectTrigger className="w-full">
              <SelectValue>
                {databases?.find((db) => db.id === selectedDatabase)
                  ?.displayName || "Select a database"}
              </SelectValue>
            </SelectTrigger>
            <SelectContent alignItemWithTrigger={false} align="start">
              {databases?.map((db) => (
                <SelectItem key={db.id} value={db.id}>
                  {db.displayName}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </Field>

        <Field>
          <FieldTitle>Migration files directory</FieldTitle>
          <Input
            placeholder="e.g. ./migrations"
            value={migrationFilesDir ?? ""}
            onChange={(e) => setMigrationFilesDir(e.target.value || null)}
          />
        </Field>

        <Field>
          <FieldTitle>Access token</FieldTitle>
          <InputGroup>
            <InputGroupInput
              placeholder="Your Skemr access token"
              value={token ?? ""}
              onChange={(e) => setToken(e.target.value || null)}
            />
            <InputGroupAddon align="inline-end">
              <AccessTokenCreationDialog
                open={isCreateDialogOpen}
                onOpenChange={setIsCreateDialogOpen}
                onCreate={handleCreateApiKey}
                trigger={
                  <InputGroupButton>
                    <PlusIcon />
                  </InputGroupButton>
                }
              />
            </InputGroupAddon>
          </InputGroup>
        </Field>

        <Field orientation="horizontal">
          <Switch
            checked={selfHosted}
            onCheckedChange={setSelfHosted}
            id="self-hosted-toggle"
          />
          <FieldTitle>
            <label htmlFor="self-hosted-toggle" className="cursor-pointer">
              Using self-hosted instance?
            </label>
          </FieldTitle>
        </Field>

        {selfHosted && (
          <Field>
            <FieldTitle>Host URI</FieldTitle>
            <Input
              placeholder="e.g. https://skemr.example.com"
              value={hostUri ?? ""}
              onChange={(e) => setHostUri(e.target.value || null)}
            />
          </Field>
        )}
      </FieldGroup>

      <div className="prose prose-neutral prose-sm dark:prose-invert">
        <h2>Integrating Skemr to your Gitlab CI pipeline</h2>
        <ul className=" list-decimal">
          <li>
            In your Gitlab repository, create or edit the{" "}
            <code>.gitlab-ci.yml</code> file.
          </li>
          <li>
            Add the provided YAML snippet to your CI/CD configuration. This
            snippet defines a job that uses the Skemr CLI Docker image to
            validate your database schema.
          </li>
          <li>
            Save and commit the changes to your repository. The next time your
            pipeline runs, it will execute the Skemr validation job as part of
            your CI/CD process.
          </li>
        </ul>
      </div>

      <div className="rounded overflow-hidden border text-sm">
        <SyntaxHighlighter language={"yaml"} style={atelierSavannaDark}>
          {gitlabDockerYaml}
        </SyntaxHighlighter>
        <div>
          <CopyButton text={gitlabDockerYaml} />
        </div>
      </div>
    </div>
  );
}
