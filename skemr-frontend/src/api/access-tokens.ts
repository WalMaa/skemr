import { apiClient, type ApiError } from "@/lib/api-client";
import { accessTokenQueryKeys } from "@/lib/queryKeys";
import type { AccessToken, AccessTokenCreationDto } from "@/types/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export function useGetAccessTokens(projectId: string) {
  return useQuery<AccessToken[]>(accessTokenListQuery(projectId));
}

export function useGetAccessToken(projectId: string, apiKeyId: string) {
  return useQuery<AccessToken>(accessTokenKeyDetailQuery(projectId, apiKeyId));
}

export function accessTokenKeyDetailQuery(projectId: string, apiKeyId: string) {
  return {
    queryKey: accessTokenQueryKeys.detail(projectId, apiKeyId),
    queryFn: async () =>
      apiClient.get<AccessToken>(`/projects/${projectId}/secrets/${apiKeyId}`),
  };
}

export function accessTokenListQuery(projectId: string) {
  return {
    queryKey: accessTokenQueryKeys.list(projectId),
    queryFn: async () =>
      apiClient.get<AccessToken[]>(`/projects/${projectId}/secrets`),
  };
}

export function useCreateAccessToken() {
  const queryClient = useQueryClient();

  return useMutation<
    {token: string;},
    ApiError,
    { projectId: string; dto: AccessTokenCreationDto }
  >({
    mutationFn: async ({ projectId, dto }) =>
      apiClient.post(`/projects/${projectId}/secrets`, dto),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: accessTokenQueryKeys.list(variables.projectId),
      });
    },
  });
}

export function useDeleteAccessToken() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, { projectId: string; accessTokenId: string }>({
    mutationFn: async ({ projectId, accessTokenId }) =>
      apiClient.delete(`/projects/${projectId}/secrets/${accessTokenId}`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: accessTokenQueryKeys.list(variables.projectId),
      });
    },
  });
}
