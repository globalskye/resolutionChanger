package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"unsafe"
)

var (
	user32                     = windows.NewLazySystemDLL("user32.dll")
	procEnumDisplaySettingsW   = user32.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsW = user32.NewProc("ChangeDisplaySettingsW")
)

const (
	ENUM_CURRENT_SETTINGS    = 0xFFFFFFFF
	CDS_UPDATEREGISTRY       = 0x00000001
	CDS_TEST                 = 0x00000002
	CDS_FULLSCREEN           = 0x00000004
	CDS_GLOBAL               = 0x00000008
	CDS_SET_PRIMARY          = 0x00000010
	CDS_VIDEOPARAMETERS      = 0x00000020
	CDS_ENABLE_UNSAFE_MODES  = 0x00000100
	CDS_DISABLE_UNSAFE_MODES = 0x00000200
	CDS_RESET                = 0x40000000
	CDS_NORESET              = 0x10000000
	DISP_CHANGE_SUCCESSFUL   = 0
	DISP_CHANGE_RESTART      = 1
	DISP_CHANGE_FAILED       = -1
)

type DEVMODE struct {
	DmDeviceName         [32]uint16
	DmSpecVersion        uint16
	DmDriverVersion      uint16
	DmSize               uint16
	DmDriverExtra        uint16
	DmFields             uint32
	DmPositionX          int32
	DmPositionY          int32
	DmDisplayOrientation uint32
	DmDisplayFixedOutput uint32
	DmColor              uint16
	DmDuplex             uint16
	DmYResolution        uint16
	DmTTOption           uint16
	DmCollate            uint16
	DmFormName           [32]uint16
	DmLogPixels          uint16
	DmBitsPerPel         uint32
	DmPelsWidth          uint32
	DmPelsHeight         uint32
	DmDisplayFlags       uint32
	DmDisplayFrequency   uint32
	DmICMMethod          uint32
	DmICMIntent          uint32
	DmMediaType          uint32
	DmDitherType         uint32
	DmReserved1          uint32
	DmReserved2          uint32
	DmPanningWidth       uint32
	DmPanningHeight      uint32
}

func getCurrentResolution() (uint32, uint32) {
	var dm DEVMODE
	dm.DmSize = uint16(unsafe.Sizeof(dm))
	ret, _, _ := procEnumDisplaySettingsW.Call(
		0,
		ENUM_CURRENT_SETTINGS,
		uintptr(unsafe.Pointer(&dm)),
	)
	if ret == 0 {
		fmt.Println("Failed to get current display settings")
		os.Exit(1)
	}
	return dm.DmPelsWidth, dm.DmPelsHeight
}

func setResolution(width, height uint32) {
	var dm DEVMODE
	dm.DmSize = uint16(unsafe.Sizeof(dm))
	ret, _, _ := procEnumDisplaySettingsW.Call(
		0,
		ENUM_CURRENT_SETTINGS,
		uintptr(unsafe.Pointer(&dm)),
	)
	if ret == 0 {
		fmt.Println("Failed to get current display settings")
		os.Exit(1)
	}
	dm.DmPelsWidth = width
	dm.DmPelsHeight = height
	dm.DmFields = 0x80000 | 0x100000

	ret, _, _ = procChangeDisplaySettingsW.Call(
		uintptr(unsafe.Pointer(&dm)),
		CDS_UPDATEREGISTRY,
	)
	if ret != DISP_CHANGE_SUCCESSFUL {
		fmt.Println("Failed to change display settings")
		os.Exit(1)
	}
}

func main() {
	width, height := getCurrentResolution()

	if width == 1280 && height == 1024 {
		setResolution(1920, 1080)
	} else {
		setResolution(1280, 1024)
	}
}
