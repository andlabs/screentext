// 18 september 2014

// cgo will include this file multiple times
#ifndef __GO_NDRAW_WINAPI_WINDOWS_H__
#define __GO_NDRAW_WINAPI_WINDOWS_H__

#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
// get Windows version right; right now Windows XP
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <stdint.h>

// pen_windows.c
// the following struct is needed because there is no ExtCreatePenIndirect() :(
struct xpen {
	DWORD style;
	DWORD width;
	LOGBRUSH brush;
	DWORD nSegments;
	DWORD *segments;
};
extern HPEN newPen(struct xpen *);

// fonts_windows.c
extern void listFonts(void *);
extern HFONT newFont(LOGFONTW *, char *, LONG);

#endif
