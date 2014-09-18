// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

HPEN newPen(struct xpen *pen)
{
	HPEN p;

	p = ExtCreatePen(pen.style, pen.width, &pen.brush, pen.nSegments, pen.segments);
	if (p == NULL)
		xpanic("error creating pen", GetLastError();
	return p;
}
