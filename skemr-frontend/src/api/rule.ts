import { apiClient, type ApiError } from "@/lib/api-client";
import { ruleQueryKeys } from "@/lib/queryKeys";
import type { Rule, RuleCreationDto } from "@/types/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export function useGetRules(
  projectId: string,
  databaseId: string,
  filters: string,
) {
  return useQuery<Rule[]>(ruleListQuery(projectId, databaseId, filters));
}

export function useGetRule(
  projectId: string,
  databaseId: string,
  ruleId: string,
) {
  return useQuery<Rule, ApiError>(
    ruleDetailQuery(projectId, databaseId, ruleId),
  );
}

export function ruleDetailQuery(
  projectId: string,
  databaseId: string,
  ruleId: string,
) {
  return {
    queryKey: ruleQueryKeys.detail(projectId, databaseId, ruleId),
    queryFn: async () =>
      apiClient.get<Rule>(
        `/projects/${projectId}/databases/${databaseId}/rules/${ruleId}`,
      ),
  };
}

export function ruleListQuery(
  projectId: string,
  databaseId: string,
  filters: string,
) {
  return {
    queryKey: ruleQueryKeys.list(projectId, databaseId, filters),
    queryFn: async () =>
      apiClient.get<Rule[]>(
        `/projects/${projectId}/databases/${databaseId}/rules?filters=${filters}`,
      ),
  };
}

export function useCreateRule() {
  const queryClient = useQueryClient();

  return useMutation<
    Rule,
    ApiError,
    { projectId: string; databaseId: string; ruleData: RuleCreationDto }
  >({
    mutationFn: async ({ projectId, databaseId, ruleData }) =>
      apiClient.post<Rule>(
        `/projects/${projectId}/databases/${databaseId}/rules`,
        ruleData,
      ),
    onSuccess: (_, { projectId, databaseId }) => {
      queryClient.invalidateQueries({
        queryKey: ruleQueryKeys.list(projectId, databaseId, ""),
      });
    },
  });
}

export function useDeleteRule() {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    ApiError,
    { projectId: string; databaseId: string; ruleId: string }
  >({
    mutationFn: async ({ projectId, databaseId, ruleId }) =>
      apiClient.delete<void>(
        `/projects/${projectId}/databases/${databaseId}/rules/${ruleId}`,
      ),
    onSuccess: (_, { projectId, databaseId }) => {
      queryClient.invalidateQueries({
        queryKey: ruleQueryKeys.list(projectId, databaseId, ""),
      });
    },
  });
}
