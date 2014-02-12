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

	w := gotk.Tk.New("button", "", //set type, id (empty string for automatic generation)
			"text",    123,        //set other parameters
			"fg",      "red",
			"bg",      "green",
			"width",   20,
			"justify", "left",
			"relief",  "groove").
		Pack("pack", ".").         //pack in root window
		Cmd("command", func() {    //set callback function
			fmt.Println("Hello...")
	})
	
	fmt.Println(w.CGet("fg"))        //get option
	w.CSet("fg","blue", "bg","gray") //set options

	//start main loop
	gotk.Tk.MainLoop()
	gotk.Exit()
}
```
