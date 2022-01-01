package kustomize

import (
	"errors"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/api/types"
)

const separator = "="

var pattern = regexp.MustCompile(`^(.*):([a-zA-Z0-9._-]*|\*)$`)
var (
	errImageNoArgs      = errors.New("no image specified")
	errImageInvalidArgs = errors.New("invalid format of image")
)

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

func parse(arg string) (types.Image, error) {

	// matches if there is an image name to overwrite
	// <image>=<new-image><:|@><new-tag>
	if s := strings.Split(arg, separator); len(s) == 2 {
		p, err := parseOverwrite(s[1], true)
		return types.Image{
			Name:    s[0],
			NewName: p.name,
			NewTag:  p.tag,
			Digest:  p.digest,
		}, err
	}

	// matches only for <tag|digest> overwrites
	// <image><:|@><new-tag>
	p, err := parseOverwrite(arg, false)
	return types.Image{
		Name:   p.name,
		NewTag: p.tag,
		Digest: p.digest,
	}, err
}

// parseOverwrite parses the overwrite parameters
// from the given arg into a struct
func parseOverwrite(arg string, overwriteImage bool) (overwrite, error) {
	// match <image>@<digest>
	if d := strings.Split(arg, "@"); len(d) > 1 {
		return overwrite{
			name:   d[0],
			digest: d[1],
		}, nil
	}

	// match <image>:<tag>
	if t := pattern.FindStringSubmatch(arg); len(t) == 3 {
		return overwrite{
			name: t[1],
			tag:  t[2],
		}, nil
	}

	// match <image>
	if len(arg) > 0 && overwriteImage {
		return overwrite{
			name: arg,
		}, nil
	}
	return overwrite{}, errImageInvalidArgs
}
