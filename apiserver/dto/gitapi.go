package dto

type GetBranchesRequest struct {
	URL string `validate:"required"`
}

type GetFilesRequest struct {
	URL    string `validate:"required"`
	Branch string `validate:"required"`
}

type GetRawRequest struct {
	URL    string `validate:"required"`
	Branch string `validate:"required"`
	Path   string `validate:"required"`
}
