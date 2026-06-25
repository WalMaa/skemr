package scope

import (
	"context"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-common/models"
)

type ProjectResolver interface {
	RequireProject(c context.Context, projectId uuid.UUID) (models.Project, error)
}

type DatabaseResolver interface {
	RequireDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID) (models.Database, error)
}

type DatabaseEntityResolver interface {
	RequireDatabaseEntity(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, entityId uuid.UUID) (models.DatabaseEntity, error)
}

type DatabaseScopeResolver interface {
	ProjectResolver
	DatabaseResolver
}

type EntityScopeResolver interface {
	DatabaseResolver
	DatabaseEntityResolver
}

type Resolver interface {
	ProjectResolver
	DatabaseResolver
	DatabaseEntityResolver
}
