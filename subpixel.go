// +build !windows

// 22 september 2014

package ndraw

// both cairo and Mac OS X use subpixel rendering, so we need to adjust coordinates for stroking
// thanks to ebassi, pq, and otaylor in irc.freenode.net/#cairo for clearing up my confusion/ignorance of the issues of subpixel rendering
// TODO windows

// adjustment should be +1 or -1 for origins and sizes, respectively
func subpixelAdjust(x0 int, y0 int, x1 int, y1 int, thickness uint, adjustment float64) (float64, float64, float64, float64) {
	fx0, fy0, fx1, fy1 := float64(x0), float64(y0), float64(x1), float64(y1)
	if thickness % 2 == 0 {
		// even thickness; no adjustment needed
		return fx0, fy0, fx1, fy1
	}
	adj := 0.5 * adjustment
	if x1 - x0 == 0 {
		// vertical line; adjust x
		return fx0 + adj, fy0, fx1 + adj, fy1
	}
	if y1 - y0 == 0 {
		// horizontal line; adjust y
		return fx0, fy0 + adj, fx1, fy1 + adj
	}
	// don't bother trying; the subpixel renderer will antialias lines anyway
	return fx0, fy0, fx1, fy1
}
