gotk
====

Go-Tk GUI binding

Donate: ** Paypal edartuz@gmail.com **

Example:

``` go
package main

import ( "gotk"; "os"; "fmt";  )

func main() {
	//init Tcl and Tk
	if !gotk.Init(true) { os.Exit(0) }

	//set root geometry
	gotk.Tk.Wm(".", "geometry", "=400x200+100+100")

	w := gotk.Tk.New("button", "", //set type, id (empty string for automatic generation
			"text",    123,        //set other parameters
			"fg",      "red",
			"bg",      "green",
			"width",   20,
			"justify", "left",
			"relief",  "groove").
		Pack("pack", ".",          //pack in root window
			"fill",    "x",        //set other pack options
			"expand",  "0",
			"side",    "bottom").         
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
