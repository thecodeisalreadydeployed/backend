package kustomize

import (
	"sigs.k8s.io/kustomize/api/types"
)

const separator = "="

type overwrite struct {
	name   string
	digest string
	tag    string
}

func replaceNewName(image types.Image, newName string) types.Image {
	return types.Image{
		Name:    image.Name,
		NewName: newName,
		NewTag:  image.NewTag,
		Digest:  image.Digest,
	}
}

func replaceNewTag(image types.Image, newTag string) types.Image {
	return types.Image{
		Name:    image.Name,
		NewName: image.NewName,
		NewTag:  newTag,
		Digest:  image.Digest,
	}
}

func replaceDigest(image types.Image, digest string) types.Image {
	return types.Image{
		Name:    image.Name,
		NewName: image.NewName,
		NewTag:  image.NewTag,
		Digest:  digest,
	}
}
