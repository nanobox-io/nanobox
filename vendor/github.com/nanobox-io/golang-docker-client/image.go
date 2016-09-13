package docker

import (
	"errors"
	"io"
	"io/ioutil"

	dockType "github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

type Image struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	RepoTags    []string `json:"repo_tags"`
	Size        int64    `json:"size"`
	VirtualSize int64    `json:"virtual_size"`
	Status      string   `json:"status"`
}

// ImageExists
func ImageExists(name string) bool {
	images, err := ImageList()
	if err != nil {
		return false
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == name+":latest" || tag == name {
				return true
			}
		}
	}

	return false
}

// pull any new image
func ImagePull(image string, output io.Writer) (Image, error) {

	privilegeFunc := func() (string, error) {
		return "", errors.New("no privilege function defined")
	}
	pullOptions := dockType.ImagePullOptions{
		PrivilegeFunc: privilegeFunc,
	}

	ctx := context.Background()
	rc, err := client.ImagePull(ctx, image, pullOptions)
	if err != nil {
		return Image{}, err
	}
	defer rc.Close()

	if output == nil {
		output = ioutil.Discard
	}

	io.Copy(output, rc)

	return ImageInspect(image)
}

// list the images i have cached on the server
func ImageList() ([]Image, error) {
	imgs := []Image{}
	dockImages, err := client.ImageList(context.Background(), dockType.ImageListOptions{})
	if err != nil {
		return imgs, err
	}
	for _, dockImage := range dockImages {
		img := Image{
			ID:          dockImage.ID,
			RepoTags:    dockImage.RepoTags,
			Size:        dockImage.Size,
			VirtualSize: dockImage.VirtualSize,
			Status:      "complete",
		}
		if len(img.RepoTags) > 0 {
			img.Slug = img.RepoTags[0]
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
}

func ImageInspect(imageID string) (Image, error) {
	// ignore the raw part of the image inspect
	dockInspect, _, err := client.ImageInspectWithRaw(context.Background(), imageID)
	img := Image{
		ID:          dockInspect.ID,
		RepoTags:    dockInspect.RepoTags,
		Size:        dockInspect.Size,
		VirtualSize: dockInspect.VirtualSize,
		Status:      "complete",
	}
	if len(img.RepoTags) > 0 {
		img.Slug = img.RepoTags[0]
	}
	return img, err
}

func ImageRemove(imageID string) error {
	_, err := client.ImageRemove(context.Background(), imageID, dockType.ImageRemoveOptions{Force: true, PruneChildren: true})
	return err
}
