This is a library for Go that draws text onto an image.Image. The text uses the underlying OS's text rendering system, so system fonts can be used. It is designed for rendering to the screen, so text is antialiased, etc.

It started as a more general-purpose vector graphics library, but technical restirctions and a general misunderstanding of the problem of vector graphics and device specificity made this unreasonable, so the scope was limited to just text strings. I plan on making a better vector graphics subsystem later.

It's sloppy :/

REQUIREMENTS
* Windows: Windows XP or newer (same as Go); uses GDI
* Unix (not Mac OS X): cairo (>=1.10) and pango (>=1.30)
* Mac OS X: Mac OS X 10.6 or newer; uses Core Graphics and Core Text directly
