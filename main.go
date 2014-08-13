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

var devices = make([]uintptr, 0)

var (
	lazyDll = syscall.NewLazyDLL(dllPath)

	procInitialize = lazyDll.NewProc("DirectOutput_Initialize")
	procDeinitialize = lazyDll.NewProc("DirectOutput_Deinitialize")
	procRegisterSoftButtonCallback = lazyDll.NewProc("DirectOutput_RegisterSoftButtonCallback")
	procRegisterPageCallback = lazyDll.NewProc("DirectOutput_RegisterPageCallback")
	procEnumerate = lazyDll.NewProc("DirectOutput_Enumerate")
	procAddPage = lazyDll.NewProc("DirectOutput_AddPage")

	procGetDeviceType = lazyDll.NewProc("DirectOutput_GetDeviceType")
	procSetLed = lazyDll.NewProc("DirectOutput_SetLed")
)

func log(a, b uintptr, err error) {
	fmt.Printf("%#x | %#x | %#x\n", a, b, err)
}

func TestBeep() {
	beepFunc := syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")
	beepFunc.Call(0xffffffff)
}

func SoftButtonChangeCallback(handle, buttons, context uintptr) int {
	fmt.Println("SoftButtonChangeCallback: all the things:", handle, buttons, context)
	return 0
}

func PageCallback(handle, page, activated uintptr) int {
	fmt.Println("PageCallback: all the things:", handle, page, activated)
	return 0
}

func EnumerateCallback(device, ctx uintptr) int {
	devices = append(devices, device)
	fmt.Println("device", device)
	return 0
}

func StrToWideString(s string) uintptr {
	runes := []rune(s)
	p := uintptr(unsafe.Pointer(&runes))
	return p
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

	fmt.Println("procInitialize")
	runes := []rune(appName)
	p := uintptr(unsafe.Pointer(&runes))
	r1, r2, lastErr := procInitialize.Call(p)
	log(r1, r2, lastErr)
	// fmt.Printf("r1, r2, lastErr: %+v, %+v, %+v\n", r1, r2, lastErr)

	fmt.Println("procGetDeviceType")
	r1, r2, lastErr = procGetDeviceType.Call(0)
	log(r1, r2, lastErr)
	// fmt.Printf("r1, r2, lastErr: %+v, %+v, %+v\n", r1, r2, lastErr)

	fmt.Println("procEnumerate")
	r1, r2, lastErr = procEnumerate.Call(syscall.NewCallback(EnumerateCallback))
	log(r1, r2, lastErr)
	// fmt.Printf("r1, r2, lastErr: %+v, %+v, %+v\n", r1, r2, lastErr)

	for len(devices) == 0 {}
	fmt.Println("devices:", devices)

	fmt.Println("procRegisterSoftButtonCallback")
	myCallback := syscall.NewCallbackCDecl(SoftButtonChangeCallback)
	r1, r2, lastErr = procRegisterSoftButtonCallback.Call(devices[0], myCallback, 6666)
	log(r1, r2, lastErr)
	// fmt.Printf("r1, r2, lastErr: %+v, %+v, %+v\n", r1, r2, lastErr)

	myPageNum := 0x00000005
	pageNumPtr := uintptr(unsafe.Pointer(&myPageNum))
	fmt.Println("procAddPage")
	r1, r2, lastErr = procAddPage.Call(devices[0], pageNumPtr, StrToWideString("foo"), 0)
	log(r1, r2, lastErr)
	// fmt.Printf("r1, r2, lastErr: %+v, %+v, %+v\n", r1, r2, lastErr)

	fmt.Println("procRegisterPageCallback")
	myCallback = syscall.NewCallbackCDecl(PageCallback)
	r1, r2, lastErr = procRegisterPageCallback.Call(devices[0], myCallback, 6666);
	log(r1, r2, lastErr)
	// fmt.Printf("r1, r2, lastErr: %+v, %+v, %+v\n", r1, r2, lastErr)

	// fmt.Println("procSetLed")
	// r1, r2, lastErr = procSetLed.Call(devices[0], pageNumPtr, 1, 1)
	// fmt.Printf("%+v\n", r1)
	// fmt.Printf("%+v\n", r2)
	// fmt.Printf("%+v\n", lastErr)

	// fmt.Println("procSetLed")
	// r1, r2, lastErr = procSetLed.Call(devices[0], pageNumPtr, 2, 1)
	// fmt.Printf("%+v\n", r1)
	// fmt.Printf("%+v\n", r2)
	// fmt.Printf("%+v\n", lastErr)


	// go func() {

	// 	for {}
	// }()

	for {}
}