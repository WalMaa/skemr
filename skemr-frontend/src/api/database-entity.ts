import { apiClient } from "@/lib/api-client";
import { databaseEntityQueryKeys } from "@/lib/queryKeys";
import type { DatabaseEntity } from "@/types/types";
import { useQuery } from "@tanstack/react-query";

export function useGetDatabaseEntities(
  projectId: string,
  databaseId: string,
  filters: string,
) {
  return useQuery<DatabaseEntity[]>(
    databaseEntityListQuery(projectId, databaseId, filters),
  );
}

export function useGetDatabaseEntity(
  projectId: string,
  databaseId: string,
  entityId: string,
) {
  return useQuery<DatabaseEntity>(
    databaseEntityDetailQuery(projectId, databaseId, entityId),
  );
}

export function databaseEntityDetailQuery(
  projectId: string,
  databaseId: string,
  entityId: string,
) {
  return {
    queryKey: databaseEntityQueryKeys.detail(projectId, databaseId, entityId),
    queryFn: async () =>
      apiClient.get<DatabaseEntity>(
        `/projects/${projectId}/databases/${databaseId}/entities/${entityId}`,
      ),
  };
}

export function databaseEntityListQuery(
  projectId: string,
  databaseId: string,
  filters: string,
) {
  return {
    queryKey: databaseEntityQueryKeys.list(projectId, databaseId, filters),
    queryFn: async () =>
      apiClient.get<DatabaseEntity[]>(
        `/projects/${projectId}/databases/${databaseId}/entities?filters=${filters}`,
      ),
  };
}
