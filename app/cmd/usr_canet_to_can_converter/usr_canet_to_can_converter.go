package main

import usrcanettocan "github.com/dmnunez1993/usr-canet200-to-can-converter"

func main() {
	_ = usrcanettocan.NewUsrCanetConverter()

	usrcanettocan.LoopForever()
}
