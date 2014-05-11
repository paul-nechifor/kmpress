package main

import (
	_ "code.google.com/p/go.image/tiff"
  "fmt"
  "log"
	"image"
	"os"
)

type Image struct {
  Image image.Image
  Vals []byte
}

type ClusterSet struct {
  Colors []byte
}

func Open(path string) (*Image, error) {
  var err error
  img := new(Image)
  img.Image, err = openImage(path)
  if err != nil {
    return nil, err
  }
  img.Vals = imageToArray(img.Image)
  return img, nil
}

func MakeRandomCluster(img *Image, nClusters int) (*ClusterSet, error) {
  c := new(ClusterSet)
  c.Colors = img.Vals
  return c, nil
}

func (c *ClusterSet) Converge(maxIterations int) {
}

func openImage(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func imageToArray(img image.Image) []byte {
	bounds := img.Bounds()
	size := bounds.Size()
	array := make([]byte, size.X*size.Y*3)

	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			array[i] = byte(r >> 8)
			array[i+1] = byte(g >> 8)
			array[i+2] = byte(b >> 8)
			i += 3
		}
	}

	return array
}

func main() {
  img, err := Open("examples/bird_small.tiff")
  if err != nil {
    log.Fatal(err)
  }

  clusterSet, err := MakeRandomCluster(img, 16)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(clusterSet.Colors)
}
