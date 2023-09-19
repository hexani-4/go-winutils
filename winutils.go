package winutils

import (
	"fmt"
	"strings"
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

/*
DWORD cbSize;
  HWND  hWnd;
  UINT  uID;
  UINT  uFlags;
  UINT  uCallbackMessage;
  HICON hIcon;
#if ...
  CHAR  szTip[64];
#else
  CHAR  szTip[128];
#endif
  DWORD dwState;
  DWORD dwStateMask;
  CHAR  szInfo[256];
  union {
    UINT uTimeout;
    UINT uVersion;
  } DUMMYUNIONNAME;
  CHAR  szInfoTitle[64];
  DWORD dwInfoFlags;
  GUID  guidItem;
  HICON hBalloonIcon;
*/

type NOTIFYICONDATA struct {
    CbSize           uint32 //The size of this structure, in bytes. 
    HWnd             syscall.Handle //A handle to the window that receives notifications associated with an icon in the notification area. 
    UID              uint32 //The application-defined identifier of the taskbar icon. The Shell uses either (hWnd plus uID) or guidItem to identify which icon to operate on when Shell_NotifyIcon is invoked. You can have multiple icons associated with a single hWnd by assigning each a different uID. If guidItem is specified, uID is ignored. 
    UFlags           uint32 //Flags that either indicate which of the other members of the structure contain valid data or provide additional information to the tooltip as to how it should display. 
    UCallbackMessage uint32 //An application-defined message identifier. The system uses this identifier to send notification messages to the window identified in hWnd. These notification messages are sent when a mouse event or hover occurs in the bounding rectangle of the icon, when the icon is selected or activated with the keyboard, or when those actions occur in the balloon notification. 
	HIcon            syscall.Handle //A handle to the icon to be added, modified, or deleted. Windows XP and later support icons of up to 32 BPP. 
	SzTip            [128]uint16 //A null-terminated string that specifies the text for a standard tooltip. It can have a maximum of 128 characters (! Windows 2000 and later !), including the terminating null character.
	DwState          uint32 //Windows 2000 and later. The state of the icon. 
	DwStateMask      uint32 //Windows 2000 and later. A value that specifies which bits of the dwState member are retrieved or modified. The possible values are the same as those for dwState. For example, setting this member to NIS_HIDDEN causes only the item's hidden state to be modified while the icon sharing bit is ignored regardless of its value. 
    SzInfo           [256]uint16 //Windows 2000 and later. A null-terminated string that specifies the text to display in a balloon notification. It can have a maximum of 256 characters, including the terminating null character, but should be restricted to 200 characters in English to accommodate localization. To remove the balloon notification from the UI, either delete the icon (with NIM_DELETE) or set the NIF_INFO flag in uFlags and set szInfo to an empty string. 
	UTimeout         uint32 //Deprecated. 
	UVersion         uint32 //Deprecated. 
    SzInfoTitle      [64]uint16 //Windows 2000 and later. A null-terminated string that specifies a title for a balloon notification. This title appears in a larger font immediately above the text. It can have a maximum of 64 characters, including the terminating null character, but should be restricted to 48 characters in English to accommodate localization. 
    DwInfoFlags      uint32 //Windows 2000 and later. Flags that can be set to modify the behavior and appearance of a balloon notification. The icon is placed to the left of the title. If the szInfoTitle member is zero-length, the icon is not shown. 
    GuidItem         syscall.GUID //Windows 7 and later: A registered GUID that identifies the icon. This value overrides uID and is the recommended method of identifying the icon. The NIF_GUID flag must be set in the uFlags member. 
    HBalloonIcon     syscall.Handle //Windows Vista and later. The handle of a customized notification icon provided by the application that should be used independently of the notification area icon. If this member is non-NULL and the NIIF_USER flag is set in the dwInfoFlags member, this icon is used as the notification icon. If this member is NULL, the legacy behavior is carried out. 
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
    nid := &NOTIFYICONDATA{}
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

	for i, c := range utf16_title {
		nid.SzInfoTitle[min(i, 63)] = c
	}
	
	for i, c := range utf16_message {
		nid.SzInfoTitle[min(i, 255)] = c
	}

    // Send the notification
    success, _, err := syscall.SyscallN(procShell_NotifyIconW.Addr(), 0x00000000, uintptr(unsafe.Pointer(nid)))

	fmt.Println(success, "-", nid)
	if success != 1 { return err }
	return err
}

//unicode NULL in string will be changed to "<NULL>" (interpreted by windows as ""). 
func ErrorMessageBox(title string, message string) (err error) {
	const(
		MB_OK = 0x00000000 //The message box contains one push button: OK. This is the default.
		MB_ICONWARNING = 0x00000030 //An exclamation-point icon appears in the message box.
		MB_SYSTEMMODAL = 0x00001000 //Same as MB_APPLMODAL except that the message box has the WS_EX_TOPMOST style. Use system-modal message boxes to notify the user of serious, potentially damaging errors that require immediate attention (for example, running out of memory). This flag has no effect on the user's ability to interact with windows other than those associated with hWnd.
	)
	if procMessageBox == nil { procMessageBox = user32.MustFindProc("MessageBoxW") }

	safe_title := strings.Join(strings.Split(title, "\x00"), "<NULL>")
	safe_message := strings.Join(strings.Split(message, "\x00"), "<NULL>")

	utf16_title, err := syscall.UTF16PtrFromString(safe_title)
	if err != nil { return err }

	utf16_message, err := syscall.UTF16PtrFromString(safe_message)
	if err != nil { return err }

	success, _, err := syscall.SyscallN(procMessageBox.Addr(), uintptr(0), uintptr(unsafe.Pointer(utf16_message)), uintptr(unsafe.Pointer(utf16_title)), uintptr(MB_SYSTEMMODAL | MB_ICONWARNING | MB_OK))

	if success == 0 { return err }
	return nil
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