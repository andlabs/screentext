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
#include <stdlib.h>

// /home/pietro/pkg/windows_386/github.com/andlabs/ui.a(_all.o): duplicate symbol reference: xpanic in both github.com/andlabs/ndraw(.text) and github.com/andlabs/ui(.text)
#define xpanic ndraw_xpanic

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
extern void penClose(HPEN);
extern HPEN penSelectInto(HPEN, HDC, COLORREF);
extern void penUnselect(HPEN, HDC, HPEN);

// fonts_windows.c
extern void listFonts(void *);
extern HFONT newFont(LOGFONTW *, char *, LONG);
extern void fontClose(HFONT);
extern HFONT fontSelectInto(HFONT, HDC);
extern void fontUnselect(HFONT, HDC, HFONT);

// image_windows.c
extern HBITMAP newBitmap(int, int, void **);
extern HDC newDCForBitmap(HBITMAP, HBITMAP *);
extern void imageClose(HBITMAP, HDC, HBITMAP);
extern void moveTo(HDC, int, int);
extern void lineTo(HDC, int, int);
extern void drawText(HDC, char *, int, int);

// common_windows.c
extern char *tostr(WCHAR *);
extern WCHAR *towstr(char *);
extern COLORREF colorref(uint8_t, uint8_t, uint8_t);

#endif
