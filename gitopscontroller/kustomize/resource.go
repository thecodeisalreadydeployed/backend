package kustomize

func AddResources(kustomizationFilePath string, paths []string) error {
	file, err := NewKustomizationFile(kustomizationFilePath)
	if err != nil {
		return err
	}

	m, err := file.Read()
	if err != nil {
		return err
	}

	for _, resource := range paths {
		if file.GetPath() != resource {
			if !StringInSlice(resource, m.Resources) {
				m.Resources = append(m.Resources, resource)
			}
		}
	}

	return file.Write(m)
}
