package w32

import (
	"syscall"
	"tbase/core/w32/dll"
	"unsafe"
)

var (
	user32 = syscall.MustLoadDLL("user32.dll")

	procInvalidateRect = user32.MustFindProc("InvalidateRect")

	procIsHungAppWindow = user32.MustFindProc("IsHungAppWindow")
	procGetMessage      = user32.MustFindProc("GetMessageW")
	procDispatchMessage = user32.MustFindProc("DispatchMessageW")
	procGetWindowRect   = user32.MustFindProc("GetWindowRect")

	procDefWindowProc = user32.MustFindProc("DefWindowProcW")

	procBeginPaint = user32.MustFindProc("BeginPaint")
	procEndPaint   = user32.MustFindProc("EndPaint")
)

func (w *Win32) _CreateIcon(icon []byte, iconWidth, iconHeight int) {
	w.icon = dll.ProcPanic(user32,
		"CreateIcon",
		uintptr(w.instance),
		uintptr(iconWidth),
		uintptr(iconHeight),
		1, 32,
		uintptr(unsafe.Pointer(&icon[0])),
		uintptr(unsafe.Pointer(&icon[0])),
	)
}

func (w *Win32) _RegisterClass(wcx *_WNDCLASSX) {
	wcx.CbSize = uint32(unsafe.Sizeof(*wcx))

	proc := user32.MustFindProc("RegisterClassExW")
	r0, _, e1 := proc.Call(uintptr(unsafe.Pointer(wcx)))
	var err error
	classAtom := uint16(r0)
	if classAtom == 0 {
		if e1 != nil {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}

	if err != nil {
		panic(err)
	}
	w.classAtom = classAtom
}

func (w *Win32) _CreateWindow(style int) {
	w.window = dll.ProcPanic(user32,
		"CreateWindowExW",
		uintptr(WS_EX_APPWINDOW), // opt.Ex_Style
		uintptr(w.classAtom),
		uintptr(unsafe.Pointer(dll.UTF16PtrFromString(title))), // TitleWindows: lpWindowName
		uintptr(style),

		uintptr(w.screen.Left),
		uintptr(w.screen.Top),
		uintptr(w.screen.Width()),
		uintptr(w.screen.Height()),

		uintptr(0), // HWND_DESKTOP => no owner window: hWndParent
		0,          // use class menu: hMenu
		uintptr(w.instance),
		0, // no window-creation data: lpParam
	)
}

func (w *Win32) _ControlClass(hndl syscall.Handle, msg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	w.wg.Add(1)
	defer w.wg.Done()
	switch msg {
	case WM_CREATE: // create Buffer and init Windows
		w.window = hndl
		// w._ShowCursor(false)
		w.buf = make([]byte, w.screen.Width()*w.screen.Height()*4)
		// w.buf, _, _ = util.ImgtoByte(util.GetImage("assets/bk2.png"))
	case WM_PAINT: // Draw
		ps := _PAINTSTRUCT{}
		defer w._EndPaint(&ps)
		w._BeginPaint(&ps)

		rw, rh := ps.Rect.Width(), ps.Rect.Height()
		if rw == 0 || rh == 0 {
			return 0
		}

		memDC := CreateCompatibleDC(ps.Hdc)
		defer DeleteDC(memDC)

		btm := CreateBitmap(ps.Hdc, rw, rh)

		b := w._Buffer(rw, rh, int(ps.Rect.Left), int(ps.Rect.Top))
		SetBitmapBits(uintptr(btm), uintptr(rw*rh*4), uintptr(unsafe.Pointer(&b[0])))

		SelectObject(memDC, btm)
		BitBlt(ps.Hdc,
			int(ps.Rect.Left), int(ps.Rect.Top),
			rw, rh,
			memDC,
			0, 0,
			SRCCOPY,
		)
		DeleteObject(btm)

	case WM_GETMINMAXINFO: // Set MinSizes of cons windows
		info := (*_MINMAXINFO)(unsafe.Pointer(lParam))
		info.MinTrackSize.X = minSizeX
		info.MinTrackSize.Y = minSizeY
		return 0
	case WM_MOVE, WM_SIZE: // Update position and size of the window
		w._UpdateSize()
		w.buf = make([]byte, w.screen.Width()*w.screen.Height()*4)
	case WM_QUIT, WM_DESTROY, WM_NCDESTROY, WM_CLOSE: // Close
		*w.isRun = false
		return w.Quit()
	default: // DefaultWindowProctocol
		return w._DefWindowProc(hndl, msg, wParam, lParam)
	}
	return 0
}

// w32.InvalidateRect invalidate a rectangular area of a window and cause a WM_PAINT message to be sent
func (w *Win32) _DrawScreen(rect *_RECT) {
	dll.ProcLazyNoZero(procInvalidateRect, uintptr(w.window), uintptr(unsafe.Pointer(rect)))
}

func (w *Win32) _GetMessage(msg *_MSG) bool {
	r0, _ := dll.ProcLazy(procGetMessage, uintptr(unsafe.Pointer(msg)))
	return int32(r0) > 0
}

func (w *Win32) _DispatchMessage(msg *_MSG) int32 {
	r0, _ := dll.ProcLazy(procDispatchMessage, uintptr(unsafe.Pointer(msg)))
	return int32(r0)
}

func (w *Win32) _DefWindowProc(hndl syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) uintptr {
	hndl, _ = dll.ProcLazy(procDefWindowProc, uintptr(hndl), uintptr(uMsg), uintptr(wParam), uintptr(lParam))
	return uintptr(hndl)
}

// GetWindowRect ( Rect -> w.screen )
func (w *Win32) _UpdateSize() {
	dll.ProcLazyPanic(procGetWindowRect, uintptr(w.window), uintptr(unsafe.Pointer(w.screen)))
}

func (w *Win32) _BeginPaint(ps *_PAINTSTRUCT) syscall.Handle {
	return dll.ProcLazyPanic(procBeginPaint, uintptr(w.window), uintptr(unsafe.Pointer(ps)))
}

func (w *Win32) _EndPaint(ps *_PAINTSTRUCT) syscall.Handle {
	return dll.ProcLazyPanic(procEndPaint, uintptr(w.window), uintptr(unsafe.Pointer(ps)))
}
