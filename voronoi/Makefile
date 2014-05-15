CXX=g++
CFLAGS=-Wall -ansi -pedantic -O3
E_INC=-I/usr/local/include/voro++
E_LIB=-L/usr/local/include/voro++

all: voronoi

voronoi: voronoi.cc
	$(CXX) $(CFLAGS) $(E_INC) $(E_LIB) -o voronoi voronoi.cc -lvoro++

clean:
	rm -f voronoi voronoi.png voronoi_p.pov voronoi_v.pov

.PHONY: all clean
