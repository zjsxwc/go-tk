package gotk

import ( "fmt" )

// #cgo CFLAGS: -I/usr/include/tcl8.6
// #cgo LDFLAGS: -ltcl8.6 -ltk8.6
// #include <tcl.h>
// #include <tk.h>
import "C"

var Interp *C.Tcl_Interp
var tk bool

func Init(useTk bool) bool {
	ok := false
	Interp = C.Tcl_CreateInterp()
	ok = C.Tcl_Init(Interp) == C.TCL_OK
	if useTk {
		tk = true
		ok = C.Tk_Init(Interp) == C.TCL_OK
	}
	return ok
}

func Exit() {
}

func SetVar( name string, val interface{} ) {
	C.Tcl_SetVar( Interp, C.CString(name), C.CString(fmt.Sprint(val)), 0 )//C.TCL_GLOBAL_ONLY )
}
func GetVar( name string ) string {
	return C.GoString(C.Tcl_GetVar( Interp, C.CString(name), 0))//C.TCL_GLOBAL_ONLY ))
}
func UnsetVar( name string ) {
	C.Tcl_UnsetVar( Interp, C.CString(name), 0)//C.TCL_GLOBAL_ONLY )
}

func Eval( script string ) bool {
	return C.Tcl_Eval( Interp, C.CString(script) ) == C.TCL_OK
}

func TkMainLoop() {
	C.Tk_MainLoop()
}

func ResultString() string {
	return C.GoString(C.Tcl_GetStringResult(Interp))
}
