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

var (
	lazyDll = syscall.NewLazyDLL(dllPath)

	procInitialize = lazyDll.NewProc("DirectOutput_Initialize")
	procRegisterDeviceChangeCallback = lazyDll.NewProc("DirectOutput_RegisterDeviceChangeCallback")
	procEnumerate = lazyDll.NewProc("DirectOutput_Enumerate")

	procGetDeviceType = lazyDll.NewProc("DirectOutput_GetDeviceType")
)

func TestBeep() {
	beepFunc := syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")
	beepFunc.Call(0xffffffff)
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

	runes := []rune(appName)
	p := uintptr(unsafe.Pointer(&runes))
	r1, r2, errNo := procInitialize.Call(p)
	fmt.Printf("%v\n", r1)
	fmt.Printf("%v\n", r2)
	fmt.Printf("%v\n", errNo)

	r1, r2, errNo = procEnumerate.Call(0)
	fmt.Printf("%+v\n", r1)
	fmt.Printf("%+v\n", r2)
	fmt.Printf("%+v\n", errNo)

	r1, r2, errNo = procGetDeviceType.Call(0)
	fmt.Printf("%+v\n", r1)
	fmt.Printf("%+v\n", r2)
	fmt.Printf("%+v\n", errNo)
}