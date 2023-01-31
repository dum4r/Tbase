package dll

import "syscall"

func UTF16PtrFromString(str string) *uint16 {
	pointer, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		panic(err)
	}
	return pointer
}

func UintptrToBool(value uintptr) bool { return value == 1 }
func BoolToUintptr(bln bool) uintptr {
	if bln {
		return 1
	}
	return 0
}

// template protocols
func Proc(dll *syscall.DLL, name string, a ...uintptr) (hndl syscall.Handle, err error) {
	proc := dll.MustFindProc(name)
	r0, _, e1 := proc.Call(a...)
	hndl = syscall.Handle(r0)
	if hndl == 0 && e1 != nil {
		err = e1
	}
	defer func() {
		proc = nil
	}()
	return
}

func ProcPanic(dll *syscall.DLL, name string, a ...uintptr) syscall.Handle {
	hndl, err := Proc(dll, name, a...)
	if err != nil {
		println("Error ProcPanic DLL => '", name, "'")
		panic(err)
	}
	return hndl
}

func ProcNoZero(dll *syscall.DLL, name string, a ...uintptr) {
	r0, err := Proc(dll, name, a...)
	if r0 == 0 {
		panic("Error ProcNoZero :'" + name + "' => " + err.Error())
	}
}

func ProcZero(dll *syscall.DLL, name string, a ...uintptr) {
	r0, err := Proc(dll, name, a...)
	if r0 != 0 {
		panic("Error ProcZero :'" + name + "' => " + err.Error())
	}
}

func ProcLazy(proc *syscall.Proc, a ...uintptr) (hndl syscall.Handle, err error) {
	r0, _, e1 := proc.Call(a...)
	hndl = syscall.Handle(r0)
	if hndl == 0 && e1 != nil {
		err = e1
	}
	return
}

func ProcLazyPanic(proc *syscall.Proc, a ...uintptr) syscall.Handle {
	hndl, err := ProcLazy(proc, a...)
	if err != nil {
		// panic(err)
		println(hndl, " => ", err.Error())
	}
	return hndl
}

func ProcLazyNoZero(proc *syscall.Proc, a ...uintptr) {
	r0, err := ProcLazy(proc, a...)
	if r0 == 0 {
		panic("Error ProcLazyNoZero => " + err.Error())
		// println(err.Error())
	}
}
