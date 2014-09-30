// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

struct image *newImage(int dx, int dy, BOOL internal)
{
	struct image *i;
	BITMAPINFO bi;

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
	// sets the background to white
	// see toImage() in main_windows.go for details
	memset(i->ppvBits, 0xFF, dx * dy * 4);

	i->dc = CreateCompatibleDC(screenDC);
	if (i->dc == NULL)
		xpanic("error creating memory DC for newImage()", GetLastError());
	i->prev = (HBITMAP) SelectObject(i->dc, i->bitmap);
	if (i->prev == NULL)
		xpanic("error selecting bitmap into memory DC for newImage()", GetLastError());

	i->width = dx;
	i->height = dy;
	return i;
}

void imageClose(struct image *i)
{
	if (SelectObject(i->dc, i->prev) != i->bitmap)
		xpanic("error restoring initial DC bitmap in imageClose()", GetLastError());
	if (DeleteDC(i->dc) == 0)
		xpanic("error removing image DC in imageClose()", GetLastError());
	if (DeleteObject(i->bitmap) == 0)
		xpanic("error removing bitmap in imageClose()", GetLastError());
	free(i);
}

static SIZE wtextSize(WCHAR *wstr, HFONT font)
{
	SIZE size;
	HFONT prevFont;

	prevFont = fontSelectInto(font, screenDC);
	if (GetTextExtentPoint32W(screenDC, wstr, wcslen(wstr), &size) == 0)
		xpanic("error getting text size", GetLastError());
	fontUnselect(font, screenDC, prevFont);
	return size;
}

struct image *drawText(char *str, HFONT font)
{
	WCHAR *wstr;
	SIZE size;
	struct image *ti;
	HFONT prevFont;

	wstr = towstr(str);
	size = wtextSize(wstr, font);
	ti = newImage(size.cx, size.cy, TRUE);
	prevFont = fontSelectInto(font, ti->dc);
	// see toImage() in main_windows.go
	if (SetTextColor(ti->dc, RGB(0, 0, 0)) == CLR_INVALID)
		xpanic("error giving text black color for alpha calculations", GetLastError());
	if (SetBkMode(ti->dc, TRANSPARENT) == 0)
		xpanic("error setting text drawing to have nonopaque background", GetLastError());
	if (TextOutW(ti->dc, 0, 0, wstr, wcslen(wstr)) == 0)
		xpanic("error drawing text path", GetLastError());
	fontUnselect(font, ti->dc, prevFont);
	free(wstr);
	if (GdiFlush() == 0)
		xpanic("error committing GDI operations to memory", GetLastError());
	return ti;
}

SIZE textSize(char *str, HFONT font)
{
	WCHAR *wstr;
	SIZE size;

	wstr = towstr(str);
	size = wtextSize(wstr, font);
	free(wstr);
	return size;
}
