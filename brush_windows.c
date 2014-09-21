// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

HBRUSH newBrush(LOGBRUSH *brush)
{
	HBRUSH b;

	b = CreateBrushIndirect(brush);
	if (b == NULL)
		xpanic("error creating brush", GetLastError());
	return b;
}

void brushClose(HBRUSH b)
{
	if (DeleteObject(b) == 0)
		xpanic("error closing Brush", GetLastError());
}

HBRUSH brushSelectInto(HBRUSH brush, HDC dc)
{
	HBRUSH prev;

	prev = (HBRUSH) SelectObject(dc, brush);
	if (prev == NULL)
		xpanic("error selecting Brush into Image DC", GetLastError());
	return prev;
}

void brushUnselect(HBRUSH brush, HDC dc, HBRUSH prev)
{
	if (SelectObject(dc, prev) != brush)
		xpanic("error unselecting Brush from Image DC", GetLastError());
}
