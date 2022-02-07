package dto

type GetBranchesRequest struct {
	Url string `validate:"required"`
}

type GetFilesRequest struct {
	Url    string `validate:"required"`
	Branch string `validate:"required"`
}

type GetRawRequest struct {
	Url    string `validate:"required"`
	Branch string `validate:"required"`
	Path   string `validate:"required"`
}
