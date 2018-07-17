package system 

import (
	"fmt"
	"unsafe"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

var (
	ntuserApiset       = windows.NewLazyDLL("ext-ms-win-ntuser-window-l1-1-0")
	procGetVersionExW  = modkernel32.NewProc("GetVersionExW")
	procGetProductInfo = modkernel32.NewProc("GetProductInfo")
)



type OSVersion struct {
	Version      uint32
	MajorVersion uint8
	MinorVersion uint8
	Build        uint16
}


type osVersionInfoEx struct {
	OSVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformID        uint32
	CSDVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	Reserve           byte
}



func GetOSVersion() OSVersion {
	var err error
	osv := OSVersion{}
	osv.Version, err = windows.GetVersion()
	if err != nil {
		
		panic(err)
	}
	osv.MajorVersion = uint8(osv.Version & 0xFF)
	osv.MinorVersion = uint8(osv.Version >> 8 & 0xFF)
	osv.Build = uint16(osv.Version >> 16)
	return osv
}

func (osv OSVersion) ToString() string {
	return fmt.Sprintf("%d.%d.%d", osv.MajorVersion, osv.MinorVersion, osv.Build)
}




func IsWindowsClient() bool {
	osviex := &osVersionInfoEx{OSVersionInfoSize: 284}
	r1, _, err := procGetVersionExW.Call(uintptr(unsafe.Pointer(osviex)))
	if r1 == 0 {
		logrus.Warnf("GetVersionExW failed - assuming server SKU: %v", err)
		return false
	}
	const verNTWorkstation = 0x00000001
	return osviex.ProductType == verNTWorkstation
}





func IsIoTCore() bool {
	var returnedProductType uint32
	r1, _, err := procGetProductInfo.Call(6, 1, 0, 0, uintptr(unsafe.Pointer(&returnedProductType)))
	if r1 == 0 {
		logrus.Warnf("GetProductInfo failed - assuming this is not IoT: %v", err)
		return false
	}
	const productIoTUAP = 0x0000007B
	const productIoTUAPCommercial = 0x00000083
	return returnedProductType == productIoTUAP || returnedProductType == productIoTUAPCommercial
}



func Unmount(dest string) error {
	return nil
}


func CommandLineToArgv(commandLine string) ([]string, error) {
	var argc int32

	argsPtr, err := windows.UTF16PtrFromString(commandLine)
	if err != nil {
		return nil, err
	}

	argv, err := windows.CommandLineToArgv(argsPtr, &argc)
	if err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(uintptr(unsafe.Pointer(argv))))

	newArgs := make([]string, argc)
	for i, v := range (*argv)[:argc] {
		newArgs[i] = string(windows.UTF16ToString((*v)[:]))
	}

	return newArgs, nil
}



func HasWin32KSupport() bool {
	
	
	
	return ntuserApiset.Load() == nil
}
