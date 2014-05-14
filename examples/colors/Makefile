CXX=g++
CFLAGS=-Wall -ansi -pedantic -O3
E_INC=-I/usr/local/include/voro++
E_LIB=-L/usr/local/include/voro++

all: colors

colors: colors.cc
	$(CXX) $(CFLAGS) $(E_INC) $(E_LIB) -o colors colors.cc -lvoro++

clean:
	rm -f colors colors.png colors_p.pov colors_v.pov

.PHONY: all clean
