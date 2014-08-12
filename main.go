package main

import (
	"fmt"
	"syscall"
	"os"
)

func main() {
	// dll, err := syscall.LoadDLL("kernel32.dll")
	dll := syscall.MustLoadDLL("DirectOutput.dll")
	// if err != nil {
	// 	fmt.Println("Error while loading DirectOutput.dll:", err)
	// 	os.Exit(1)
	// }
	fmt.Println(dll)

	proc, err := dll.FindProc("DirectOutput_Initialize")

	if err != nil {
		fmt.Println("Error while trying to find DirectOutput_Initialize:", err);
		os.Exit(1)
	}

	fmt.Println(proc)

	err = dll.Release()
	if err != nil {
		fmt.Println("error while releasing DLL:", err)
		os.Exit(1)
	}

	os.Exit(0)
}