import { apiClient, type ApiError } from "@/lib/api-client";
import {
  databaseEntityQueryKeys,
  databaseQueryKeys,
  ruleQueryKeys,
} from "@/lib/queryKeys";
import type {
  Database,
  DatabaseCreationDto,
  DatabaseUpdateDto,
} from "@/types/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export function useGetDatabases(projectId: string, filters: string) {
  return useQuery<Database[]>(databaseListQuery(projectId, filters));
}

export function useGetDatabase(projectId: string, databaseId: string) {
  return useQuery<Database>(databaseDetailQuery(projectId, databaseId));
}

export function databaseDetailQuery(projectId: string, databaseId: string) {
  return {
    queryKey: databaseQueryKeys.detail(projectId, databaseId),
    queryFn: async () =>
      apiClient.get<Database>(`/projects/${projectId}/databases/${databaseId}`),
  };
}

export function databaseListQuery(projectId: string, filters: string) {
  return {
    queryKey: databaseQueryKeys.list(projectId, filters),
    queryFn: async () =>
      apiClient.get<Database[]>(
        `/projects/${projectId}/databases?filters=${filters}`,
      ),
  };
}

export function useCreateDatabase() {
  const queryClient = useQueryClient();

  return useMutation<
    Database,
    ApiError,
    { projectId: string; databaseData: DatabaseCreationDto }
  >({
    mutationFn: async ({ projectId, databaseData }) =>
      apiClient.post<Database>(
        `/projects/${projectId}/databases`,
        databaseData,
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: databaseQueryKeys.list(variables.projectId, ""),
      });
    },
  });
}

export function useUpdateDatabase() {
  const queryClient = useQueryClient();

  return useMutation<
    Database,
    ApiError,
    { projectId: string; databaseId: string; updateData: DatabaseUpdateDto }
  >({
    mutationFn: async ({ projectId, databaseId, updateData }) =>
      apiClient.patch<Database>(
        `/projects/${projectId}/databases/${databaseId}`,
        updateData,
      ),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: databaseQueryKeys.detail(
          variables.projectId,
          variables.databaseId,
        ),
      });
      queryClient.invalidateQueries({
        queryKey: databaseQueryKeys.list(variables.projectId, ""),
      });
    },
  });
}

export function useSyncDatabaseSchema() {
  const queryClient = useQueryClient();
  return useMutation<void, ApiError, { projectId: string; databaseId: string }>(
    {
      mutationFn: async ({ projectId, databaseId }) =>
        apiClient.post<void>(
          `/projects/${projectId}/databases/${databaseId}/sync`,
        ),
      onSuccess: (_, variables) => {
        // Invalidate the database entities after a short delay to allow the backend to process the sync
        setTimeout(() => {
          queryClient.invalidateQueries({
            queryKey: databaseEntityQueryKeys.list(
              variables.projectId,
              variables.databaseId,
              "",
            ),
          });
          queryClient.invalidateQueries({
            queryKey: databaseQueryKeys.detail(
              variables.projectId,
              variables.databaseId,
            ),
          });

          queryClient.invalidateQueries({
            queryKey: ruleQueryKeys.list(
              variables.projectId,
              variables.databaseId,
              "",
            ),
          });
        }, 2000);
      },
    },
  );
}

export function useDeleteDatabase() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, { projectId: string; databaseId: string }>(
    {
      mutationFn: async ({ projectId, databaseId }) =>
        apiClient.delete<void>(
          `/projects/${projectId}/databases/${databaseId}`,
        ),
      onSuccess: (_, variables) => {
        queryClient.invalidateQueries({
          queryKey: databaseQueryKeys.list(variables.projectId, ""),
        });
        queryClient.removeQueries({
          queryKey: databaseQueryKeys.detail(
            variables.projectId,
            variables.databaseId,
          ),
        });
      },
    },
  );
}
