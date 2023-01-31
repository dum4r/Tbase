package w32

import (
	"syscall"
)

type _WNDCLASSX struct {
	CbSize        uint32
	Style         int32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

// the _RECT is RECT of user32.dll
type _RECT struct {
	Left, Top, Right, Bottom int32
}

func (r *_RECT) Width() int  { return int(r.Right - r.Left) }
func (r *_RECT) Height() int { return int(r.Bottom - r.Top) }

// add the Default function: that matches the values to CW_USEDEFAULT
// add the Fullscreen function: exec GetSystemMetrics and assign new values
type Screen _RECT

func (r *Screen) Rect() *_RECT {
	return (*_RECT)(r)
}

// Separar Default  y FullScreen de Rect ????????
func (r *Screen) Default() {
	r.Left, r.Top, r.Right, r.Bottom = CW_USEDEFAULT, CW_USEDEFAULT, 0, 0
}
func (r *Screen) Fullscreen() {
	proc := user32.MustFindProc("GetSystemMetrics")
	w, _, _ := proc.Call(0) // SM_CXSCREEN = 0
	h, _, _ := proc.Call(1) // SM_CYSCREEN = 1
	r.Left, r.Top, r.Right, r.Bottom = 0, 0, int32(w), int32(h)
}

type _MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 } // POINT user32.dll
}

type _PAINTSTRUCT struct {
	Hdc   syscall.Handle
	Erase bool
	Rect  _RECT
	// Used internally by the system.
	restore     bool
	incUpdate   bool
	rgbReserved [32]byte
}

type POINT struct {
	X int32
	Y int32
}
type _MINMAXINFO struct {
	Reserved     POINT
	MaxSize      POINT
	MaxPosition  POINT
	MinTrackSize POINT
	MaxTrackSize POINT
}

// {0x28, 0x28, 0x2e, 0x0},
// {0x6c, 0x56, 0x71, 0x0},
// {0xd9, 0xc8, 0xbf, 0x0},
// {0xf9, 0x82, 0x84, 0x0},
// {0xb0, 0xa9, 0xe4, 0x0},
// {0xac, 0xcc, 0xe4, 0x0},
// {0xb3, 0xe3, 0xda, 0x0},
// {0xfe, 0xaa, 0xe4, 0x0},
// {0x87, 0xa8, 0x89, 0x0},
// {0xb0, 0xeb, 0x93, 0x0},
// {0xe9, 0xf5, 0x9d, 0x0},
// {0xff, 0xe6, 0xc6, 0x0},
// {0xde, 0xa3, 0x8b, 0x0},
// {0xff, 0xc3, 0x84, 0x0},
// {0xff, 0xf7, 0xa0, 0x0},
// {0xff, 0xf7, 0xe4, 0x0},
