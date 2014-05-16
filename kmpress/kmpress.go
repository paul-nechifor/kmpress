package main

// I hate, as much as you do, the fact that I haven't yet updated this code to
// look and perform better.

import (
  "code.google.com/p/go.image/tiff"
  "flag"
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
  img.Image, err = OpenImage(path)
  if err != nil {
    return nil, err
  }
  img.Vals = ImageToArray(img.Image)
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

func PrintIteration(colors []byte, i int) {
  fmt.Println("iteration>>>", i, colors)
}

func (c *ClusterSet) Converge(img *Image, maxIterations int) {
  nClusters := len(c.Colors) / 3
  nPixels := len(img.Vals) / 3
  colors := img.Vals
  oldc := c.Colors

  PrintIteration(oldc, 0)

  for i := 0; i < maxIterations; i++ {
    newAssigned := make([]int, nClusters)
    newSums := make([]int, nClusters * 3)

    for j := 0; j < nPixels; j++ {
      k := GetMinDist(colors, oldc, j, nClusters)
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

    PrintIteration(newc, i + 1)

    // Test convergence.
    modified := false
    for j := 0; j < nClusters; j++ {
      if newc[j] != oldc[j] {
        modified = true
        break
      }
    }

    oldc = newc

    if !modified {
      break
    }
  }

  c.Colors = oldc
}

func GetMinDist(colors, oldc []byte, j, nClusters int) int {
  minDist := math.MaxFloat64
  minK := -1

  for k := 0; k < nClusters; k++ {
    d1 := float64(colors[3 * j])     - float64(oldc[3 * k])
    d2 := float64(colors[3 * j + 1]) - float64(oldc[3 * k + 1])
    d3 := float64(colors[3 * j + 2]) - float64(oldc[3 * k + 2])
    dist := math.Sqrt(d1 * d1 + d2 * d2 + d3 * d3)
    if dist < minDist {
      minDist = dist
      minK = k
    }
  }

  return minK
}

func OpenImage(path string) (image.Image, error) {
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

func WriteImage(img image.Image, path string) error {
  out, err := os.Create(path)
  if err != nil {
    return err
  }
  defer out.Close()
  tiff.Encode(out, img, &tiff.Options{tiff.Deflate, true})
  return nil
}

func ImageToArray(img image.Image) []byte {
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

func ArrayToImage(colors []byte, w, h int) image.Image {
  nPixels := len(colors) / 3
  img := image.NewRGBA(image.Rect(0, 0, w, h))

  for i := 0; i < nPixels; i++ {
    c := color.RGBA{colors[3 * i], colors[3 * i + 1], colors[3 * i + 2], 255}
    img.Set(i % w, i / w, c)
  }

  return img
}

func EncodeToCluster(colors []byte, c *ClusterSet) []byte {
  nClusters := len(c.Colors) / 3
  nPixels := len(colors) / 3
  out := make([]byte, nPixels * 3)

  for i := 0; i < nPixels; i++ {
    k := GetMinDist(colors, c.Colors, i, nClusters)
    out[3 * i] = c.Colors[3 * k]
    out[3 * i + 1] = c.Colors[3 * k + 1]
    out[3 * i + 2] = c.Colors[3 * k + 2]
  }
  return out
}

func Encode(nClusters, maxIterations int, i, o string) {
  img, err := Open(i)
  if err != nil {
    log.Fatal(err)
  }

  clusterSet, err := MakeRandomCluster(img, nClusters)
  if err != nil {
    log.Fatal(err)
  }

  clusterSet.Converge(img, maxIterations)

  EncodeSave(img, clusterSet, o)
}

func EncodeSave(img *Image, clusterSet *ClusterSet, o string) {
  colors := EncodeToCluster(img.Vals, clusterSet)
  s := img.Image.Bounds().Size()
  img2 := ArrayToImage(colors, s.X, s.Y)
  WriteImage(img2, o)
}

func RenderCluster(i, o string) {
  img, err := Open(i)
  if err != nil {
    log.Fatal(err)
  }

  var n int
  fmt.Scanf("%d", &n)
  colors := make([]byte, n)
  for i := 0; i < n; i++ {
    fmt.Scanf("%d", &colors[i])
  }
  fmt.Println("Number", n, colors)

  clusterSet := new(ClusterSet)
  clusterSet.Colors = colors

  EncodeSave(img, clusterSet, o)
}

func main() {
  nClusters := flag.Int("clusters", 16, "number of clusters")
  maxIterations := flag.Int("max", 16, "maximum number of iterations")
  i := flag.String("i", "", "input file")
  o := flag.String("o", "", "output file")
  renderCluster := flag.Bool("renderCluster", false, "")
  flag.Parse()

  if *renderCluster {
    RenderCluster(*i, *o)
    return
  }

  Encode(*nClusters, *maxIterations, *i, *o)
}