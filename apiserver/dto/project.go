package dto

type CreateProjectRequest struct {
	Name string `validate:"required"`
}
