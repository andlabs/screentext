// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

HPEN newPen(struct xpen *pen)
{
	HPEN p;

	p = ExtCreatePen(pen->style, pen->width, &pen->brush, pen->nSegments, pen->segments);
	if (p == NULL)
		xpanic("error creating pen", GetLastError());
	return p;
}

void penClose(HPEN p)
{
	if (DeleteObject(p) == 0)
		xpanic("error closing Pen", GetLastError());
}

HPEN penSelectInto(HPEN pen, HDC dc, COLORREF color)
{
	HPEN prev;

	prev = (HPEN) SelectObject(dc, pen);
	if (prev == NULL)
		xpanic("error selecting Pen into Image DC", GetLastError());
	if (SetTextColor(dc, color) == CLR_INVALID)
		xpanic("error selecting text color from Pen into Image DC", GetLastError());
	return prev;
}

void penUnselect(HPEN pen, HDC dc, HPEN prev)
{
	if (SelectObject(dc, prev) != pen)
		xpanic("error unselecting Pen from Image DC", GetLastError());
}
