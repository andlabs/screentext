// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

struct image *newImage(int dx, int dy)
{
	struct image *i;
	BITMAPINFO bi;
	HDC screen;

	i = (struct image *) malloc(sizeof (struct image));
	if (i == NULL)		// TODO errno
		xpanic("memory exhausted allocating image data in newImage()", GetLastError());
	ZeroMemory(i, sizeof (struct image));

	ZeroMemory(&bi, sizeof (BITMAPINFO));
	bi.bmiHeader.biSize = sizeof (BITMAPINFOHEADER);
	bi.bmiHeader.biWidth = (LONG) dx;
	bi.bmiHeader.biHeight = -((LONG) dy);                   // negative height to force top-down drawing
	bi.bmiHeader.biPlanes = 1;
	bi.bmiHeader.biBitCount = 32;
	bi.bmiHeader.biCompression = BI_RGB;
	bi.bmiHeader.biSizeImage = (DWORD) (dx * dy * 4);
	i->bitmap = CreateDIBSection(NULL, &bi, DIB_RGB_COLORS, &i->ppvBits, NULL, 0);
	if (i->bitmap == NULL)
		xpanic("error creating image in newImage()", GetLastError());
	// see image.Image() in image_windows.go for details
	memset(i->ppvBits, 0xFF, dx * dy * 4);

	screen = GetDC(NULL);
	if (screen == NULL)
		xpanic("error getting screen DC for newImage()", GetLastError());
	i->dc = CreateCompatibleDC(screen);
	if (i->dc == NULL)
		xpanic("error creating memory DC for newImage()", GetLastError());
	i->prev = (HBITMAP) SelectObject(i->dc, i->bitmap);
	if (i->prev == NULL)
		xpanic("error selecting bitmap into memory DC for newImage()", GetLastError());
	if (ReleaseDC(NULL, screen) == 0)
		xpanic("error releasing screen DC for newImage()", GetLastError());

	return i;
}

void imageClose(struct image *i)
{
	if (SelectObject(i->dc, i->prev) != i->bitmap)
		xpanic("error restoring initial DC bitmap in Image.Close()", GetLastError());
	if (DeleteDC(i->dc) == 0)
		xpanic("error removing image DC in Image.Close()", GetLastError());
	if (DeleteObject(i->bitmap) == 0)
		xpanic("error removing bitmap in Image.Close()", GetLastError());
	free(i);
}

void line(struct image *i, int x0, int y0, int x1, int y1)
{
	if (MoveToEx(i->dc, x0, y0, NULL) == 0)
		xpanic("error moving to point", GetLastError());
	if (LineTo(i->dc, x1, y1) == 0)
		xpanic("error drawing line to point", GetLastError());
}

void drawText(struct image *i, char *str, int x, int y)
{
	WCHAR *wstr;

	wstr = towstr(str);
	if (SetBkMode(i->dc, TRANSPARENT) == 0)
		xpanic("error setting text drawing to be transparent", GetLastError());
	if (TextOutW(i->dc, x, y, wstr, wcslen(wstr)) == 0)
		xpanic("error drawing text", GetLastError());
	free(str);
}
