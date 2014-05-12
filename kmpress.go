package main

import (
  "code.google.com/p/go.image/tiff"
  "fmt"
  "image"
  "image/color"
  "log"
  "math"
  "math/rand"
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

// TODO: Make sure colors aren't the same.
func MakeRandomCluster(img *Image, nClusters int) (*ClusterSet, error) {
  c := new(ClusterSet)
  c.Colors = make([]byte, nClusters * 3)

  nPixels := len(img.Vals) / 3

  for i := 0; i < nClusters; i++ {
    k := rand.Intn(nPixels)
    c.Colors[3 * i]     = img.Vals[3 * k]
    c.Colors[3 * i + 1] = img.Vals[3 * k + 1]
    c.Colors[3 * i + 2] = img.Vals[3 * k + 2]
  }

  return c, nil
}

func (c *ClusterSet) Converge(img *Image, maxIterations int) {
  nClusters := len(c.Colors) / 3
  nPixels := len(img.Vals) / 3
  colors := img.Vals
  oldc := c.Colors

  for i := 0; i < maxIterations; i++ {
    newAssigned := make([]int, nClusters)
    newSums := make([]int, nClusters * 3)

    for j := 0; j < nPixels; j++ {
      k := getMinDist(colors, oldc, j, nClusters)
      newAssigned[k]++
      newSums[3 * k]     += int(colors[3 * j])
      newSums[3 * k + 1] += int(colors[3 * j + 1])
      newSums[3 * k + 2] += int(colors[3 * j + 2])
    }

    newc := make([]byte, len(oldc))
    
    for j := 0; j < nClusters; j++ {
      if newAssigned[j] != 0 {
        newc[3 * j]     = byte(newSums[3 * j]     / newAssigned[j])
        newc[3 * j + 1] = byte(newSums[3 * j + 1] / newAssigned[j])
        newc[3 * j + 2] = byte(newSums[3 * j + 2] / newAssigned[j])
      }
    }

    oldc = newc
  }

  c.Colors = oldc
}

func getMinDist(colors, oldc []byte, j, nClusters int) int {
  minDist := math.MaxFloat64
  minK := -1

  for k := 0; k < nClusters; k++ {
    d1 := colors[3 * j]     - oldc[3 * k]
    d2 := colors[3 * j + 1] - oldc[3 * k + 1]
    d3 := colors[3 * j + 2] - oldc[3 * k + 2]
    dist := math.Sqrt(float64(d1 * d1 + d2 * d2 + d3 * d3))
    if dist < minDist {
      minDist = dist
      minK = k
    }
  }

  return minK
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

func writeImage(img image.Image, path string) error {
  out, err := os.Create(path)
  if err != nil {
    return err
  }
  defer out.Close()
  tiff.Encode(out, img, &tiff.Options{tiff.Deflate, true})
  return nil
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

func arrayToImage(colors []byte, w, h int) image.Image {
  nPixels := len(colors) / 3
  img := image.NewRGBA(image.Rect(0, 0, w, h))

  for i := 0; i < nPixels; i++ {
    c := color.RGBA{colors[3 * i], colors[3 * i + 1], colors[3 * i + 2], 255}
    img.Set(i % w, i / w, c)
  }

  return img
}

func encodeToCluster(colors []byte, c *ClusterSet) []byte {
  nClusters := len(c.Colors) / 3
  nPixels := len(colors) / 3
  out := make([]byte, nPixels * 3)

  for i := 0; i < nPixels; i++ {
    k := getMinDist(colors, c.Colors, i, nClusters)
    out[3 * i] = c.Colors[3 * k]
    out[3 * i + 1] = c.Colors[3 * k + 1]
    out[3 * i + 2] = c.Colors[3 * k + 2]
  }
  return out
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

  clusterSet.Converge(img, 100)

  colors := encodeToCluster(img.Vals, clusterSet)
  s := img.Image.Bounds().Size()
  img2 := arrayToImage(colors, s.X, s.Y)
  writeImage(img2, "examples/bird_small_done.tiff")

  fmt.Println("Done.")
}
