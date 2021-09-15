package model

type BuildConfiguration struct {
	// Dockerfile instructions for building container image.
	BuildScript string `json:"build_script"`

	// Flag to specify whether the BuildScript should be parsed.
	ParseBuildScript bool `json:"parse_build_script"`

	// Directory for running Dockerfile instructions.
	// If not set, the working directory defaults to the root of the repository.
	WorkingDirectory string `json:"working_directory"`

	// Command to run to install the application's dependencies.
	// For example, `npm install`.
	InstallCommand string `json:"install_command"`

	// Command to run to build the application.
	// For example, `npm run build`.
	BuildCommand string `json:"build_command"`

	// Directory (relative to the root of the repository) that contains the
	// outputs of the BuildCommand. For example, `dist`.
	OutputDirectory string `json:"output_directory"`

	// Command to run to start the application.
	// For example, `npm start`.
	StartCommand string `json:"start_command"`
}
