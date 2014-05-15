#include "voro++.hh"
using namespace voro;

// Set up constants for the container geometry
const double x_min=-5,x_max=5;
const double y_min=-5,y_max=5;
const double z_min=-5,z_max=5;

// Set up the number of blocks that the container is divided into.
const int n_x=6,n_y=6,n_z=6;

int main() {
	// Create a container with the geometry given above, and make it
	// non-periodic in each of the three coordinates. Allocate space for
	// eight particles within each computational block.
	container con(x_min,x_max,y_min,y_max,z_min,z_max,n_x,n_y,n_z,
			false,false,false,8);

	con.import("list");

	// Save the Voronoi network POV-Ray format.
	con.draw_cells_pov("colors_v.pov");

	// Output the particles in POV-Ray format.
	con.draw_particles_pov("colors_p.pov");
}
