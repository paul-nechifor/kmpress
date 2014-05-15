# KMPress

A program which uses *k*-means clustering to encode an image using fewer colors.

![KMPress](screenshot.png)

This is one of my [Machine Learning][course] homeworks ([see
another][facetrain]). The requirements are from [here][ex9].

## Install

### KMPress binary

    cd kmpress
    export GOPATH="`pwd`/lib"
    mkdir -p "`pwd`/lib"
    go get ./...
    cd ..

### Voronoi

You have to install [Voro++][voro] and [POV-Ray][povray].

Build the executable:

    cd voronoi
    make

## License

MIT
[course]: http://thor.info.uaic.ro/~ciortuz/teaching.html
[facetrain]: https://github.com/paul-nechifor/facetrain
[ex9]: http://openclassroom.stanford.edu/MainFolder/DocumentPage.php?course=MachineLearning&doc=exercises/ex9/ex9.html
[voro]: http://math.lbl.gov/voro++/download/
[povray]: http://www.povray.org/download/
