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
extern struct image *drawText(char *, HFONT, uint8_t, uint8_t, uint8_t);
extern SIZE textSize(char *, HFONT);

// fonts_windows.c
extern void listFonts(void *);
extern HFONT newFont(LOGFONTW *, char *, LONG);
extern void fontClose(HFONT);
extern HFONT fontSelectInto(HFONT, HDC);
extern void fontUnselect(HFONT, HDC, HFONT);

// common_windows.c
extern HDC screenDC;
extern void init(void);
extern char *tostr(WCHAR *);
extern WCHAR *towstr(char *);
extern COLORREF colorref(uint8_t, uint8_t, uint8_t);

#endif
