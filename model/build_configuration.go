package model

type BuildConfiguration struct {
	// Dockerfile instructions for building container image.
	BuildScript string `json:"buildScript"`

	// Directory for running Dockerfile instructions.
	// If not set, the working directory defaults to the root of the repository.
	WorkingDirectory string `json:"workingDirectory"`

	// Command to run to install the application's dependencies.
	// For example, `npm install`.
	InstallCommand string `json:"installCommand"`

	// Command to run to build the application.
	// For example, `npm run build`.
	BuildCommand string `json:"buildCommand"`

	// Directory (relative to the root of the repository) that contains the
	// outputs of the BuildCommand. For example, `dist`.
	OutputDirectory string `json:"outputDirectory"`

	// Command to run to start the application.
	// For example, `npm start`.
	StartCommand string `json:"startCommand"`
}
