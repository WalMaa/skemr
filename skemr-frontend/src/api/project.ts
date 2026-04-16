import { apiClient, type ApiError } from "@/lib/api-client";
import { logger } from "@/lib/logger";
import { projectQueryKeys } from "@/lib/queryKeys";
import type { Project } from "@/types/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export function useGetProjects() {
  return useQuery<Project[], ApiError>({
    queryKey: projectQueryKeys.lists(),
    queryFn: async () => apiClient.get<Project[]>("/projects"),
  });
}

export function useGetProject(id: string) {
  return useQuery<Project, ApiError>(projectDetailQuery(id));
}

export function projectDetailQuery(id: string) {
  return {
    queryKey: projectQueryKeys.detail(id),
    queryFn: async () => apiClient.get<Project>(`/projects/${id}`),
  };
}

export function useCreateProject() {
  const queryClient = useQueryClient();
  
  return useMutation<Project, ApiError, { name: string }>({
    mutationFn: async (projectData) =>
      apiClient.post<Project>("/projects", projectData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.lists() });
    },
    onError: (error) => {
      logger.error(error, "Failed to create project");
    },
  });
}

export function useDeleteProject() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, string>({
    mutationFn: async (projectId) =>
      apiClient.delete<void>(`/projects/${projectId}`),
    onSuccess: (_data, projectId) => {
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.lists() });
      queryClient.removeQueries({ queryKey: projectQueryKeys.detail(projectId) });
    },
  });
}