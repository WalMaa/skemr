package accesstoken

type AccessTokenCreationDto struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	ExpiresAt string `json:"expiresAt" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}
