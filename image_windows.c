// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

struct image *newImage(int dx, int dy, BOOL internal)
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
	if (internal)
		// see imageInternalBlend() below for details
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

	i->width = dx;
	i->height = dy;
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

// GDI doesn't natively support alpha blending outside of AlphaBlend(), but there's a trick: GDI sets the alpha byte of any written pixel to 0x00
// this means we can manually patch in alpha values
// this is important since we must also premultiply
// so what we do is: all drawing operations actually draw to an internal bitmap that's set to all 0xFF (see internal in newImage() above) and then we patch in the alpha ourselves here before calling AlphaBlend()
// in the end, the only function that should ever be called to a non-internal image is AlphaBlend()
static void imageInternalBlend(struct image *dest, struct image *src, uint8_t alpha)
{
	uint8_t *sp;
	int x, y;
	BLENDFUNCTION bf;

	// first make sure GDI has written everything
	// TODO this is per-thread...
	if (GdiFlush() == 0)
		xpanic("error flushing GDI buffers", GetLastError());
	sp = (uint8_t *) src->ppvBits;
	for (y = 0; y < src->height; y++)
		for (x = 0; x < src->width; x++)
			if (*(sp + 3) == 0xFF) {		// not written by GDI; make transparent
				*sp++ = 0;
				*sp++ = 0;
				*sp++ = 0;
				*sp++ = 0;
			} else {					// written by GDI; premultiply and set alpha
				// premultiplication steps from http://msdn.microsoft.com/en-us/library/dd183393%28v=vs.85%29.aspx
				*sp = (*sp * alpha) / 0xFF;		// R
				sp++;
				*sp = (*sp * alpha) / 0xFF;		// G
				sp++;
				*sp = (*sp * alpha) / 0xFF;		// B
				sp++;
				*sp = alpha;					// A
				sp++;
			}
	// and now blend
	ZeroMemory(&bf, sizeof (BLENDFUNCTION));
	bf.BlendOp = AC_SRC_OVER;
	bf.BlendFlags = 0;
	bf.SourceConstantAlpha = 255;		// per-pixel alpha
	bf.AlphaFormat = AC_SRC_ALPHA;
	if (AlphaBlend(dest->dc, 0, 0, dest->width, dest->height,
		src->dc, 0, 0, src->width, src->height,
		bf) == FALSE)
		xpanic("error doing internal alpha-blending of image draw", GetLastError());
}

void line(struct image *i, int x0, int y0, int x1, int y1, HPEN pen, uint8_t alpha)
{
	struct image *li;
	HPEN prev;

	li = newImage(i->width, i->height, TRUE);
	prev = penSelectInto(pen, li->dc);
	if (MoveToEx(li->dc, x0, y0, NULL) == 0)
		xpanic("error moving to point", GetLastError());
	if (LineTo(li->dc, x1, y1) == 0)
		xpanic("error drawing line to point", GetLastError());
	imageInternalBlend(i, li, alpha);
	penUnselect(pen, li->dc, prev);
	imageClose(li);
}

void strokeText(struct image *i, char *str, int x, int y, HFONT font, HPEN pen, uint8_t alpha)
{
	WCHAR *wstr;
	struct image *ti;
	HPEN prevPen;
	HFONT prevFont;

	wstr = towstr(str);
	ti = newImage(i->width, i->height, TRUE);
	prevFont = fontSelectInto(font, ti->dc);
	prevPen = penSelectInto(pen, ti->dc);
	if (BeginPath(ti->dc) == 0)
		xpanic("error beginning text drawing path", GetLastError());
	if (SetBkMode(ti->dc, TRANSPARENT) == 0)
		xpanic("error setting text drawing to have nonopaque background", GetLastError());
	if (TextOutW(ti->dc, x, y, wstr, wcslen(wstr)) == 0)
		xpanic("error drawing text path", GetLastError());
	if (EndPath(ti->dc) == 0)
		xpanic("error ending text drawing path", GetLastError());
	if (StrokePath(ti->dc) == 0)
		xpanic("error stroking text drawing path", GetLastError());
	imageInternalBlend(i, ti, alpha);
	penUnselect(pen, ti->dc, prevPen);
	fontUnselect(font, ti->dc, prevFont);
	imageClose(ti);
	free(wstr);
}

// TODO merge with strokeText()
void fillText(struct image *i, char *str, int x, int y, HFONT font, HBRUSH brush, uint8_t alpha)
{
	WCHAR *wstr;
	struct image *ti;
	HPEN prevBrush;
	HFONT prevFont;

	wstr = towstr(str);
	ti = newImage(i->width, i->height, TRUE);
	prevFont = fontSelectInto(font, ti->dc);
	prevBrush = brushSelectInto(brush, ti->dc);
	if (BeginPath(ti->dc) == 0)
		xpanic("error beginning text drawing path", GetLastError());
	if (SetBkMode(ti->dc, TRANSPARENT) == 0)
		xpanic("error setting text drawing to have nonopaque background", GetLastError());
	if (TextOutW(ti->dc, x, y, wstr, wcslen(wstr)) == 0)
		xpanic("error drawing text path", GetLastError());
	if (EndPath(ti->dc) == 0)
		xpanic("error ending text drawing path", GetLastError());
	if (FillPath(ti->dc) == 0)
		xpanic("error filling text drawing path", GetLastError());
	imageInternalBlend(i, ti, alpha);
	brushUnselect(brush, ti->dc, prevBrush);
	fontUnselect(font, ti->dc, prevFont);
	imageClose(ti);
	free(wstr);
}

SIZE textSize(struct image *i, char *str, HFONT font)
{
	WCHAR *wstr;
	SIZE size;
	HFONT prevFont;

	wstr = towstr(str);
	prevFont = fontSelectInto(font, i->dc);
	if (GetTextExtentPoint32W(i->dc, wstr, wcslen(wstr), &size) == 0)
		xpanic("error getting text size", GetLastError());
	fontUnselect(font, i->dc, prevFont);
	free(wstr);
	return size;
}
