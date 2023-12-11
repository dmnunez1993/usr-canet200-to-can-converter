package main

import usrcanettocan "github.com/dmnunez1993/usr-canet200-to-can-converter"

func main() {
	converter := usrcanettocan.NewUsrCanetConverter()
	go converter.Run()

	go usrcanettocan.ServeRestApi()

	usrcanettocan.LoopForever()
}
