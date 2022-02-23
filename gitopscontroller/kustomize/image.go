package kustomize

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"sigs.k8s.io/kustomize/api/types"
)

const separator = "="

var pattern = regexp.MustCompile(`^(.*):([a-zA-Z0-9._-]*|\*)$`)
var preserveSeparator = "*"
var (
	errImageInvalidArgs = errors.New("invalid format of image")
)

type overwrite struct {
	name   string
	digest string
	tag    string
}

func SetImage(kustomizationFilePath string, image string, newImage string) error {
	mf, err := NewKustomizationFile(kustomizationFilePath)
	if err != nil {
		return err
	}

	m, err := mf.Read()
	if err != nil {
		return err
	}

	img, err := parse(fmt.Sprintf("%s=%s", image, newImage))
	if err != nil {
		return err
	}

	imageMap := map[string]types.Image{}
	imageMap[img.Name] = img

	for _, im := range m.Images {
		if im.Name == img.Name {
			if img.NewName == preserveSeparator {
				img = replaceNewName(img, im.NewName)
			}
			if img.NewTag == preserveSeparator {
				img = replaceNewTag(img, im.NewTag)
			}
			if img.Digest == preserveSeparator {
				img = replaceDigest(img, im.Digest)
			}
			imageMap[im.Name] = img
			continue
		}
		imageMap[im.Name] = im
	}

	var images []types.Image
	for _, im := range imageMap {
		if im.NewName == preserveSeparator {
			im = replaceNewName(im, "")
		}
		if im.NewTag == preserveSeparator {
			im = replaceNewTag(im, "")
		}
		if im.Digest == preserveSeparator {
			im = replaceDigest(im, "")
		}
		images = append(images, im)
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].Name < images[j].Name
	})

	m.Images = images

	return mf.Write(m)
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
