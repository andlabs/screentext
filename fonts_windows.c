// 18 september 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

struct lfd {
	void *golist;
	HDC dc;
};

static int CALLBACK enumFontProc(const LOGFONTW *lf, const TEXTMETRICW *tm, DWORD type, LPARAM lParam)
{
	struct lfd *lfd = (struct lfd *) lParam;
	LONG h;

	// convert lfHeight to point size
	h = lf->lfHeight;
	if (h < 0) {
		// solve height = -MulDiv(points, base, 72dpi) for points to get this
		h = -h;
		h = MulDiv(h, 72, GetDeviceCaps(lfd->dc, LOGPIXELSY));
	} else if (h > 0) {
		// solve the first formula in http://support.microsoft.com/kb/74299 for point size
		h = MulDiv(h - tm->tmInternalLeading, 72, GetDeviceCaps(lfd->dc, LOGPIXELSY));
	} else
		xpanic("don't know how to handle lfHeight == 0 for ListFonts()", 0);

	listFontsAdd(lfd->golist, lf, tostr(lf->lfFaceName), h);
	return 1;
}

void listFonts(void *golist)
{
	struct lfd lfd;
	LOGFONTW spec;

	lfd.golist = golist;
	lfd.dc = GetDC(NULL);
	if (lfd.dc == NULL)
		xpanic("error getting screen DC for ListFonts()", GetLastError());
	ZeroMemory(&spec, sizeof (LOGFONTW));
	spec.lfCharSet = DEFAULT_CHARSET;		// all character sets
	spec.lfFaceName[0] = L'\0';				// all faces
	spec.lfPitchAndFamily = 0;
	EnumFontFamiliesExW(lfd.dc, &spec, enumFontProc, (LPARAM) (&lfd), 0);
	if (ReleaseDC(NULL, lfd.dc) == 0)
		xpanic("error releasing screen DC for ListFonts()", GetLastError());
}

HFONT newFont(LOGFONTW *lf, char *family, LONG size)
{
	HDC dc;
	HFONT f;

	dc = GetDC(NULL);
	if (dc == NULL)
		xpanic("error getting screen DC for NewFont() (needed for size calculation)", GetLastError());
	if (MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS,
		family, -1,
		lf->lfFaceName, LF_FACESIZE) == 0)
		xpanic("error loading FontSpec.Family into LOGFONTW for NewFont()", GetLastError());
	// see http://msdn.microsoft.com/en-us/library/windows/desktop/dd145037%28v=vs.85%29.aspx
	lf->lfHeight = -MulDiv(size, GetDeviceCaps(dc, LOGPIXELSY), 72);
	f = CreateFontIndirectW(lf);
	if (f == NULL)
		xpanic("error creating font in NewFont()", GetLastError());
	if (ReleaseDC(NULL, dc) == 0)
		xpanic("error releasing screen DC for NewFont() (needed for size calculation)", GetLastError());
	return f;
}

void fontClose(HFONT f)
{
	if (DeleteObject(f) == 0)
		xpanic("error closing Font", GetLastError());
}

HFONT fontSelectInto(HFONT font, HDC dc)
{
	HFONT prev;

	prev = (HFONT) SelectObject(dc, font);
	if (prev == NULL)
		xpanic("error selecting Font into Image DC", GetLastError());
	return prev;
}

void fontUnselect(HFONT font, HDC dc, HFONT prev)
{
	if (SelectObject(dc, prev) != font)
		xpanic("error unselecting Font from Image DC", GetLastError());
}
