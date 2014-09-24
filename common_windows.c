// 18 september 2014

#include "winapi_windows.h"

HDC screenDC;

void init(void)
{
	screenDC = GetDC(NULL);
	if (screenDC == NULL)
		xpanic("error getting screen DC", GetLastError());
}

char *tostr(WCHAR *wstr)
{
	char *buf;
	int n;

	// alas WC_ERR_INVALID_CHARS is Vista-only
	// and WC_NO_BEST_FIT_CHARS is unsupported with CP_UTF8
	n = WideCharToMultiByte(CP_UTF8, 0, wstr, -1, NULL, 0, NULL, NULL);
	if (n == 0)
		xpanic("error getting buffer size in tostr()", GetLastError());
	// n includes the null terminator
	buf = (char *) malloc(n * sizeof (char));
	if (buf == NULL)		// TODO errno
		xpanic("error allocating buffer in tostr()", GetLastError());
	if (WideCharToMultiByte(CP_UTF8, 0, wstr, -1, buf, n, NULL, NULL) == 0)
		xpanic("error converting wide string to string in tostr()", GetLastError());
	return buf;
}

WCHAR *towstr(char *str)
{
	WCHAR *buf;
	int n;

	n = MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, str, -1, NULL, 0);
	if (n == 0)
		xpanic("error getting buffer size in towstr()", GetLastError());
	// n includes the null terminator
	buf = (WCHAR *) malloc(n * sizeof (WCHAR));
	if (buf == NULL)		// TODO errno
		xpanic("error allocating buffer in towstr()", GetLastError());
	if (MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, str, -1, buf, n) == 0)
		xpanic("error converting string to wide string in towstr()", GetLastError());
	return buf;
}

COLORREF colorref(uint8_t r, uint8_t g, uint8_t b)
{
	return RGB(r, g, b);
}
