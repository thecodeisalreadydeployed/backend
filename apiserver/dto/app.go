package dto

type CreateAppRequest struct {
	ProjectID       string `validate:"required"`
	Name            string `validate:"required"`
	RepositoryURL   string `validate:"required"`
	BuildCommand    string `validate:"required"`
	OutputDirectory string `validate:"required"`
	InstallCommand  string `validate:"required"`
	StartCommand    string `validate:"required"`
}
