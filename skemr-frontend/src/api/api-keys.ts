import { apiClient, type ApiError } from "@/lib/api-client";
import { apiKeyQueryKeys } from "@/lib/queryKeys";
import type { ApiKey, ApiKeyCreationDto } from "@/types/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export function useGetApiKeys(projectId: string) {
  return useQuery<ApiKey[]>(apiKeyListQuery(projectId));
}

export function useGetApiKey(projectId: string, apiKeyId: string) {
  return useQuery<ApiKey>(apiKeyDetailQuery(projectId, apiKeyId));
}

export function apiKeyDetailQuery(projectId: string, apiKeyId: string) {
  return {
    queryKey: apiKeyQueryKeys.detail(projectId, apiKeyId),
    queryFn: async () =>
      apiClient.get<ApiKey>(`/projects/${projectId}/secrets/${apiKeyId}`),
  };
}

export function apiKeyListQuery(projectId: string) {
  return {
    queryKey: apiKeyQueryKeys.list(projectId),
    queryFn: async () =>
      apiClient.get<ApiKey[]>(`/projects/${projectId}/secrets`),
  };
}

export function useCreateApiKey() {
  const queryClient = useQueryClient();

  return useMutation<
    {token: string;},
    ApiError,
    { projectId: string; dto: ApiKeyCreationDto }
  >({
    mutationFn: async ({ projectId, dto }) =>
      apiClient.post(`/projects/${projectId}/secrets`, dto),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: apiKeyQueryKeys.list(variables.projectId),
      });
    },
  });
}

export function useDeleteApiKey() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, { projectId: string; apiKeyId: string }>({
    mutationFn: async ({ projectId, apiKeyId }) =>
      apiClient.delete(`/projects/${projectId}/secrets/${apiKeyId}`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: apiKeyQueryKeys.list(variables.projectId),
      });
    },
  });
}
