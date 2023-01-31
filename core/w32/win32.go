package w32

import (
	"runtime"
	"sync"
	"syscall"
	"tbase/core/util"
	"tbase/core/w32/dll"
)

const (
	clasName           string = "CORE_TBASE"
	title              string = "TyzyBase"
	urlIcon            string = "assets/icon.png"
	minSizeX, minSizeY int32  = 900, 600
)

type Win32 struct {
	isRun  *bool
	screen *_RECT

	instance  syscall.Handle
	window    syscall.Handle
	icon      syscall.Handle
	classAtom uint16

	msg         _MSG
	chanUpdates chan _RECT
	buf         []byte
	wg          sync.WaitGroup
}

func (w *Win32) Create(isRun *bool, style int, screen *_RECT) error {
	if screen == nil || isRun == nil {
		return &syscall.DLLError{Msg: "Error: Create Window 32 can't be nil"}
	}
	w.screen = screen
	w.isRun = isRun
	w.chanUpdates = make(chan _RECT, 2)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	icon, iconWidth, iconHeight := util.IcontoByte(util.GetImage(urlIcon))

	w._CreateIcon(icon, iconWidth, iconHeight)

	// RegisterClassExW
	w._RegisterClass(
		&_WNDCLASSX{
			Style:         CS_CLASSDC | CS_DBLCLKS | CS_SAVEBITS,
			LpfnWndProc:   syscall.NewCallback(w._ControlClass),
			HInstance:     syscall.Handle(w.instance),
			HIcon:         w.icon,
			HbrBackground: syscall.Handle(1),
			LpszMenuName:  new(uint16),
			LpszClassName: dll.UTF16PtrFromString(clasName),
			HIconSm:       w.icon,
		},
	)

	// CreateWindowExW
	w._CreateWindow(style)

	go func() {
		for rect := range w.chanUpdates {
			w._DrawScreen(&rect)
		}
	}()
	return nil
}

func (w *Win32) Destroy() { dll.ProcNoZero(user32, "DestroyWindow", uintptr(w.window)) }

func (w *Win32) Manager() {
	// msg := w.msg
	// S.O. Retrieves a message from the calling thread's message queue.
	if w._GetMessage(&w.msg) {
		// w32.TranslateMessage(msg)
		w._DispatchMessage(&w.msg)
	}
}

func (w *Win32) IsHung() bool {
	r0, err := dll.ProcLazy(procIsHungAppWindow, uintptr(w.window))
	if r0 > 1 && err != nil {
		panic(err.Error())
	}
	return dll.UintptrToBool(uintptr(r0))
}

func (w *Win32) Wait() { w.wg.Wait() }

func (w *Win32) Quit() uintptr {
	return uintptr(dll.ProcPanic(user32, "PostQuitMessage", uintptr(w.window), WM_QUIT, 0, 0))
}
