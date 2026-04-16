export interface Project {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface Database {
  id: string;
  displayName: string;
  dbName: string;
  username: string;
  host: string;
  port: number;
  databaseType: DatabaseType;
  sslMode?: DatabaseSslMode;
  lastSyncedAt: string | null;
  failedConnectionAttempts: number;
  lastSyncError: string | null;
}

export interface DatabaseCreationDto {
  displayName: string;
  dbName: string;
  username: string;
  password: string;
  host: string;
  port: number;
  databaseType: DatabaseType;
  sslMode?: DatabaseSslMode;
}

export interface DatabaseUpdateDto {
  displayName?: string;
  dbName?: string;
  username?: string;
  password?: string;
  host?: string;
  port?: number;
  databaseType?: DatabaseType;
  sslMode?: DatabaseSslMode;
}

type DatabaseType = "postgres";
type DatabaseSslMode =
  | "disable"
  | "allow"
  | "prefer"
  | "require"
  | "verify-ca"
  | "verify-full";

export interface Rule {
  id: string;
  name: string;
  ruleType: DatabaseRuleType;
  createdAt: string;
  databaseEntity: DatabaseEntity;
}

export type RuleCreationDto = {
  name: string;
  ruleType: DatabaseRuleType;
  databaseEntityId: string;
};

export type DatabaseRuleType = "locked" | "warn" | "advisory" | "deprecated";

type ColumnAttributes = {
  dataType: string;
  default: string | null;
  nullable: string;
  updatable: string;
};

type DatabaseEntityStatus = "active" | "deleted"

export interface DatabaseEntity {
  id: string;
  name: string;
  type: DatabaseEntityType;
  parentId: string | null;
  status: DatabaseEntityStatus;
  deletedAt: string | null;
  firstSeenAt: string;
  createdAt: string;
  attributes: Record<string, string> | ColumnAttributes;
}

export type DatabaseEntityType = "database" | "schema" | "table" | "column";

export interface ApiKey {
  id: string;
  name: string;
  lastUsedAt: string | null;
  expiresAt: string | null;
  createdAt: string;
  updatedAt: string;
}

export type ApiKeyCreationDto = {
  name: string;
  expiresAt: string | null;
};
