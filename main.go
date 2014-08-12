package main

import (
	"fmt"
	"syscall"
	"os"
	"unsafe"
	// "strings"
)

var deviceTypeX52Pro = "29DAD506-F93B-4F20-85FA-1E02C04FAC17";
var appName = "x52p-mfd-x3-golang"
var dllPath = "DirectOutput.dll"

func TestBeep() {
	beepFunc := syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")
	beepFunc.Call(0xffffffff)
}

func LazyCallProc(a... uintptr) (r1, r2 uintptr, lastErr string) {

}

func SoftButtonChangeCallback(a... uintptr) {
	fmt.Printf("all the things: %q", a)
}

func main() {
	TestBeep()

	lazyDll := syscall.NewLazyDLL(dllPath)
	err := lazyDll.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(lazyDll)

	proc := lazyDll.NewProc("DirectOutput_Initialize")

	fmt.Println(proc)

	runes := []rune(appName)
	p := uintptr(unsafe.Pointer(&runes))
	r1, r2, errNo := proc.Call(p)

	// Takes a long time!

	fmt.Printf("%q\n", r1)
	fmt.Printf("%q\n", r2)
	fmt.Printf("%q\n", errNo)

	// cbackPtr := syscall.NewCallback(SoftButtonChangeCallback)

	// setCallbackProc := lazyDll.NewProc("DirectOutput_RegisterDeviceChangeCallback")
	// fmt.Println(setCallbackProc)

	// setCallbackProc.Call()

	enumerateProc := lazyDll.NewProc("DirectOutput_Enumerate")
	r1, r2, errNo = enumerateProc.Call(0)
	fmt.Printf("%q\n", r1)
	fmt.Printf("%q\n", r2)
	fmt.Printf("%q\n", errNo)

	os.Exit(0)
}