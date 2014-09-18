// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

HBITMAP newBitmap(int dx, int dy, void **ppvBits)
{
	BITMAPINFO bi;
	HBITMAP b;

	ZeroMemory(&bi, sizeof (BITMAPINFO));
	bi.bmiHeader.biSize = sizeof (BITMAPINFOHEADER);
	bi.bmiHeader.biWidth = (LONG) dx;
	bi.bmiHeader.biHeight = -((LONG) dy);                   // negative height to force top-down drawing
	bi.bmiHeader.biPlanes = 1;
	bi.bmiHeader.biBitCount = 32;
	bi.bmiHeader.biCompression = BI_RGB;
	bi.bmiHeader.biSizeImage = (DWORD) (dx * dy * 4);
	b = CreateDIBSection(NULL, &bi, DIB_RGB_COLORS, (VOID **) ppvBits, NULL, 0);
	if (b == NULL)
		xpanic("error creating Image", GetLastError());
	return b;
}

HDC newDCForBitmap(HBITMAP bitmap, HBITMAP *prev)
{
	HDC screen, dc;

	screen = GetDC(NULL);
	if (screen == NULL)
		xpanic("error getting screen DC for NewImage()", GetLastError());
	dc = CreateCompatibleDC(screen);
	if (dc == NULL)
		xpanic("error creating memory DC for NewImage()", GetLastError());
	*prev = (HBITMAP) SelectObject(dc, bitmap);
	if (*prev == NULL)
		xpanic("error selecting bitmap into memory DC for NewImage()", GetLastError());
	if (ReleaseDC(NULL, screen) == 0)
		xpanic("error releasing screen DC for NewImage()", GetLastError());
	return dc;
}

void imageClose(HBITMAP bitmap, HDC dc, HBITMAP prev)
{
	if (SelectObject(dc, prev) != bitmap)
		xpanic("error restoring initial DC bitmap in Image.Close()", GetLastError());
	if (DeleteDC(dc) == 0)
		xpanic("error removing image DC in Image.Close()", GetLastError());
	if (DeleteObject(bitmap) == 0)
		xpanic("error removing bitmap in Image.Close()", GetLastError());
}

void moveTo(HDC dc, int x, int y)
{
	if (MoveToEx(dc, x, y, NULL) == 0)
		xpanic("error moving to point", GetLastError());
}

void lineTo(HDC dc, int x, int y)
{
	if (LineTo(dc, x, y) == 0)
		xpanic("error drawing line to point", GetLastError());
}

#define drawTextStyle (DT_LEFT | DT_TOP | DT_NOPREFIX | DT_SINGLELINE)

void drawText(HDC dc, char *str, int x, int y)
{
	WCHAR *wstr;
	RECT r;

	wstr = towstr(str);
	r.left = (LONG) x;
	r.top = (LONG) y;
	r.right = r.left;
	r.bottom = r.top;
	if (DrawTextW(dc, wstr, -1, &r, DT_CALCRECT | drawTextStyle) == 0)
		xpanic("error computing text bounding box", GetLastError());
	if (DrawTextW(dc, wstr, -1, &r, drawTextStyle) == 0)
		xpanic("error drawing text", GetLastError());
	freewstr(str);
}
