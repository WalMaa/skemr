package databases

import "github.com/walmaa/skemr-common/models"

type DatabaseCreationDto struct {
	DisplayName  string              `json:"displayName" validate:"required"`
	DbName       *string             `json:"dbName"`
	Username     *string             `json:"username"`
	Password     *string             `json:"password"`
	Host         *string             `json:"host"`
	Port         int32               `json:"port"`
	SslMode      *string             `json:"sslMode"`
	DatabaseType models.DatabaseType `json:"databaseType" validate:"required,oneof=postgres"`
}

type DatabaseUpdateDto struct {
	DisplayName  *string              `json:"displayName"`
	DbName       *string              `json:"dbName"`
	Username     *string              `json:"username"`
	Password     *string              `json:"password"`
	Host         *string              `json:"host"`
	Port         *int32               `json:"port"`
	SslMode      *string              `json:"sslMode"`
	DatabaseType *models.DatabaseType `json:"databaseType"`
}
