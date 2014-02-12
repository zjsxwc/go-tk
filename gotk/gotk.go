package gotk

import ( "fmt"; "strings"; )//"unsafe" )

// #cgo CFLAGS: -I/usr/include/tcl8.6
// #cgo LDFLAGS: -ltcl8.6 -ltk8.6
// #include <tcl.h>
// #include <tk.h>
// #include <stdio.h>
// extern void CmdDispatch(unsigned int cb);
// static inline int CmdCallback(ClientData clientData, Tcl_Interp *interp, int objc, Tcl_Obj *const objv[]) {
//    CmdDispatch((unsigned int)clientData);
//    return 0;
// }
// static inline void RegisterCmd(Tcl_Interp *interp, char *cmdName, unsigned int cb) {
//    Tcl_CreateObjCommand( interp, cmdName, CmdCallback, (void *)cb, NULL );
// }
import "C"


var Interp *C.Tcl_Interp

func Init(useTk bool) bool {
	ok := false
	Interp = C.Tcl_CreateInterp()
	ok = C.Tcl_Init(Interp) == C.TCL_OK
	if useTk {
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

func ResultString() string {
	return C.GoString(C.Tcl_GetStringResult(Interp))
}


////////////////////////////////////////////////////////////////////////////////
type tk struct {
	wid uint //next widget ID
	widgets map[string]*Widget
	cmds []func()
}
var Tk = tk{ wid:0, widgets:make(map[string]*Widget) }

//get next widget ID
func (tk *tk) NewId() uint {
	tk.wid++
	return tk.wid-1
}
//get next widget ID as a string
func (tk *tk) NewIdS() string { return fmt.Sprintf(".gotk_%04d", tk.NewId()) }

//start Tk main loop
func (tk *tk) MainLoop() {
	C.Tk_MainLoop()
}
func (tk *tk) MainWindow() C.Tk_Window {
	return C.Tk_MainWindow(Interp)
}

//Tk options////////////////////////////////////////////////////////////////////
func (tk *tk) WindowingSystem() string {
	Eval("tk windowingsystem")
	return ResultString()
}
func (tk *tk) AppName() string {
	Eval("tk appname")
	return ResultString()
}
func (tk *tk) SetAppName(name string) {
	Eval("tk appname " + name)
}

//window options////////////////////////////////////////////////////////////////
func (tk *tk) Option(id string, name string, class string) string {
	Eval("option get " + id + " " + name + " " + class)
	return ResultString()
}
//widget options////////////////////////////////////////////////////////////////
func (tk *tk) Configure(id string, name string, val interface{}) {
	Eval( id + " configure -" + name + " " + fmt.Sprint(val) )
}

//add new widget to 
func (tk *tk) AddWidget(ws ...*Widget) {
	for _,w := range ws {
		tk.widgets[w.Id] = w
	}
}
func (tk *tk) ById(id string) *Widget {
	w,ok := tk.widgets[id]
	if ok { return w }
	return nil
}

////////////////////////////////////////////////////////////////////////////////
//export CmdDispatch
func CmdDispatch(cb C.uint) {
	Tk.cmds[cb]()
}
func (tk *tk) AddCmd(name string, cb func()) string {
	if name == "" { name = tk.NewIdS() }
	if tk.cmds == nil { tk.cmds = make([]func(), 0) }
	tk.cmds = append(tk.cmds, cb)
	C.RegisterCmd(Interp, C.CString(name), C.uint(len(tk.cmds)-1))
	return name
}


////////////////////////////////////////////////////////////////////////////////
//Message boxes/////////////////////////////////////////////////////////////////
//messagebox types
const (
	MBT_ABORTRETRYIGNORE = "abortretryignore"
	MBT_ARI              = "abortretryignore"
	MBT_OK               = "ok"
	MBT_OKCANCEL         = "okcancel"
	MBT_RETRYCANCEL      = "retrycancel"
	MBT_YESNO            = "yesno"
	MBT_YESNOCANCEL      = "yesnocancel"
)
//messagebox icon types
const (
	MBI_ERROR    = "error"
	MBI_INFO     = "info"
	MBI_QUESTION = "question"
	MBI_WARNING  = "warning"
)
//display message box and return result
func (tk *tk) MessageBox(typ string, title string, msg string, detailMsg string, icon string, defaultButton string, parent string) string {
	if (parent == "") { parent = "." }
	cmd := "tk_messageBox -parent " + parent
	if typ != "" { cmd += " -type " + typ }
	if title != "" { cmd += " -title " + title }
	if msg != "" { cmd += " -message " + msg }
	if detailMsg != "" { cmd += " -detail " + detailMsg }
	if icon != "" { cmd += " -icon " + icon }
	if defaultButton != "" { cmd += " -default " + defaultButton }
	Eval(cmd)
	return ResultString()
}


////////////////////////////////////////////////////////////////////////////////
type Widget struct {
	Type    string //widget type
	Id      string //this widget id
	Parent  string //parent window ID
	Options map[string]interface{}
}

func (tk *tk) New(typ string, id string, opts ...interface{}) *Widget {
	w := Widget{}
	w.Init(typ, id, opts...)
	w.Create()
	return &w
}

func (w *Widget) Reset() *Widget {
	w.Options = make(map[string]interface{})
	return w
}
func (w *Widget) Set(opt ...interface{}) *Widget {
	for pn:=0; pn+1<len(opt); pn += 2 {
		w.Options[ fmt.Sprint(opt[pn]) ] = opt[pn+1]
	}
	return w
}
func (w *Widget) SetMap(m map[string]interface{}) *Widget {
	for k,v := range m { w.Options[k] = v }
	return w
}

func (w *Widget) Init(typ string, id string, opt ...interface{}) *Widget {
	if typ != "" { w.Type = typ }
	if id != "" { w.Id = id
	} else if w.Id == "" { w.Id = Tk.NewIdS() }
	w.Reset()
	w.Set(opt...)
	return w
}

func (w *Widget) Create(opt ...interface{}) *Widget {
	w.Set(opt...)
	cmd := []string{w.Type, w.Id}
	for k,v := range w.Options {
		cmd = append(cmd, "-" + k, fmt.Sprint(v))
	}
	w.Reset()
	Eval(strings.Join(cmd, " "))
	Tk.AddWidget(w)
	return w
}

func (w *Widget) Pack(typ string, parent interface{}, opt ...interface{}) *Widget {
	par := "."
	switch p := parent.(type) {
		case *Window: par = p.Id
		case string:  if p != "" { par = p }
	}
	if typ == "" { typ = "pack" }
	cmd := []string{typ, w.Id, "-in", par}
	w.Set(opt...)
	for k,v := range w.Options {
		cmd = append( cmd, "-" + k, fmt.Sprint(v) )
	}
	w.Reset()
	Eval( strings.Join(cmd, " ") )
	return w
}

func (w *Widget) CGet(opt string) string {
	Eval(w.Id + " cget -" + opt)
	return ResultString()
}
func (w *Widget) CSet(opt ...interface{}) *Widget {
	cmd := []string{w.Id,"configure"}
	if (len(opt) > 1) {
		for n:=0; n+1<len(opt); n += 2 {
			cmd = append(cmd, "-" + fmt.Sprint(opt[n]), fmt.Sprint(opt[n+1]))
		}
		Eval( strings.Join(cmd, " ") )
	}
	return w
}

//events and callbacks//////////////////////////////////////////////////////////
func (w *Widget) Bind(event string, cb func()) *Widget {
	return w
}
func (w *Widget) Cmd(cmd string, cb func()) *Widget {
	w.CSet(cmd, Tk.AddCmd("", cb))
	return w
}

////////////////////////////////////////////////////////////////////////////////
type Window struct {
	Id     string //this widget id
}

func (w *Window) Init(id string) *Window {
	if id == "" { id = "." }
	w.Id = id
	
	return w
}
