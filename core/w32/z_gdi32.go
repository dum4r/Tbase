package w32

import (
	"syscall"
	"tbase/core/w32/dll"
)

var (
	gdi32 = syscall.MustLoadDLL("gdi32.dll")

	procCreateCompatibleDC     = gdi32.MustFindProc("CreateCompatibleDC")
	procDeleteDC               = gdi32.MustFindProc("DeleteDC")
	procCreateCompatibleBitmap = gdi32.MustFindProc("CreateCompatibleBitmap")
	procSetBitmapBits          = gdi32.MustFindProc("SetBitmapBits")
	procDeleteObject           = gdi32.MustFindProc("DeleteObject")
	procSelectObject           = gdi32.MustFindProc("SelectObject")

	procBitBlt   = gdi32.MustFindProc("BitBlt")
	procSetPixel = gdi32.MustFindProc("SetPixelV")
)

func CreateCompatibleDC(hDC syscall.Handle) syscall.Handle {
	return dll.ProcLazyPanic(procCreateCompatibleDC, uintptr(hDC))
}
func DeleteDC(hDC syscall.Handle) { dll.ProcLazyNoZero(procDeleteDC, uintptr(hDC)) }

func DeleteObject(hndl syscall.Handle) { dll.ProcLazyNoZero(procDeleteObject, uintptr(hndl)) }

func SelectObject(hDC syscall.Handle, hndl syscall.Handle) syscall.Handle {
	return dll.ProcLazyPanic(procSelectObject, uintptr(hDC), uintptr(hndl))
}

func BitBlt(hDC syscall.Handle, x, y, w1, h1 int, hsrc syscall.Handle, w2, h2 int, dword uintptr) syscall.Handle {
	return dll.ProcLazyPanic(procBitBlt,
		uintptr(hDC),
		uintptr(x),
		uintptr(y),
		uintptr(w1),
		uintptr(h1),
		uintptr(hsrc),
		uintptr(w2),
		uintptr(h2),
		dword,
	)
}

func SetPixel(hDC syscall.Handle, x, y int, color uintptr) {
	dll.ProcLazyNoZero(procSetPixel, uintptr(hDC), uintptr(x), uintptr(y), color)
}

// CreateBitmap
func CreateBitmap(hDC syscall.Handle, x, y int) syscall.Handle {
	return dll.ProcLazyPanic(procCreateCompatibleBitmap, uintptr(hDC), uintptr(x), uintptr(y))
}

func SetBitmapBits(btm, t, ter uintptr) syscall.Handle {
	return dll.ProcLazyPanic(procSetBitmapBits, btm, t, ter)
}
