package winutils

import (
	"fmt"
	"syscall"
	"unsafe"
)

const(
	ZPos_BOTTOM uint    = 1 //Places the window at the bottom of the Z order. If the hWnd parameter identifies a topmost window, the window loses its topmost status and is placed at the bottom of all other windows.
	ZPos_NOTOPMOST uint = ^uint(0) - 1 //Places the window above all non-topmost windows (that is, behind all topmost windows). This flag has no effect if the window is already a non-topmost window.
	ZPos_TOP uint       = 0 //Places the window at the top of the Z order.
	ZPos_TOPMOST uint   = ^uint(0) //Places the window above all non-topmost windows. The window maintains its topmost position even when it is deactivated.
)

var(
	user32                = syscall.MustLoadDLL("user32.dll")
	shell32               = syscall.MustLoadDLL("shell32.dll")

	procEnumWindows       *syscall.Proc
	procGetWindowTextW    *syscall.Proc
	procSetWindowPos      *syscall.Proc
	procGetWindowInfo     *syscall.Proc
	procGetWindowRect     *syscall.Proc
	procIsWindow          *syscall.Proc
	procMessageBox        *syscall.Proc
	procShell_NotifyIconW *syscall.Proc
)

type NOTIFYICONDATA struct {
    CbSize           uint32
    HWnd             syscall.Handle
    UID              uint32
    DwState          uint32
    DwStateMask      uint32
    SzInfo           [256]uint16
    UVersion         uint32
    SzInfoTitle      [64]uint16
    DwInfoFlags      uint32
    GuidItem         syscall.GUID
    HIcon            syscall.Handle
    SzTip            [128]uint16
    DwFlags          uint32
    DwCallbackMessage uint32
    HBalloonIcon     syscall.Handle
}

type WINDOWINFO struct {
	CbSize uint32
  	RcWindow DIST_RECT
    RcClient DIST_RECT
  	DwStyle uint32
 	DwExStyle uint32
	DwWindowStatus uint32
  	CxWindowBorders uint32
	CyWindowBorders uint32
  	AtomWindowType uint16
	WCreatorVersion uint16
}

type DIST_RECT struct {
    Left   uint32 //Distance from the left edge of the screen
    Top    uint32 //Distance from the top edge of the screen
    Right  uint32 //Distance from the right edge of the screen
    Bottom uint32 //Distance from the bottom edge of the screen
}

type SIZE_RECT struct {
	X   uint32 //Distance from the left edge of the screen
    Y    uint32 //Distance from the top edge of the screen
	Width uint32 //Width
	Height uint32 //Height
}


func GetWindowInfo(hwnd syscall.Handle) (w_info WINDOWINFO, err error) {
	if procGetWindowInfo == nil { procGetWindowInfo = user32.MustFindProc("GetWindowInfo") }

	var window_info WINDOWINFO

	success, _, err := syscall.SyscallN(procGetWindowInfo.Addr(), uintptr(hwnd), uintptr(unsafe.Pointer(&window_info)))

	if success == 0 { return window_info, err }
	return window_info, nil
}

func GetWindowRect(hwnd syscall.Handle) (window_rect DIST_RECT, err error) {
	if procGetWindowRect == nil { procGetWindowRect = user32.MustFindProc("GetWindowRect") }

	var rect DIST_RECT
	success, _, err := syscall.SyscallN(procGetWindowRect.Addr(), uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))

	if success == 0 { return rect, err }
	return rect, nil
}

func IsWindow(hwnd syscall.Handle) (result bool) {
	if procIsWindow == nil { procIsWindow = user32.MustFindProc("IsWindow") }

	success, _, _ := syscall.SyscallN(procIsWindow.Addr(), uintptr(hwnd))
	
	if success == 0 { return false
	} else { return true }
}

func SetWindowZPos(hwnd syscall.Handle, z_pos uint) (err error) {
	if procSetWindowPos == nil { procSetWindowPos = user32.MustFindProc("SetWindowPos")}

	rect, err := GetWindowRect(hwnd)
	if err != nil { return err }

	success, _, err := syscall.SyscallN(procSetWindowPos.Addr(), uintptr(hwnd), uintptr(z_pos), uintptr(rect.Left), uintptr(rect.Top), uintptr(rect.Right - rect.Left), uintptr(rect.Bottom - rect.Top), uintptr(0x0040))
	
	if success == 0 { return err }
	return nil
}

func SetWindowPos(hwnd syscall.Handle, rect SIZE_RECT) (err error) {
	if procSetWindowPos == nil { procSetWindowPos = user32.MustFindProc("SetWindowPos") }

	success, _, err := syscall.SyscallN(procSetWindowPos.Addr(), uintptr(hwnd), uintptr(0), uintptr(rect.X), uintptr(rect.Y), uintptr(rect.Width), uintptr(rect.Height), uintptr(0x0010 | 0x0004))
	if success == 0 { return err }
	return nil
}

func ErrorMessageNotification(message string) (err error) {
	if procShell_NotifyIconW == nil { procShell_NotifyIconW = shell32.MustFindProc("Shell_NotifyIconW") }

	// Initialize NOTIFYICONDATA
    var nid NOTIFYICONDATA
    nid.CbSize = uint32(unsafe.Sizeof(nid))
    nid.HWnd = syscall.Handle(0)
    nid.UID = 1 // Unique ID for the notification
    nid.DwInfoFlags = 0x00000010
    nid.GuidItem = syscall.GUID{}
    nid.HIcon = 0
    nid.SzTip = [128]uint16{}
    nid.DwState = 0
    nid.DwStateMask = 0
    nid.SzInfoTitle = [64]uint16{}
    nid.SzInfo = [256]uint16{}
    nid.DwInfoFlags |= 0x00000004 // NIIF_USER

    // Set the tooltip (notification title) and info (notification message)
	utf16_title, err := syscall.UTF16FromString("TestTitle")
	if err != nil { return err }

	utf16_message, err := syscall.UTF16FromString(message)
	if err != nil { return err }

	nid.SzInfoTitle = [64]uint16(utf16_title)
	nid.SzInfo = [256]uint16(utf16_message)

    // Send the notification
    success, _, err := syscall.SyscallN(procShell_NotifyIconW.Addr(), 0x00000000, uintptr(unsafe.Pointer(&nid)))

	if success != 0 { return err }
	return err
}

func ErrorMessageBox(title string, message string) (err error) {
	if procMessageBox == nil { procMessageBox = user32.MustFindProc("MessageBox") }

	utf16_title, err := syscall.UTF16FromString(title)
	if err != nil { return err }

	utf16_message, err := syscall.UTF16FromString(message)
	if err != nil { return err }

	_, _, err = syscall.SyscallN(procMessageBox.Addr(), uintptr(unsafe.Pointer(&utf16_message)), uintptr(unsafe.Pointer(&utf16_title)), uintptr(0x00001000))

	return err
}

//thx to EliCDavis (github)
func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	if procEnumWindows == nil { procEnumWindows = user32.MustFindProc("EnumWindows") }

	r1, _, e1 := syscall.SyscallN(procEnumWindows.Addr(), uintptr(enumFunc), uintptr(lparam))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

//thx to EliCDavis (github)
func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	if procGetWindowTextW == nil { procGetWindowTextW = user32.MustFindProc("GetWindowTextW") }

	r0, _, e1 := syscall.SyscallN(procGetWindowTextW.Addr(), uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

//thx to EliCDavis (github)
func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("no window with title '%s' found", title)
	}
	return hwnd, nil
}