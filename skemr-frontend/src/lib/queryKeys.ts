export const projectQueryKeys = {
  all: ["projects"] as const,
  lists: () => [...projectQueryKeys.all, "list"] as const,
  list: (filters: string) =>
    [...projectQueryKeys.lists(), { filters }] as const,
  detail: (id: string) => [...projectQueryKeys.all, "detail", id] as const,
};

export const databaseQueryKeys = {
  all: ["databases"] as const,
  lists: () => [...databaseQueryKeys.all, "list"] as const,
  list: (projectId: string, filters: string) =>
    [...databaseQueryKeys.lists(), projectId, { filters }] as const,
  detail: (projectId: string, databaseId: string) =>
    [...databaseQueryKeys.all, "detail", projectId, databaseId] as const,
};

export const ruleQueryKeys = {
  all: ["rules"] as const,
  lists: () => [...ruleQueryKeys.all, "list"] as const,
  list: (projectId: string, databaseId: string, filters: string) =>
    [...ruleQueryKeys.lists(), projectId, databaseId, { filters }] as const,
  detail: (projectId: string, databaseId: string, ruleId: string) =>
    [...ruleQueryKeys.all, "detail", projectId, databaseId, ruleId] as const,
};

export const databaseEntityQueryKeys = {
  all: ["databaseEntities"] as const,
  detail: (projectId: string, databaseId: string, entityId: string) =>
    [
      ...databaseEntityQueryKeys.all,
      "detail",
      projectId,
      databaseId,
      entityId,
    ] as const,
  list: (projectId: string, databaseId: string, filters: string) =>
    [
      ...databaseEntityQueryKeys.all,
      "list",
      projectId,
      databaseId,
      { filters },
    ] as const,
  lists: () => [...databaseEntityQueryKeys.all, "list"] as const,
};

export const accessTokenQueryKeys = {
  all: ["accessTokens"] as const,
  list: (projectId: string) =>
    [...accessTokenQueryKeys.all, "list", projectId] as const,
  detail: (projectId: string, apiKeyId: string) =>
    [...accessTokenQueryKeys.all, "detail", projectId, apiKeyId] as const,
};
