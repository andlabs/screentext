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

// image_windows.c
struct image {
	HBITMAP bitmap;
	HDC dc;
	HBITMAP prev;
	VOID *ppvBits;
	int width;
	int height;
};
extern struct image *newImage(int, int, BOOL);
extern void imageClose(struct image *);
extern void line(struct image *, int, int, int, int, HPEN, uint8_t);
extern void strokeText(struct image *, char *, int, int, HFONT, HPEN, uint8_t);
extern void fillText(struct image *, char *, int, int, HFONT, HBRUSH, uint8_t);
extern SIZE textSize(struct image *, char *, HFONT);

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
extern HPEN penSelectInto(HPEN, HDC);
extern void penUnselect(HPEN, HDC, HPEN);

// fonts_windows.c
extern void listFonts(void *);
extern HFONT newFont(LOGFONTW *, char *, LONG);
extern void fontClose(HFONT);
extern HFONT fontSelectInto(HFONT, HDC);
extern void fontUnselect(HFONT, HDC, HFONT);

// common_windows.c
extern char *tostr(WCHAR *);
extern WCHAR *towstr(char *);
extern COLORREF colorref(uint8_t, uint8_t, uint8_t);

// brush_windows.c
extern HBRUSH newBrush(LOGBRUSH *);
extern void brushClose(HBRUSH);
extern HBRUSH brushSelectInto(HBRUSH, HDC);
extern void brushUnselect(HBRUSH, HDC, HBRUSH);

#endif
