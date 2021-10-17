package dto

type ValidateBuildScriptRequest struct {
	BuildScript string `validate:"required" json:"buildScript"`
}
