gotk
====

Go-Tk GUI binding

Example:

``` go
package main

import ( "gotk"; "os"; "fmt";  )

func main() {
	//init Tcl and Tk
	if !gotk.Init(true) { os.Exit(0) }

	w := gotk.Widget{} //create new widget
	w.Init("label", "", //configure its type and parameters, and pack
		"text",123,
		"fg","red",
		"bg","green",
		"width",20,
		"justify","left",
		"relief","groove").Create().Pack("pack", ".")
	
	fmt.Println(w.CGet("fg"))        //get option
	w.CSet("fg","blue", "bg","gray") //set options

	//start main loop
	gotk.Tk.MainLoop()
	gotk.Exit()
}
```
