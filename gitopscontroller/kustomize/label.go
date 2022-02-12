package kustomize

import "sigs.k8s.io/kustomize/api/types"

func SetLabel(kustomizationFilePath string, metadata map[string]string) error {
	mf, err := NewKustomizationFile(kustomizationFilePath)
	if err != nil {
		return err
	}

	m, err := mf.Read()
	if err != nil {
		return err
	}

	err = setLabels(m, metadata)
	if err != nil {
		return err
	}

	return mf.Write(m)
}

func setLabels(m *types.Kustomization, metadata map[string]string) error {
	if m.CommonLabels == nil {
		m.CommonLabels = make(map[string]string)
	}
	return writeToMap(metadata, m.CommonLabels)
}

func writeToMap(metadata map[string]string, m map[string]string) error {
	for k, v := range metadata {
		m[k] = v
	}
	return nil
}
