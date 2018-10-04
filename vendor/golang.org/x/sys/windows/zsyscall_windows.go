

package windows

import (
	"syscall"
	"unsafe"
)

var _ unsafe.Pointer



const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)



func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	
	
	
	return e
}

var (
	modadvapi32 = NewLazySystemDLL("advapi32.dll")
	modkernel32 = NewLazySystemDLL("kernel32.dll")
	modshell32  = NewLazySystemDLL("shell32.dll")
	modmswsock  = NewLazySystemDLL("mswsock.dll")
	modcrypt32  = NewLazySystemDLL("crypt32.dll")
	modws2_32   = NewLazySystemDLL("ws2_32.dll")
	moddnsapi   = NewLazySystemDLL("dnsapi.dll")
	modiphlpapi = NewLazySystemDLL("iphlpapi.dll")
	modsecur32  = NewLazySystemDLL("secur32.dll")
	modnetapi32 = NewLazySystemDLL("netapi32.dll")
	moduserenv  = NewLazySystemDLL("userenv.dll")

	procRegisterEventSourceW               = modadvapi32.NewProc("RegisterEventSourceW")
	procDeregisterEventSource              = modadvapi32.NewProc("DeregisterEventSource")
	procReportEventW                       = modadvapi32.NewProc("ReportEventW")
	procOpenSCManagerW                     = modadvapi32.NewProc("OpenSCManagerW")
	procCloseServiceHandle                 = modadvapi32.NewProc("CloseServiceHandle")
	procCreateServiceW                     = modadvapi32.NewProc("CreateServiceW")
	procOpenServiceW                       = modadvapi32.NewProc("OpenServiceW")
	procDeleteService                      = modadvapi32.NewProc("DeleteService")
	procStartServiceW                      = modadvapi32.NewProc("StartServiceW")
	procQueryServiceStatus                 = modadvapi32.NewProc("QueryServiceStatus")
	procControlService                     = modadvapi32.NewProc("ControlService")
	procStartServiceCtrlDispatcherW        = modadvapi32.NewProc("StartServiceCtrlDispatcherW")
	procSetServiceStatus                   = modadvapi32.NewProc("SetServiceStatus")
	procChangeServiceConfigW               = modadvapi32.NewProc("ChangeServiceConfigW")
	procQueryServiceConfigW                = modadvapi32.NewProc("QueryServiceConfigW")
	procChangeServiceConfig2W              = modadvapi32.NewProc("ChangeServiceConfig2W")
	procQueryServiceConfig2W               = modadvapi32.NewProc("QueryServiceConfig2W")
	procEnumServicesStatusExW              = modadvapi32.NewProc("EnumServicesStatusExW")
	procQueryServiceStatusEx               = modadvapi32.NewProc("QueryServiceStatusEx")
	procGetLastError                       = modkernel32.NewProc("GetLastError")
	procLoadLibraryW                       = modkernel32.NewProc("LoadLibraryW")
	procLoadLibraryExW                     = modkernel32.NewProc("LoadLibraryExW")
	procFreeLibrary                        = modkernel32.NewProc("FreeLibrary")
	procGetProcAddress                     = modkernel32.NewProc("GetProcAddress")
	procGetVersion                         = modkernel32.NewProc("GetVersion")
	procFormatMessageW                     = modkernel32.NewProc("FormatMessageW")
	procExitProcess                        = modkernel32.NewProc("ExitProcess")
	procCreateFileW                        = modkernel32.NewProc("CreateFileW")
	procReadFile                           = modkernel32.NewProc("ReadFile")
	procWriteFile                          = modkernel32.NewProc("WriteFile")
	procSetFilePointer                     = modkernel32.NewProc("SetFilePointer")
	procCloseHandle                        = modkernel32.NewProc("CloseHandle")
	procGetStdHandle                       = modkernel32.NewProc("GetStdHandle")
	procSetStdHandle                       = modkernel32.NewProc("SetStdHandle")
	procFindFirstFileW                     = modkernel32.NewProc("FindFirstFileW")
	procFindNextFileW                      = modkernel32.NewProc("FindNextFileW")
	procFindClose                          = modkernel32.NewProc("FindClose")
	procGetFileInformationByHandle         = modkernel32.NewProc("GetFileInformationByHandle")
	procGetCurrentDirectoryW               = modkernel32.NewProc("GetCurrentDirectoryW")
	procSetCurrentDirectoryW               = modkernel32.NewProc("SetCurrentDirectoryW")
	procCreateDirectoryW                   = modkernel32.NewProc("CreateDirectoryW")
	procRemoveDirectoryW                   = modkernel32.NewProc("RemoveDirectoryW")
	procDeleteFileW                        = modkernel32.NewProc("DeleteFileW")
	procMoveFileW                          = modkernel32.NewProc("MoveFileW")
	procMoveFileExW                        = modkernel32.NewProc("MoveFileExW")
	procGetComputerNameW                   = modkernel32.NewProc("GetComputerNameW")
	procGetComputerNameExW                 = modkernel32.NewProc("GetComputerNameExW")
	procSetEndOfFile                       = modkernel32.NewProc("SetEndOfFile")
	procGetSystemTimeAsFileTime            = modkernel32.NewProc("GetSystemTimeAsFileTime")
	procGetSystemTimePreciseAsFileTime     = modkernel32.NewProc("GetSystemTimePreciseAsFileTime")
	procGetTimeZoneInformation             = modkernel32.NewProc("GetTimeZoneInformation")
	procCreateIoCompletionPort             = modkernel32.NewProc("CreateIoCompletionPort")
	procGetQueuedCompletionStatus          = modkernel32.NewProc("GetQueuedCompletionStatus")
	procPostQueuedCompletionStatus         = modkernel32.NewProc("PostQueuedCompletionStatus")
	procCancelIo                           = modkernel32.NewProc("CancelIo")
	procCancelIoEx                         = modkernel32.NewProc("CancelIoEx")
	procCreateProcessW                     = modkernel32.NewProc("CreateProcessW")
	procOpenProcess                        = modkernel32.NewProc("OpenProcess")
	procTerminateProcess                   = modkernel32.NewProc("TerminateProcess")
	procGetExitCodeProcess                 = modkernel32.NewProc("GetExitCodeProcess")
	procGetStartupInfoW                    = modkernel32.NewProc("GetStartupInfoW")
	procGetCurrentProcess                  = modkernel32.NewProc("GetCurrentProcess")
	procGetProcessTimes                    = modkernel32.NewProc("GetProcessTimes")
	procDuplicateHandle                    = modkernel32.NewProc("DuplicateHandle")
	procWaitForSingleObject                = modkernel32.NewProc("WaitForSingleObject")
	procGetTempPathW                       = modkernel32.NewProc("GetTempPathW")
	procCreatePipe                         = modkernel32.NewProc("CreatePipe")
	procGetFileType                        = modkernel32.NewProc("GetFileType")
	procCryptAcquireContextW               = modadvapi32.NewProc("CryptAcquireContextW")
	procCryptReleaseContext                = modadvapi32.NewProc("CryptReleaseContext")
	procCryptGenRandom                     = modadvapi32.NewProc("CryptGenRandom")
	procGetEnvironmentStringsW             = modkernel32.NewProc("GetEnvironmentStringsW")
	procFreeEnvironmentStringsW            = modkernel32.NewProc("FreeEnvironmentStringsW")
	procGetEnvironmentVariableW            = modkernel32.NewProc("GetEnvironmentVariableW")
	procSetEnvironmentVariableW            = modkernel32.NewProc("SetEnvironmentVariableW")
	procSetFileTime                        = modkernel32.NewProc("SetFileTime")
	procGetFileAttributesW                 = modkernel32.NewProc("GetFileAttributesW")
	procSetFileAttributesW                 = modkernel32.NewProc("SetFileAttributesW")
	procGetFileAttributesExW               = modkernel32.NewProc("GetFileAttributesExW")
	procGetCommandLineW                    = modkernel32.NewProc("GetCommandLineW")
	procCommandLineToArgvW                 = modshell32.NewProc("CommandLineToArgvW")
	procLocalFree                          = modkernel32.NewProc("LocalFree")
	procSetHandleInformation               = modkernel32.NewProc("SetHandleInformation")
	procFlushFileBuffers                   = modkernel32.NewProc("FlushFileBuffers")
	procGetFullPathNameW                   = modkernel32.NewProc("GetFullPathNameW")
	procGetLongPathNameW                   = modkernel32.NewProc("GetLongPathNameW")
	procGetShortPathNameW                  = modkernel32.NewProc("GetShortPathNameW")
	procCreateFileMappingW                 = modkernel32.NewProc("CreateFileMappingW")
	procMapViewOfFile                      = modkernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile                    = modkernel32.NewProc("UnmapViewOfFile")
	procFlushViewOfFile                    = modkernel32.NewProc("FlushViewOfFile")
	procVirtualLock                        = modkernel32.NewProc("VirtualLock")
	procVirtualUnlock                      = modkernel32.NewProc("VirtualUnlock")
	procVirtualAlloc                       = modkernel32.NewProc("VirtualAlloc")
	procVirtualFree                        = modkernel32.NewProc("VirtualFree")
	procVirtualProtect                     = modkernel32.NewProc("VirtualProtect")
	procTransmitFile                       = modmswsock.NewProc("TransmitFile")
	procReadDirectoryChangesW              = modkernel32.NewProc("ReadDirectoryChangesW")
	procCertOpenSystemStoreW               = modcrypt32.NewProc("CertOpenSystemStoreW")
	procCertOpenStore                      = modcrypt32.NewProc("CertOpenStore")
	procCertEnumCertificatesInStore        = modcrypt32.NewProc("CertEnumCertificatesInStore")
	procCertAddCertificateContextToStore   = modcrypt32.NewProc("CertAddCertificateContextToStore")
	procCertCloseStore                     = modcrypt32.NewProc("CertCloseStore")
	procCertGetCertificateChain            = modcrypt32.NewProc("CertGetCertificateChain")
	procCertFreeCertificateChain           = modcrypt32.NewProc("CertFreeCertificateChain")
	procCertCreateCertificateContext       = modcrypt32.NewProc("CertCreateCertificateContext")
	procCertFreeCertificateContext         = modcrypt32.NewProc("CertFreeCertificateContext")
	procCertVerifyCertificateChainPolicy   = modcrypt32.NewProc("CertVerifyCertificateChainPolicy")
	procRegOpenKeyExW                      = modadvapi32.NewProc("RegOpenKeyExW")
	procRegCloseKey                        = modadvapi32.NewProc("RegCloseKey")
	procRegQueryInfoKeyW                   = modadvapi32.NewProc("RegQueryInfoKeyW")
	procRegEnumKeyExW                      = modadvapi32.NewProc("RegEnumKeyExW")
	procRegQueryValueExW                   = modadvapi32.NewProc("RegQueryValueExW")
	procGetCurrentProcessId                = modkernel32.NewProc("GetCurrentProcessId")
	procGetConsoleMode                     = modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode                     = modkernel32.NewProc("SetConsoleMode")
	procGetConsoleScreenBufferInfo         = modkernel32.NewProc("GetConsoleScreenBufferInfo")
	procWriteConsoleW                      = modkernel32.NewProc("WriteConsoleW")
	procReadConsoleW                       = modkernel32.NewProc("ReadConsoleW")
	procCreateToolhelp32Snapshot           = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW                    = modkernel32.NewProc("Process32FirstW")
	procProcess32NextW                     = modkernel32.NewProc("Process32NextW")
	procDeviceIoControl                    = modkernel32.NewProc("DeviceIoControl")
	procCreateSymbolicLinkW                = modkernel32.NewProc("CreateSymbolicLinkW")
	procCreateHardLinkW                    = modkernel32.NewProc("CreateHardLinkW")
	procGetCurrentThreadId                 = modkernel32.NewProc("GetCurrentThreadId")
	procCreateEventW                       = modkernel32.NewProc("CreateEventW")
	procCreateEventExW                     = modkernel32.NewProc("CreateEventExW")
	procOpenEventW                         = modkernel32.NewProc("OpenEventW")
	procSetEvent                           = modkernel32.NewProc("SetEvent")
	procResetEvent                         = modkernel32.NewProc("ResetEvent")
	procPulseEvent                         = modkernel32.NewProc("PulseEvent")
	procDefineDosDeviceW                   = modkernel32.NewProc("DefineDosDeviceW")
	procDeleteVolumeMountPointW            = modkernel32.NewProc("DeleteVolumeMountPointW")
	procFindFirstVolumeW                   = modkernel32.NewProc("FindFirstVolumeW")
	procFindFirstVolumeMountPointW         = modkernel32.NewProc("FindFirstVolumeMountPointW")
	procFindNextVolumeW                    = modkernel32.NewProc("FindNextVolumeW")
	procFindNextVolumeMountPointW          = modkernel32.NewProc("FindNextVolumeMountPointW")
	procFindVolumeClose                    = modkernel32.NewProc("FindVolumeClose")
	procFindVolumeMountPointClose          = modkernel32.NewProc("FindVolumeMountPointClose")
	procGetDriveTypeW                      = modkernel32.NewProc("GetDriveTypeW")
	procGetLogicalDrives                   = modkernel32.NewProc("GetLogicalDrives")
	procGetLogicalDriveStringsW            = modkernel32.NewProc("GetLogicalDriveStringsW")
	procGetVolumeInformationW              = modkernel32.NewProc("GetVolumeInformationW")
	procGetVolumeInformationByHandleW      = modkernel32.NewProc("GetVolumeInformationByHandleW")
	procGetVolumeNameForVolumeMountPointW  = modkernel32.NewProc("GetVolumeNameForVolumeMountPointW")
	procGetVolumePathNameW                 = modkernel32.NewProc("GetVolumePathNameW")
	procGetVolumePathNamesForVolumeNameW   = modkernel32.NewProc("GetVolumePathNamesForVolumeNameW")
	procQueryDosDeviceW                    = modkernel32.NewProc("QueryDosDeviceW")
	procSetVolumeLabelW                    = modkernel32.NewProc("SetVolumeLabelW")
	procSetVolumeMountPointW               = modkernel32.NewProc("SetVolumeMountPointW")
	procWSAStartup                         = modws2_32.NewProc("WSAStartup")
	procWSACleanup                         = modws2_32.NewProc("WSACleanup")
	procWSAIoctl                           = modws2_32.NewProc("WSAIoctl")
	procsocket                             = modws2_32.NewProc("socket")
	procsetsockopt                         = modws2_32.NewProc("setsockopt")
	procgetsockopt                         = modws2_32.NewProc("getsockopt")
	procbind                               = modws2_32.NewProc("bind")
	procconnect                            = modws2_32.NewProc("connect")
	procgetsockname                        = modws2_32.NewProc("getsockname")
	procgetpeername                        = modws2_32.NewProc("getpeername")
	proclisten                             = modws2_32.NewProc("listen")
	procshutdown                           = modws2_32.NewProc("shutdown")
	procclosesocket                        = modws2_32.NewProc("closesocket")
	procAcceptEx                           = modmswsock.NewProc("AcceptEx")
	procGetAcceptExSockaddrs               = modmswsock.NewProc("GetAcceptExSockaddrs")
	procWSARecv                            = modws2_32.NewProc("WSARecv")
	procWSASend                            = modws2_32.NewProc("WSASend")
	procWSARecvFrom                        = modws2_32.NewProc("WSARecvFrom")
	procWSASendTo                          = modws2_32.NewProc("WSASendTo")
	procgethostbyname                      = modws2_32.NewProc("gethostbyname")
	procgetservbyname                      = modws2_32.NewProc("getservbyname")
	procntohs                              = modws2_32.NewProc("ntohs")
	procgetprotobyname                     = modws2_32.NewProc("getprotobyname")
	procDnsQuery_W                         = moddnsapi.NewProc("DnsQuery_W")
	procDnsRecordListFree                  = moddnsapi.NewProc("DnsRecordListFree")
	procDnsNameCompare_W                   = moddnsapi.NewProc("DnsNameCompare_W")
	procGetAddrInfoW                       = modws2_32.NewProc("GetAddrInfoW")
	procFreeAddrInfoW                      = modws2_32.NewProc("FreeAddrInfoW")
	procGetIfEntry                         = modiphlpapi.NewProc("GetIfEntry")
	procGetAdaptersInfo                    = modiphlpapi.NewProc("GetAdaptersInfo")
	procSetFileCompletionNotificationModes = modkernel32.NewProc("SetFileCompletionNotificationModes")
	procWSAEnumProtocolsW                  = modws2_32.NewProc("WSAEnumProtocolsW")
	procGetAdaptersAddresses               = modiphlpapi.NewProc("GetAdaptersAddresses")
	procGetACP                             = modkernel32.NewProc("GetACP")
	procMultiByteToWideChar                = modkernel32.NewProc("MultiByteToWideChar")
	procTranslateNameW                     = modsecur32.NewProc("TranslateNameW")
	procGetUserNameExW                     = modsecur32.NewProc("GetUserNameExW")
	procNetUserGetInfo                     = modnetapi32.NewProc("NetUserGetInfo")
	procNetGetJoinInformation              = modnetapi32.NewProc("NetGetJoinInformation")
	procNetApiBufferFree                   = modnetapi32.NewProc("NetApiBufferFree")
	procLookupAccountSidW                  = modadvapi32.NewProc("LookupAccountSidW")
	procLookupAccountNameW                 = modadvapi32.NewProc("LookupAccountNameW")
	procConvertSidToStringSidW             = modadvapi32.NewProc("ConvertSidToStringSidW")
	procConvertStringSidToSidW             = modadvapi32.NewProc("ConvertStringSidToSidW")
	procGetLengthSid                       = modadvapi32.NewProc("GetLengthSid")
	procCopySid                            = modadvapi32.NewProc("CopySid")
	procAllocateAndInitializeSid           = modadvapi32.NewProc("AllocateAndInitializeSid")
	procFreeSid                            = modadvapi32.NewProc("FreeSid")
	procEqualSid                           = modadvapi32.NewProc("EqualSid")
	procCheckTokenMembership               = modadvapi32.NewProc("CheckTokenMembership")
	procOpenProcessToken                   = modadvapi32.NewProc("OpenProcessToken")
	procGetTokenInformation                = modadvapi32.NewProc("GetTokenInformation")
	procGetUserProfileDirectoryW           = moduserenv.NewProc("GetUserProfileDirectoryW")
)

func RegisterEventSource(uncServerName *uint16, sourceName *uint16) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procRegisterEventSourceW.Addr(), 2, uintptr(unsafe.Pointer(uncServerName)), uintptr(unsafe.Pointer(sourceName)), 0)
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DeregisterEventSource(handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procDeregisterEventSource.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ReportEvent(log Handle, etype uint16, category uint16, eventId uint32, usrSId uintptr, numStrings uint16, dataSize uint32, strings **uint16, rawData *byte) (err error) {
	r1, _, e1 := syscall.Syscall9(procReportEventW.Addr(), 9, uintptr(log), uintptr(etype), uintptr(category), uintptr(eventId), uintptr(usrSId), uintptr(numStrings), uintptr(dataSize), uintptr(unsafe.Pointer(strings)), uintptr(unsafe.Pointer(rawData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenSCManager(machineName *uint16, databaseName *uint16, access uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procOpenSCManagerW.Addr(), 3, uintptr(unsafe.Pointer(machineName)), uintptr(unsafe.Pointer(databaseName)), uintptr(access))
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CloseServiceHandle(handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procCloseServiceHandle.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateService(mgr Handle, serviceName *uint16, displayName *uint16, access uint32, srvType uint32, startType uint32, errCtl uint32, pathName *uint16, loadOrderGroup *uint16, tagId *uint32, dependencies *uint16, serviceStartName *uint16, password *uint16) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall15(procCreateServiceW.Addr(), 13, uintptr(mgr), uintptr(unsafe.Pointer(serviceName)), uintptr(unsafe.Pointer(displayName)), uintptr(access), uintptr(srvType), uintptr(startType), uintptr(errCtl), uintptr(unsafe.Pointer(pathName)), uintptr(unsafe.Pointer(loadOrderGroup)), uintptr(unsafe.Pointer(tagId)), uintptr(unsafe.Pointer(dependencies)), uintptr(unsafe.Pointer(serviceStartName)), uintptr(unsafe.Pointer(password)), 0, 0)
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenService(mgr Handle, serviceName *uint16, access uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procOpenServiceW.Addr(), 3, uintptr(mgr), uintptr(unsafe.Pointer(serviceName)), uintptr(access))
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DeleteService(service Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procDeleteService.Addr(), 1, uintptr(service), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func StartService(service Handle, numArgs uint32, argVectors **uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procStartServiceW.Addr(), 3, uintptr(service), uintptr(numArgs), uintptr(unsafe.Pointer(argVectors)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func QueryServiceStatus(service Handle, status *SERVICE_STATUS) (err error) {
	r1, _, e1 := syscall.Syscall(procQueryServiceStatus.Addr(), 2, uintptr(service), uintptr(unsafe.Pointer(status)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ControlService(service Handle, control uint32, status *SERVICE_STATUS) (err error) {
	r1, _, e1 := syscall.Syscall(procControlService.Addr(), 3, uintptr(service), uintptr(control), uintptr(unsafe.Pointer(status)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func StartServiceCtrlDispatcher(serviceTable *SERVICE_TABLE_ENTRY) (err error) {
	r1, _, e1 := syscall.Syscall(procStartServiceCtrlDispatcherW.Addr(), 1, uintptr(unsafe.Pointer(serviceTable)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetServiceStatus(service Handle, serviceStatus *SERVICE_STATUS) (err error) {
	r1, _, e1 := syscall.Syscall(procSetServiceStatus.Addr(), 2, uintptr(service), uintptr(unsafe.Pointer(serviceStatus)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ChangeServiceConfig(service Handle, serviceType uint32, startType uint32, errorControl uint32, binaryPathName *uint16, loadOrderGroup *uint16, tagId *uint32, dependencies *uint16, serviceStartName *uint16, password *uint16, displayName *uint16) (err error) {
	r1, _, e1 := syscall.Syscall12(procChangeServiceConfigW.Addr(), 11, uintptr(service), uintptr(serviceType), uintptr(startType), uintptr(errorControl), uintptr(unsafe.Pointer(binaryPathName)), uintptr(unsafe.Pointer(loadOrderGroup)), uintptr(unsafe.Pointer(tagId)), uintptr(unsafe.Pointer(dependencies)), uintptr(unsafe.Pointer(serviceStartName)), uintptr(unsafe.Pointer(password)), uintptr(unsafe.Pointer(displayName)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func QueryServiceConfig(service Handle, serviceConfig *QUERY_SERVICE_CONFIG, bufSize uint32, bytesNeeded *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procQueryServiceConfigW.Addr(), 4, uintptr(service), uintptr(unsafe.Pointer(serviceConfig)), uintptr(bufSize), uintptr(unsafe.Pointer(bytesNeeded)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ChangeServiceConfig2(service Handle, infoLevel uint32, info *byte) (err error) {
	r1, _, e1 := syscall.Syscall(procChangeServiceConfig2W.Addr(), 3, uintptr(service), uintptr(infoLevel), uintptr(unsafe.Pointer(info)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func QueryServiceConfig2(service Handle, infoLevel uint32, buff *byte, buffSize uint32, bytesNeeded *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procQueryServiceConfig2W.Addr(), 5, uintptr(service), uintptr(infoLevel), uintptr(unsafe.Pointer(buff)), uintptr(buffSize), uintptr(unsafe.Pointer(bytesNeeded)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func EnumServicesStatusEx(mgr Handle, infoLevel uint32, serviceType uint32, serviceState uint32, services *byte, bufSize uint32, bytesNeeded *uint32, servicesReturned *uint32, resumeHandle *uint32, groupName *uint16) (err error) {
	r1, _, e1 := syscall.Syscall12(procEnumServicesStatusExW.Addr(), 10, uintptr(mgr), uintptr(infoLevel), uintptr(serviceType), uintptr(serviceState), uintptr(unsafe.Pointer(services)), uintptr(bufSize), uintptr(unsafe.Pointer(bytesNeeded)), uintptr(unsafe.Pointer(servicesReturned)), uintptr(unsafe.Pointer(resumeHandle)), uintptr(unsafe.Pointer(groupName)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func QueryServiceStatusEx(service Handle, infoLevel uint32, buff *byte, buffSize uint32, bytesNeeded *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procQueryServiceStatusEx.Addr(), 5, uintptr(service), uintptr(infoLevel), uintptr(unsafe.Pointer(buff)), uintptr(buffSize), uintptr(unsafe.Pointer(bytesNeeded)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetLastError() (lasterr error) {
	r0, _, _ := syscall.Syscall(procGetLastError.Addr(), 0, 0, 0, 0)
	if r0 != 0 {
		lasterr = syscall.Errno(r0)
	}
	return
}

func LoadLibrary(libname string) (handle Handle, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(libname)
	if err != nil {
		return
	}
	return _LoadLibrary(_p0)
}

func _LoadLibrary(libname *uint16) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procLoadLibraryW.Addr(), 1, uintptr(unsafe.Pointer(libname)), 0, 0)
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func LoadLibraryEx(libname string, zero Handle, flags uintptr) (handle Handle, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(libname)
	if err != nil {
		return
	}
	return _LoadLibraryEx(_p0, zero, flags)
}

func _LoadLibraryEx(libname *uint16, zero Handle, flags uintptr) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procLoadLibraryExW.Addr(), 3, uintptr(unsafe.Pointer(libname)), uintptr(zero), uintptr(flags))
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FreeLibrary(handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procFreeLibrary.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetProcAddress(module Handle, procname string) (proc uintptr, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(procname)
	if err != nil {
		return
	}
	return _GetProcAddress(module, _p0)
}

func _GetProcAddress(module Handle, procname *byte) (proc uintptr, err error) {
	r0, _, e1 := syscall.Syscall(procGetProcAddress.Addr(), 2, uintptr(module), uintptr(unsafe.Pointer(procname)), 0)
	proc = uintptr(r0)
	if proc == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetVersion() (ver uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetVersion.Addr(), 0, 0, 0, 0)
	ver = uint32(r0)
	if ver == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FormatMessage(flags uint32, msgsrc uintptr, msgid uint32, langid uint32, buf []uint16, args *byte) (n uint32, err error) {
	var _p0 *uint16
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r0, _, e1 := syscall.Syscall9(procFormatMessageW.Addr(), 7, uintptr(flags), uintptr(msgsrc), uintptr(msgid), uintptr(langid), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(args)), 0, 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ExitProcess(exitcode uint32) {
	syscall.Syscall(procExitProcess.Addr(), 1, uintptr(exitcode), 0, 0)
	return
}

func CreateFile(name *uint16, access uint32, mode uint32, sa *SecurityAttributes, createmode uint32, attrs uint32, templatefile int32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall9(procCreateFileW.Addr(), 7, uintptr(unsafe.Pointer(name)), uintptr(access), uintptr(mode), uintptr(unsafe.Pointer(sa)), uintptr(createmode), uintptr(attrs), uintptr(templatefile), 0, 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ReadFile(handle Handle, buf []byte, done *uint32, overlapped *Overlapped) (err error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r1, _, e1 := syscall.Syscall6(procReadFile.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(done)), uintptr(unsafe.Pointer(overlapped)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WriteFile(handle Handle, buf []byte, done *uint32, overlapped *Overlapped) (err error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r1, _, e1 := syscall.Syscall6(procWriteFile.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(done)), uintptr(unsafe.Pointer(overlapped)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetFilePointer(handle Handle, lowoffset int32, highoffsetptr *int32, whence uint32) (newlowoffset uint32, err error) {
	r0, _, e1 := syscall.Syscall6(procSetFilePointer.Addr(), 4, uintptr(handle), uintptr(lowoffset), uintptr(unsafe.Pointer(highoffsetptr)), uintptr(whence), 0, 0)
	newlowoffset = uint32(r0)
	if newlowoffset == 0xffffffff {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CloseHandle(handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procCloseHandle.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetStdHandle(stdhandle uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetStdHandle.Addr(), 1, uintptr(stdhandle), 0, 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetStdHandle(stdhandle uint32, handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func findFirstFile1(name *uint16, data *win32finddata1) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procFindFirstFileW.Addr(), 2, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(data)), 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func findNextFile1(handle Handle, data *win32finddata1) (err error) {
	r1, _, e1 := syscall.Syscall(procFindNextFileW.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(data)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindClose(handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procFindClose.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetFileInformationByHandle(handle Handle, data *ByHandleFileInformation) (err error) {
	r1, _, e1 := syscall.Syscall(procGetFileInformationByHandle.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(data)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetCurrentDirectory(buflen uint32, buf *uint16) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetCurrentDirectoryW.Addr(), 2, uintptr(buflen), uintptr(unsafe.Pointer(buf)), 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetCurrentDirectory(path *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procSetCurrentDirectoryW.Addr(), 1, uintptr(unsafe.Pointer(path)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateDirectory(path *uint16, sa *SecurityAttributes) (err error) {
	r1, _, e1 := syscall.Syscall(procCreateDirectoryW.Addr(), 2, uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(sa)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func RemoveDirectory(path *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procRemoveDirectoryW.Addr(), 1, uintptr(unsafe.Pointer(path)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DeleteFile(path *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procDeleteFileW.Addr(), 1, uintptr(unsafe.Pointer(path)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func MoveFile(from *uint16, to *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procMoveFileW.Addr(), 2, uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(to)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func MoveFileEx(from *uint16, to *uint16, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procMoveFileExW.Addr(), 3, uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(to)), uintptr(flags))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetComputerName(buf *uint16, n *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetComputerNameW.Addr(), 2, uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(n)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetComputerNameEx(nametype uint32, buf *uint16, n *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetComputerNameExW.Addr(), 3, uintptr(nametype), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(n)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetEndOfFile(handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procSetEndOfFile.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetSystemTimeAsFileTime(time *Filetime) {
	syscall.Syscall(procGetSystemTimeAsFileTime.Addr(), 1, uintptr(unsafe.Pointer(time)), 0, 0)
	return
}

func GetSystemTimePreciseAsFileTime(time *Filetime) {
	syscall.Syscall(procGetSystemTimePreciseAsFileTime.Addr(), 1, uintptr(unsafe.Pointer(time)), 0, 0)
	return
}

func GetTimeZoneInformation(tzi *Timezoneinformation) (rc uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetTimeZoneInformation.Addr(), 1, uintptr(unsafe.Pointer(tzi)), 0, 0)
	rc = uint32(r0)
	if rc == 0xffffffff {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateIoCompletionPort(filehandle Handle, cphandle Handle, key uint32, threadcnt uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procCreateIoCompletionPort.Addr(), 4, uintptr(filehandle), uintptr(cphandle), uintptr(key), uintptr(threadcnt), 0, 0)
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetQueuedCompletionStatus(cphandle Handle, qty *uint32, key *uint32, overlapped **Overlapped, timeout uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procGetQueuedCompletionStatus.Addr(), 5, uintptr(cphandle), uintptr(unsafe.Pointer(qty)), uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(overlapped)), uintptr(timeout), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func PostQueuedCompletionStatus(cphandle Handle, qty uint32, key uint32, overlapped *Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall6(procPostQueuedCompletionStatus.Addr(), 4, uintptr(cphandle), uintptr(qty), uintptr(key), uintptr(unsafe.Pointer(overlapped)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CancelIo(s Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procCancelIo.Addr(), 1, uintptr(s), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CancelIoEx(s Handle, o *Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall(procCancelIoEx.Addr(), 2, uintptr(s), uintptr(unsafe.Pointer(o)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateProcess(appName *uint16, commandLine *uint16, procSecurity *SecurityAttributes, threadSecurity *SecurityAttributes, inheritHandles bool, creationFlags uint32, env *uint16, currentDir *uint16, startupInfo *StartupInfo, outProcInfo *ProcessInformation) (err error) {
	var _p0 uint32
	if inheritHandles {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := syscall.Syscall12(procCreateProcessW.Addr(), 10, uintptr(unsafe.Pointer(appName)), uintptr(unsafe.Pointer(commandLine)), uintptr(unsafe.Pointer(procSecurity)), uintptr(unsafe.Pointer(threadSecurity)), uintptr(_p0), uintptr(creationFlags), uintptr(unsafe.Pointer(env)), uintptr(unsafe.Pointer(currentDir)), uintptr(unsafe.Pointer(startupInfo)), uintptr(unsafe.Pointer(outProcInfo)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenProcess(da uint32, inheritHandle bool, pid uint32) (handle Handle, err error) {
	var _p0 uint32
	if inheritHandle {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := syscall.Syscall(procOpenProcess.Addr(), 3, uintptr(da), uintptr(_p0), uintptr(pid))
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func TerminateProcess(handle Handle, exitcode uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procTerminateProcess.Addr(), 2, uintptr(handle), uintptr(exitcode), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetExitCodeProcess(handle Handle, exitcode *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetExitCodeProcess.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(exitcode)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetStartupInfo(startupInfo *StartupInfo) (err error) {
	r1, _, e1 := syscall.Syscall(procGetStartupInfoW.Addr(), 1, uintptr(unsafe.Pointer(startupInfo)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetCurrentProcess() (pseudoHandle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetCurrentProcess.Addr(), 0, 0, 0, 0)
	pseudoHandle = Handle(r0)
	if pseudoHandle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetProcessTimes(handle Handle, creationTime *Filetime, exitTime *Filetime, kernelTime *Filetime, userTime *Filetime) (err error) {
	r1, _, e1 := syscall.Syscall6(procGetProcessTimes.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(creationTime)), uintptr(unsafe.Pointer(exitTime)), uintptr(unsafe.Pointer(kernelTime)), uintptr(unsafe.Pointer(userTime)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DuplicateHandle(hSourceProcessHandle Handle, hSourceHandle Handle, hTargetProcessHandle Handle, lpTargetHandle *Handle, dwDesiredAccess uint32, bInheritHandle bool, dwOptions uint32) (err error) {
	var _p0 uint32
	if bInheritHandle {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := syscall.Syscall9(procDuplicateHandle.Addr(), 7, uintptr(hSourceProcessHandle), uintptr(hSourceHandle), uintptr(hTargetProcessHandle), uintptr(unsafe.Pointer(lpTargetHandle)), uintptr(dwDesiredAccess), uintptr(_p0), uintptr(dwOptions), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WaitForSingleObject(handle Handle, waitMilliseconds uint32) (event uint32, err error) {
	r0, _, e1 := syscall.Syscall(procWaitForSingleObject.Addr(), 2, uintptr(handle), uintptr(waitMilliseconds), 0)
	event = uint32(r0)
	if event == 0xffffffff {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetTempPath(buflen uint32, buf *uint16) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetTempPathW.Addr(), 2, uintptr(buflen), uintptr(unsafe.Pointer(buf)), 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreatePipe(readhandle *Handle, writehandle *Handle, sa *SecurityAttributes, size uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procCreatePipe.Addr(), 4, uintptr(unsafe.Pointer(readhandle)), uintptr(unsafe.Pointer(writehandle)), uintptr(unsafe.Pointer(sa)), uintptr(size), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetFileType(filehandle Handle) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetFileType.Addr(), 1, uintptr(filehandle), 0, 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CryptAcquireContext(provhandle *Handle, container *uint16, provider *uint16, provtype uint32, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procCryptAcquireContextW.Addr(), 5, uintptr(unsafe.Pointer(provhandle)), uintptr(unsafe.Pointer(container)), uintptr(unsafe.Pointer(provider)), uintptr(provtype), uintptr(flags), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CryptReleaseContext(provhandle Handle, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procCryptReleaseContext.Addr(), 2, uintptr(provhandle), uintptr(flags), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CryptGenRandom(provhandle Handle, buflen uint32, buf *byte) (err error) {
	r1, _, e1 := syscall.Syscall(procCryptGenRandom.Addr(), 3, uintptr(provhandle), uintptr(buflen), uintptr(unsafe.Pointer(buf)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetEnvironmentStrings() (envs *uint16, err error) {
	r0, _, e1 := syscall.Syscall(procGetEnvironmentStringsW.Addr(), 0, 0, 0, 0)
	envs = (*uint16)(unsafe.Pointer(r0))
	if envs == nil {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FreeEnvironmentStrings(envs *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procFreeEnvironmentStringsW.Addr(), 1, uintptr(unsafe.Pointer(envs)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetEnvironmentVariable(name *uint16, buffer *uint16, size uint32) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetEnvironmentVariableW.Addr(), 3, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(buffer)), uintptr(size))
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetEnvironmentVariable(name *uint16, value *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procSetEnvironmentVariableW.Addr(), 2, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(value)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetFileTime(handle Handle, ctime *Filetime, atime *Filetime, wtime *Filetime) (err error) {
	r1, _, e1 := syscall.Syscall6(procSetFileTime.Addr(), 4, uintptr(handle), uintptr(unsafe.Pointer(ctime)), uintptr(unsafe.Pointer(atime)), uintptr(unsafe.Pointer(wtime)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetFileAttributes(name *uint16) (attrs uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetFileAttributesW.Addr(), 1, uintptr(unsafe.Pointer(name)), 0, 0)
	attrs = uint32(r0)
	if attrs == INVALID_FILE_ATTRIBUTES {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetFileAttributes(name *uint16, attrs uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procSetFileAttributesW.Addr(), 2, uintptr(unsafe.Pointer(name)), uintptr(attrs), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetFileAttributesEx(name *uint16, level uint32, info *byte) (err error) {
	r1, _, e1 := syscall.Syscall(procGetFileAttributesExW.Addr(), 3, uintptr(unsafe.Pointer(name)), uintptr(level), uintptr(unsafe.Pointer(info)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetCommandLine() (cmd *uint16) {
	r0, _, _ := syscall.Syscall(procGetCommandLineW.Addr(), 0, 0, 0, 0)
	cmd = (*uint16)(unsafe.Pointer(r0))
	return
}

func CommandLineToArgv(cmd *uint16, argc *int32) (argv *[8192]*[8192]uint16, err error) {
	r0, _, e1 := syscall.Syscall(procCommandLineToArgvW.Addr(), 2, uintptr(unsafe.Pointer(cmd)), uintptr(unsafe.Pointer(argc)), 0)
	argv = (*[8192]*[8192]uint16)(unsafe.Pointer(r0))
	if argv == nil {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func LocalFree(hmem Handle) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procLocalFree.Addr(), 1, uintptr(hmem), 0, 0)
	handle = Handle(r0)
	if handle != 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetHandleInformation(handle Handle, mask uint32, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procSetHandleInformation.Addr(), 3, uintptr(handle), uintptr(mask), uintptr(flags))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FlushFileBuffers(handle Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procFlushFileBuffers.Addr(), 1, uintptr(handle), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetFullPathName(path *uint16, buflen uint32, buf *uint16, fname **uint16) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall6(procGetFullPathNameW.Addr(), 4, uintptr(unsafe.Pointer(path)), uintptr(buflen), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(fname)), 0, 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetLongPathName(path *uint16, buf *uint16, buflen uint32) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetLongPathNameW.Addr(), 3, uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(buf)), uintptr(buflen))
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetShortPathName(longpath *uint16, shortpath *uint16, buflen uint32) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetShortPathNameW.Addr(), 3, uintptr(unsafe.Pointer(longpath)), uintptr(unsafe.Pointer(shortpath)), uintptr(buflen))
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateFileMapping(fhandle Handle, sa *SecurityAttributes, prot uint32, maxSizeHigh uint32, maxSizeLow uint32, name *uint16) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procCreateFileMappingW.Addr(), 6, uintptr(fhandle), uintptr(unsafe.Pointer(sa)), uintptr(prot), uintptr(maxSizeHigh), uintptr(maxSizeLow), uintptr(unsafe.Pointer(name)))
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func MapViewOfFile(handle Handle, access uint32, offsetHigh uint32, offsetLow uint32, length uintptr) (addr uintptr, err error) {
	r0, _, e1 := syscall.Syscall6(procMapViewOfFile.Addr(), 5, uintptr(handle), uintptr(access), uintptr(offsetHigh), uintptr(offsetLow), uintptr(length), 0)
	addr = uintptr(r0)
	if addr == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func UnmapViewOfFile(addr uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procUnmapViewOfFile.Addr(), 1, uintptr(addr), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FlushViewOfFile(addr uintptr, length uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procFlushViewOfFile.Addr(), 2, uintptr(addr), uintptr(length), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func VirtualLock(addr uintptr, length uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procVirtualLock.Addr(), 2, uintptr(addr), uintptr(length), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func VirtualUnlock(addr uintptr, length uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procVirtualUnlock.Addr(), 2, uintptr(addr), uintptr(length), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func VirtualAlloc(address uintptr, size uintptr, alloctype uint32, protect uint32) (value uintptr, err error) {
	r0, _, e1 := syscall.Syscall6(procVirtualAlloc.Addr(), 4, uintptr(address), uintptr(size), uintptr(alloctype), uintptr(protect), 0, 0)
	value = uintptr(r0)
	if value == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func VirtualFree(address uintptr, size uintptr, freetype uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procVirtualFree.Addr(), 3, uintptr(address), uintptr(size), uintptr(freetype))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func VirtualProtect(address uintptr, size uintptr, newprotect uint32, oldprotect *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procVirtualProtect.Addr(), 4, uintptr(address), uintptr(size), uintptr(newprotect), uintptr(unsafe.Pointer(oldprotect)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func TransmitFile(s Handle, handle Handle, bytesToWrite uint32, bytsPerSend uint32, overlapped *Overlapped, transmitFileBuf *TransmitFileBuffers, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall9(procTransmitFile.Addr(), 7, uintptr(s), uintptr(handle), uintptr(bytesToWrite), uintptr(bytsPerSend), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(transmitFileBuf)), uintptr(flags), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ReadDirectoryChanges(handle Handle, buf *byte, buflen uint32, watchSubTree bool, mask uint32, retlen *uint32, overlapped *Overlapped, completionRoutine uintptr) (err error) {
	var _p0 uint32
	if watchSubTree {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := syscall.Syscall9(procReadDirectoryChangesW.Addr(), 8, uintptr(handle), uintptr(unsafe.Pointer(buf)), uintptr(buflen), uintptr(_p0), uintptr(mask), uintptr(unsafe.Pointer(retlen)), uintptr(unsafe.Pointer(overlapped)), uintptr(completionRoutine), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertOpenSystemStore(hprov Handle, name *uint16) (store Handle, err error) {
	r0, _, e1 := syscall.Syscall(procCertOpenSystemStoreW.Addr(), 2, uintptr(hprov), uintptr(unsafe.Pointer(name)), 0)
	store = Handle(r0)
	if store == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertOpenStore(storeProvider uintptr, msgAndCertEncodingType uint32, cryptProv uintptr, flags uint32, para uintptr) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procCertOpenStore.Addr(), 5, uintptr(storeProvider), uintptr(msgAndCertEncodingType), uintptr(cryptProv), uintptr(flags), uintptr(para), 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertEnumCertificatesInStore(store Handle, prevContext *CertContext) (context *CertContext, err error) {
	r0, _, e1 := syscall.Syscall(procCertEnumCertificatesInStore.Addr(), 2, uintptr(store), uintptr(unsafe.Pointer(prevContext)), 0)
	context = (*CertContext)(unsafe.Pointer(r0))
	if context == nil {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertAddCertificateContextToStore(store Handle, certContext *CertContext, addDisposition uint32, storeContext **CertContext) (err error) {
	r1, _, e1 := syscall.Syscall6(procCertAddCertificateContextToStore.Addr(), 4, uintptr(store), uintptr(unsafe.Pointer(certContext)), uintptr(addDisposition), uintptr(unsafe.Pointer(storeContext)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertCloseStore(store Handle, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procCertCloseStore.Addr(), 2, uintptr(store), uintptr(flags), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertGetCertificateChain(engine Handle, leaf *CertContext, time *Filetime, additionalStore Handle, para *CertChainPara, flags uint32, reserved uintptr, chainCtx **CertChainContext) (err error) {
	r1, _, e1 := syscall.Syscall9(procCertGetCertificateChain.Addr(), 8, uintptr(engine), uintptr(unsafe.Pointer(leaf)), uintptr(unsafe.Pointer(time)), uintptr(additionalStore), uintptr(unsafe.Pointer(para)), uintptr(flags), uintptr(reserved), uintptr(unsafe.Pointer(chainCtx)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertFreeCertificateChain(ctx *CertChainContext) {
	syscall.Syscall(procCertFreeCertificateChain.Addr(), 1, uintptr(unsafe.Pointer(ctx)), 0, 0)
	return
}

func CertCreateCertificateContext(certEncodingType uint32, certEncoded *byte, encodedLen uint32) (context *CertContext, err error) {
	r0, _, e1 := syscall.Syscall(procCertCreateCertificateContext.Addr(), 3, uintptr(certEncodingType), uintptr(unsafe.Pointer(certEncoded)), uintptr(encodedLen))
	context = (*CertContext)(unsafe.Pointer(r0))
	if context == nil {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertFreeCertificateContext(ctx *CertContext) (err error) {
	r1, _, e1 := syscall.Syscall(procCertFreeCertificateContext.Addr(), 1, uintptr(unsafe.Pointer(ctx)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CertVerifyCertificateChainPolicy(policyOID uintptr, chain *CertChainContext, para *CertChainPolicyPara, status *CertChainPolicyStatus) (err error) {
	r1, _, e1 := syscall.Syscall6(procCertVerifyCertificateChainPolicy.Addr(), 4, uintptr(policyOID), uintptr(unsafe.Pointer(chain)), uintptr(unsafe.Pointer(para)), uintptr(unsafe.Pointer(status)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func RegOpenKeyEx(key Handle, subkey *uint16, options uint32, desiredAccess uint32, result *Handle) (regerrno error) {
	r0, _, _ := syscall.Syscall6(procRegOpenKeyExW.Addr(), 5, uintptr(key), uintptr(unsafe.Pointer(subkey)), uintptr(options), uintptr(desiredAccess), uintptr(unsafe.Pointer(result)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func RegCloseKey(key Handle) (regerrno error) {
	r0, _, _ := syscall.Syscall(procRegCloseKey.Addr(), 1, uintptr(key), 0, 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func RegQueryInfoKey(key Handle, class *uint16, classLen *uint32, reserved *uint32, subkeysLen *uint32, maxSubkeyLen *uint32, maxClassLen *uint32, valuesLen *uint32, maxValueNameLen *uint32, maxValueLen *uint32, saLen *uint32, lastWriteTime *Filetime) (regerrno error) {
	r0, _, _ := syscall.Syscall12(procRegQueryInfoKeyW.Addr(), 12, uintptr(key), uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(classLen)), uintptr(unsafe.Pointer(reserved)), uintptr(unsafe.Pointer(subkeysLen)), uintptr(unsafe.Pointer(maxSubkeyLen)), uintptr(unsafe.Pointer(maxClassLen)), uintptr(unsafe.Pointer(valuesLen)), uintptr(unsafe.Pointer(maxValueNameLen)), uintptr(unsafe.Pointer(maxValueLen)), uintptr(unsafe.Pointer(saLen)), uintptr(unsafe.Pointer(lastWriteTime)))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func RegEnumKeyEx(key Handle, index uint32, name *uint16, nameLen *uint32, reserved *uint32, class *uint16, classLen *uint32, lastWriteTime *Filetime) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegEnumKeyExW.Addr(), 8, uintptr(key), uintptr(index), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(nameLen)), uintptr(unsafe.Pointer(reserved)), uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(classLen)), uintptr(unsafe.Pointer(lastWriteTime)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func RegQueryValueEx(key Handle, name *uint16, reserved *uint32, valtype *uint32, buf *byte, buflen *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall6(procRegQueryValueExW.Addr(), 6, uintptr(key), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(reserved)), uintptr(unsafe.Pointer(valtype)), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(buflen)))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func getCurrentProcessId() (pid uint32) {
	r0, _, _ := syscall.Syscall(procGetCurrentProcessId.Addr(), 0, 0, 0, 0)
	pid = uint32(r0)
	return
}

func GetConsoleMode(console Handle, mode *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(console), uintptr(unsafe.Pointer(mode)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetConsoleMode(console Handle, mode uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(console), uintptr(mode), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetConsoleScreenBufferInfo(console Handle, info *ConsoleScreenBufferInfo) (err error) {
	r1, _, e1 := syscall.Syscall(procGetConsoleScreenBufferInfo.Addr(), 2, uintptr(console), uintptr(unsafe.Pointer(info)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WriteConsole(console Handle, buf *uint16, towrite uint32, written *uint32, reserved *byte) (err error) {
	r1, _, e1 := syscall.Syscall6(procWriteConsoleW.Addr(), 5, uintptr(console), uintptr(unsafe.Pointer(buf)), uintptr(towrite), uintptr(unsafe.Pointer(written)), uintptr(unsafe.Pointer(reserved)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ReadConsole(console Handle, buf *uint16, toread uint32, read *uint32, inputControl *byte) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadConsoleW.Addr(), 5, uintptr(console), uintptr(unsafe.Pointer(buf)), uintptr(toread), uintptr(unsafe.Pointer(read)), uintptr(unsafe.Pointer(inputControl)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateToolhelp32Snapshot(flags uint32, processId uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procCreateToolhelp32Snapshot.Addr(), 2, uintptr(flags), uintptr(processId), 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Process32First(snapshot Handle, procEntry *ProcessEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procProcess32FirstW.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(procEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Process32Next(snapshot Handle, procEntry *ProcessEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procProcess32NextW.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(procEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DeviceIoControl(handle Handle, ioControlCode uint32, inBuffer *byte, inBufferSize uint32, outBuffer *byte, outBufferSize uint32, bytesReturned *uint32, overlapped *Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall9(procDeviceIoControl.Addr(), 8, uintptr(handle), uintptr(ioControlCode), uintptr(unsafe.Pointer(inBuffer)), uintptr(inBufferSize), uintptr(unsafe.Pointer(outBuffer)), uintptr(outBufferSize), uintptr(unsafe.Pointer(bytesReturned)), uintptr(unsafe.Pointer(overlapped)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateSymbolicLink(symlinkfilename *uint16, targetfilename *uint16, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procCreateSymbolicLinkW.Addr(), 3, uintptr(unsafe.Pointer(symlinkfilename)), uintptr(unsafe.Pointer(targetfilename)), uintptr(flags))
	if r1&0xff == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateHardLink(filename *uint16, existingfilename *uint16, reserved uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procCreateHardLinkW.Addr(), 3, uintptr(unsafe.Pointer(filename)), uintptr(unsafe.Pointer(existingfilename)), uintptr(reserved))
	if r1&0xff == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetCurrentThreadId() (id uint32) {
	r0, _, _ := syscall.Syscall(procGetCurrentThreadId.Addr(), 0, 0, 0, 0)
	id = uint32(r0)
	return
}

func CreateEvent(eventAttrs *SecurityAttributes, manualReset uint32, initialState uint32, name *uint16) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procCreateEventW.Addr(), 4, uintptr(unsafe.Pointer(eventAttrs)), uintptr(manualReset), uintptr(initialState), uintptr(unsafe.Pointer(name)), 0, 0)
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateEventEx(eventAttrs *SecurityAttributes, name *uint16, flags uint32, desiredAccess uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procCreateEventExW.Addr(), 4, uintptr(unsafe.Pointer(eventAttrs)), uintptr(unsafe.Pointer(name)), uintptr(flags), uintptr(desiredAccess), 0, 0)
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenEvent(desiredAccess uint32, inheritHandle bool, name *uint16) (handle Handle, err error) {
	var _p0 uint32
	if inheritHandle {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := syscall.Syscall(procOpenEventW.Addr(), 3, uintptr(desiredAccess), uintptr(_p0), uintptr(unsafe.Pointer(name)))
	handle = Handle(r0)
	if handle == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetEvent(event Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procSetEvent.Addr(), 1, uintptr(event), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ResetEvent(event Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procResetEvent.Addr(), 1, uintptr(event), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func PulseEvent(event Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procPulseEvent.Addr(), 1, uintptr(event), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DefineDosDevice(flags uint32, deviceName *uint16, targetPath *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procDefineDosDeviceW.Addr(), 3, uintptr(flags), uintptr(unsafe.Pointer(deviceName)), uintptr(unsafe.Pointer(targetPath)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DeleteVolumeMountPoint(volumeMountPoint *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procDeleteVolumeMountPointW.Addr(), 1, uintptr(unsafe.Pointer(volumeMountPoint)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindFirstVolume(volumeName *uint16, bufferLength uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procFindFirstVolumeW.Addr(), 2, uintptr(unsafe.Pointer(volumeName)), uintptr(bufferLength), 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindFirstVolumeMountPoint(rootPathName *uint16, volumeMountPoint *uint16, bufferLength uint32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procFindFirstVolumeMountPointW.Addr(), 3, uintptr(unsafe.Pointer(rootPathName)), uintptr(unsafe.Pointer(volumeMountPoint)), uintptr(bufferLength))
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindNextVolume(findVolume Handle, volumeName *uint16, bufferLength uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procFindNextVolumeW.Addr(), 3, uintptr(findVolume), uintptr(unsafe.Pointer(volumeName)), uintptr(bufferLength))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindNextVolumeMountPoint(findVolumeMountPoint Handle, volumeMountPoint *uint16, bufferLength uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procFindNextVolumeMountPointW.Addr(), 3, uintptr(findVolumeMountPoint), uintptr(unsafe.Pointer(volumeMountPoint)), uintptr(bufferLength))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindVolumeClose(findVolume Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procFindVolumeClose.Addr(), 1, uintptr(findVolume), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindVolumeMountPointClose(findVolumeMountPoint Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procFindVolumeMountPointClose.Addr(), 1, uintptr(findVolumeMountPoint), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetDriveType(rootPathName *uint16) (driveType uint32) {
	r0, _, _ := syscall.Syscall(procGetDriveTypeW.Addr(), 1, uintptr(unsafe.Pointer(rootPathName)), 0, 0)
	driveType = uint32(r0)
	return
}

func GetLogicalDrives() (drivesBitMask uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetLogicalDrives.Addr(), 0, 0, 0, 0)
	drivesBitMask = uint32(r0)
	if drivesBitMask == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetLogicalDriveStrings(bufferLength uint32, buffer *uint16) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetLogicalDriveStringsW.Addr(), 2, uintptr(bufferLength), uintptr(unsafe.Pointer(buffer)), 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetVolumeInformation(rootPathName *uint16, volumeNameBuffer *uint16, volumeNameSize uint32, volumeNameSerialNumber *uint32, maximumComponentLength *uint32, fileSystemFlags *uint32, fileSystemNameBuffer *uint16, fileSystemNameSize uint32) (err error) {
	r1, _, e1 := syscall.Syscall9(procGetVolumeInformationW.Addr(), 8, uintptr(unsafe.Pointer(rootPathName)), uintptr(unsafe.Pointer(volumeNameBuffer)), uintptr(volumeNameSize), uintptr(unsafe.Pointer(volumeNameSerialNumber)), uintptr(unsafe.Pointer(maximumComponentLength)), uintptr(unsafe.Pointer(fileSystemFlags)), uintptr(unsafe.Pointer(fileSystemNameBuffer)), uintptr(fileSystemNameSize), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetVolumeInformationByHandle(file Handle, volumeNameBuffer *uint16, volumeNameSize uint32, volumeNameSerialNumber *uint32, maximumComponentLength *uint32, fileSystemFlags *uint32, fileSystemNameBuffer *uint16, fileSystemNameSize uint32) (err error) {
	r1, _, e1 := syscall.Syscall9(procGetVolumeInformationByHandleW.Addr(), 8, uintptr(file), uintptr(unsafe.Pointer(volumeNameBuffer)), uintptr(volumeNameSize), uintptr(unsafe.Pointer(volumeNameSerialNumber)), uintptr(unsafe.Pointer(maximumComponentLength)), uintptr(unsafe.Pointer(fileSystemFlags)), uintptr(unsafe.Pointer(fileSystemNameBuffer)), uintptr(fileSystemNameSize), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetVolumeNameForVolumeMountPoint(volumeMountPoint *uint16, volumeName *uint16, bufferlength uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetVolumeNameForVolumeMountPointW.Addr(), 3, uintptr(unsafe.Pointer(volumeMountPoint)), uintptr(unsafe.Pointer(volumeName)), uintptr(bufferlength))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetVolumePathName(fileName *uint16, volumePathName *uint16, bufferLength uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetVolumePathNameW.Addr(), 3, uintptr(unsafe.Pointer(fileName)), uintptr(unsafe.Pointer(volumePathName)), uintptr(bufferLength))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetVolumePathNamesForVolumeName(volumeName *uint16, volumePathNames *uint16, bufferLength uint32, returnLength *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procGetVolumePathNamesForVolumeNameW.Addr(), 4, uintptr(unsafe.Pointer(volumeName)), uintptr(unsafe.Pointer(volumePathNames)), uintptr(bufferLength), uintptr(unsafe.Pointer(returnLength)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func QueryDosDevice(deviceName *uint16, targetPath *uint16, max uint32) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procQueryDosDeviceW.Addr(), 3, uintptr(unsafe.Pointer(deviceName)), uintptr(unsafe.Pointer(targetPath)), uintptr(max))
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetVolumeLabel(rootPathName *uint16, volumeName *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procSetVolumeLabelW.Addr(), 2, uintptr(unsafe.Pointer(rootPathName)), uintptr(unsafe.Pointer(volumeName)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetVolumeMountPoint(volumeMountPoint *uint16, volumeName *uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procSetVolumeMountPointW.Addr(), 2, uintptr(unsafe.Pointer(volumeMountPoint)), uintptr(unsafe.Pointer(volumeName)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WSAStartup(verreq uint32, data *WSAData) (sockerr error) {
	r0, _, _ := syscall.Syscall(procWSAStartup.Addr(), 2, uintptr(verreq), uintptr(unsafe.Pointer(data)), 0)
	if r0 != 0 {
		sockerr = syscall.Errno(r0)
	}
	return
}

func WSACleanup() (err error) {
	r1, _, e1 := syscall.Syscall(procWSACleanup.Addr(), 0, 0, 0, 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WSAIoctl(s Handle, iocc uint32, inbuf *byte, cbif uint32, outbuf *byte, cbob uint32, cbbr *uint32, overlapped *Overlapped, completionRoutine uintptr) (err error) {
	r1, _, e1 := syscall.Syscall9(procWSAIoctl.Addr(), 9, uintptr(s), uintptr(iocc), uintptr(unsafe.Pointer(inbuf)), uintptr(cbif), uintptr(unsafe.Pointer(outbuf)), uintptr(cbob), uintptr(unsafe.Pointer(cbbr)), uintptr(unsafe.Pointer(overlapped)), uintptr(completionRoutine))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func socket(af int32, typ int32, protocol int32) (handle Handle, err error) {
	r0, _, e1 := syscall.Syscall(procsocket.Addr(), 3, uintptr(af), uintptr(typ), uintptr(protocol))
	handle = Handle(r0)
	if handle == InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Setsockopt(s Handle, level int32, optname int32, optval *byte, optlen int32) (err error) {
	r1, _, e1 := syscall.Syscall6(procsetsockopt.Addr(), 5, uintptr(s), uintptr(level), uintptr(optname), uintptr(unsafe.Pointer(optval)), uintptr(optlen), 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Getsockopt(s Handle, level int32, optname int32, optval *byte, optlen *int32) (err error) {
	r1, _, e1 := syscall.Syscall6(procgetsockopt.Addr(), 5, uintptr(s), uintptr(level), uintptr(optname), uintptr(unsafe.Pointer(optval)), uintptr(unsafe.Pointer(optlen)), 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func bind(s Handle, name unsafe.Pointer, namelen int32) (err error) {
	r1, _, e1 := syscall.Syscall(procbind.Addr(), 3, uintptr(s), uintptr(name), uintptr(namelen))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func connect(s Handle, name unsafe.Pointer, namelen int32) (err error) {
	r1, _, e1 := syscall.Syscall(procconnect.Addr(), 3, uintptr(s), uintptr(name), uintptr(namelen))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getsockname(s Handle, rsa *RawSockaddrAny, addrlen *int32) (err error) {
	r1, _, e1 := syscall.Syscall(procgetsockname.Addr(), 3, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getpeername(s Handle, rsa *RawSockaddrAny, addrlen *int32) (err error) {
	r1, _, e1 := syscall.Syscall(procgetpeername.Addr(), 3, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func listen(s Handle, backlog int32) (err error) {
	r1, _, e1 := syscall.Syscall(proclisten.Addr(), 2, uintptr(s), uintptr(backlog), 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func shutdown(s Handle, how int32) (err error) {
	r1, _, e1 := syscall.Syscall(procshutdown.Addr(), 2, uintptr(s), uintptr(how), 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Closesocket(s Handle) (err error) {
	r1, _, e1 := syscall.Syscall(procclosesocket.Addr(), 1, uintptr(s), 0, 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func AcceptEx(ls Handle, as Handle, buf *byte, rxdatalen uint32, laddrlen uint32, raddrlen uint32, recvd *uint32, overlapped *Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall9(procAcceptEx.Addr(), 8, uintptr(ls), uintptr(as), uintptr(unsafe.Pointer(buf)), uintptr(rxdatalen), uintptr(laddrlen), uintptr(raddrlen), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(overlapped)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetAcceptExSockaddrs(buf *byte, rxdatalen uint32, laddrlen uint32, raddrlen uint32, lrsa **RawSockaddrAny, lrsalen *int32, rrsa **RawSockaddrAny, rrsalen *int32) {
	syscall.Syscall9(procGetAcceptExSockaddrs.Addr(), 8, uintptr(unsafe.Pointer(buf)), uintptr(rxdatalen), uintptr(laddrlen), uintptr(raddrlen), uintptr(unsafe.Pointer(lrsa)), uintptr(unsafe.Pointer(lrsalen)), uintptr(unsafe.Pointer(rrsa)), uintptr(unsafe.Pointer(rrsalen)), 0)
	return
}

func WSARecv(s Handle, bufs *WSABuf, bufcnt uint32, recvd *uint32, flags *uint32, overlapped *Overlapped, croutine *byte) (err error) {
	r1, _, e1 := syscall.Syscall9(procWSARecv.Addr(), 7, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(flags)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WSASend(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, overlapped *Overlapped, croutine *byte) (err error) {
	r1, _, e1 := syscall.Syscall9(procWSASend.Addr(), 7, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(sent)), uintptr(flags), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WSARecvFrom(s Handle, bufs *WSABuf, bufcnt uint32, recvd *uint32, flags *uint32, from *RawSockaddrAny, fromlen *int32, overlapped *Overlapped, croutine *byte) (err error) {
	r1, _, e1 := syscall.Syscall9(procWSARecvFrom.Addr(), 9, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(flags)), uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(fromlen)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WSASendTo(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, to *RawSockaddrAny, tolen int32, overlapped *Overlapped, croutine *byte) (err error) {
	r1, _, e1 := syscall.Syscall9(procWSASendTo.Addr(), 9, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(sent)), uintptr(flags), uintptr(unsafe.Pointer(to)), uintptr(tolen), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetHostByName(name string) (h *Hostent, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(name)
	if err != nil {
		return
	}
	return _GetHostByName(_p0)
}

func _GetHostByName(name *byte) (h *Hostent, err error) {
	r0, _, e1 := syscall.Syscall(procgethostbyname.Addr(), 1, uintptr(unsafe.Pointer(name)), 0, 0)
	h = (*Hostent)(unsafe.Pointer(r0))
	if h == nil {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetServByName(name string, proto string) (s *Servent, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(name)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = syscall.BytePtrFromString(proto)
	if err != nil {
		return
	}
	return _GetServByName(_p0, _p1)
}

func _GetServByName(name *byte, proto *byte) (s *Servent, err error) {
	r0, _, e1 := syscall.Syscall(procgetservbyname.Addr(), 2, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(proto)), 0)
	s = (*Servent)(unsafe.Pointer(r0))
	if s == nil {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Ntohs(netshort uint16) (u uint16) {
	r0, _, _ := syscall.Syscall(procntohs.Addr(), 1, uintptr(netshort), 0, 0)
	u = uint16(r0)
	return
}

func GetProtoByName(name string) (p *Protoent, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(name)
	if err != nil {
		return
	}
	return _GetProtoByName(_p0)
}

func _GetProtoByName(name *byte) (p *Protoent, err error) {
	r0, _, e1 := syscall.Syscall(procgetprotobyname.Addr(), 1, uintptr(unsafe.Pointer(name)), 0, 0)
	p = (*Protoent)(unsafe.Pointer(r0))
	if p == nil {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DnsQuery(name string, qtype uint16, options uint32, extra *byte, qrs **DNSRecord, pr *byte) (status error) {
	var _p0 *uint16
	_p0, status = syscall.UTF16PtrFromString(name)
	if status != nil {
		return
	}
	return _DnsQuery(_p0, qtype, options, extra, qrs, pr)
}

func _DnsQuery(name *uint16, qtype uint16, options uint32, extra *byte, qrs **DNSRecord, pr *byte) (status error) {
	r0, _, _ := syscall.Syscall6(procDnsQuery_W.Addr(), 6, uintptr(unsafe.Pointer(name)), uintptr(qtype), uintptr(options), uintptr(unsafe.Pointer(extra)), uintptr(unsafe.Pointer(qrs)), uintptr(unsafe.Pointer(pr)))
	if r0 != 0 {
		status = syscall.Errno(r0)
	}
	return
}

func DnsRecordListFree(rl *DNSRecord, freetype uint32) {
	syscall.Syscall(procDnsRecordListFree.Addr(), 2, uintptr(unsafe.Pointer(rl)), uintptr(freetype), 0)
	return
}

func DnsNameCompare(name1 *uint16, name2 *uint16) (same bool) {
	r0, _, _ := syscall.Syscall(procDnsNameCompare_W.Addr(), 2, uintptr(unsafe.Pointer(name1)), uintptr(unsafe.Pointer(name2)), 0)
	same = r0 != 0
	return
}

func GetAddrInfoW(nodename *uint16, servicename *uint16, hints *AddrinfoW, result **AddrinfoW) (sockerr error) {
	r0, _, _ := syscall.Syscall6(procGetAddrInfoW.Addr(), 4, uintptr(unsafe.Pointer(nodename)), uintptr(unsafe.Pointer(servicename)), uintptr(unsafe.Pointer(hints)), uintptr(unsafe.Pointer(result)), 0, 0)
	if r0 != 0 {
		sockerr = syscall.Errno(r0)
	}
	return
}

func FreeAddrInfoW(addrinfo *AddrinfoW) {
	syscall.Syscall(procFreeAddrInfoW.Addr(), 1, uintptr(unsafe.Pointer(addrinfo)), 0, 0)
	return
}

func GetIfEntry(pIfRow *MibIfRow) (errcode error) {
	r0, _, _ := syscall.Syscall(procGetIfEntry.Addr(), 1, uintptr(unsafe.Pointer(pIfRow)), 0, 0)
	if r0 != 0 {
		errcode = syscall.Errno(r0)
	}
	return
}

func GetAdaptersInfo(ai *IpAdapterInfo, ol *uint32) (errcode error) {
	r0, _, _ := syscall.Syscall(procGetAdaptersInfo.Addr(), 2, uintptr(unsafe.Pointer(ai)), uintptr(unsafe.Pointer(ol)), 0)
	if r0 != 0 {
		errcode = syscall.Errno(r0)
	}
	return
}

func SetFileCompletionNotificationModes(handle Handle, flags uint8) (err error) {
	r1, _, e1 := syscall.Syscall(procSetFileCompletionNotificationModes.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WSAEnumProtocols(protocols *int32, protocolBuffer *WSAProtocolInfo, bufferLength *uint32) (n int32, err error) {
	r0, _, e1 := syscall.Syscall(procWSAEnumProtocolsW.Addr(), 3, uintptr(unsafe.Pointer(protocols)), uintptr(unsafe.Pointer(protocolBuffer)), uintptr(unsafe.Pointer(bufferLength)))
	n = int32(r0)
	if n == -1 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetAdaptersAddresses(family uint32, flags uint32, reserved uintptr, adapterAddresses *IpAdapterAddresses, sizePointer *uint32) (errcode error) {
	r0, _, _ := syscall.Syscall6(procGetAdaptersAddresses.Addr(), 5, uintptr(family), uintptr(flags), uintptr(reserved), uintptr(unsafe.Pointer(adapterAddresses)), uintptr(unsafe.Pointer(sizePointer)), 0)
	if r0 != 0 {
		errcode = syscall.Errno(r0)
	}
	return
}

func GetACP() (acp uint32) {
	r0, _, _ := syscall.Syscall(procGetACP.Addr(), 0, 0, 0, 0)
	acp = uint32(r0)
	return
}

func MultiByteToWideChar(codePage uint32, dwFlags uint32, str *byte, nstr int32, wchar *uint16, nwchar int32) (nwrite int32, err error) {
	r0, _, e1 := syscall.Syscall6(procMultiByteToWideChar.Addr(), 6, uintptr(codePage), uintptr(dwFlags), uintptr(unsafe.Pointer(str)), uintptr(nstr), uintptr(unsafe.Pointer(wchar)), uintptr(nwchar))
	nwrite = int32(r0)
	if nwrite == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func TranslateName(accName *uint16, accNameFormat uint32, desiredNameFormat uint32, translatedName *uint16, nSize *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procTranslateNameW.Addr(), 5, uintptr(unsafe.Pointer(accName)), uintptr(accNameFormat), uintptr(desiredNameFormat), uintptr(unsafe.Pointer(translatedName)), uintptr(unsafe.Pointer(nSize)), 0)
	if r1&0xff == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetUserNameEx(nameFormat uint32, nameBuffre *uint16, nSize *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetUserNameExW.Addr(), 3, uintptr(nameFormat), uintptr(unsafe.Pointer(nameBuffre)), uintptr(unsafe.Pointer(nSize)))
	if r1&0xff == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func NetUserGetInfo(serverName *uint16, userName *uint16, level uint32, buf **byte) (neterr error) {
	r0, _, _ := syscall.Syscall6(procNetUserGetInfo.Addr(), 4, uintptr(unsafe.Pointer(serverName)), uintptr(unsafe.Pointer(userName)), uintptr(level), uintptr(unsafe.Pointer(buf)), 0, 0)
	if r0 != 0 {
		neterr = syscall.Errno(r0)
	}
	return
}

func NetGetJoinInformation(server *uint16, name **uint16, bufType *uint32) (neterr error) {
	r0, _, _ := syscall.Syscall(procNetGetJoinInformation.Addr(), 3, uintptr(unsafe.Pointer(server)), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(bufType)))
	if r0 != 0 {
		neterr = syscall.Errno(r0)
	}
	return
}

func NetApiBufferFree(buf *byte) (neterr error) {
	r0, _, _ := syscall.Syscall(procNetApiBufferFree.Addr(), 1, uintptr(unsafe.Pointer(buf)), 0, 0)
	if r0 != 0 {
		neterr = syscall.Errno(r0)
	}
	return
}

func LookupAccountSid(systemName *uint16, sid *SID, name *uint16, nameLen *uint32, refdDomainName *uint16, refdDomainNameLen *uint32, use *uint32) (err error) {
	r1, _, e1 := syscall.Syscall9(procLookupAccountSidW.Addr(), 7, uintptr(unsafe.Pointer(systemName)), uintptr(unsafe.Pointer(sid)), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(nameLen)), uintptr(unsafe.Pointer(refdDomainName)), uintptr(unsafe.Pointer(refdDomainNameLen)), uintptr(unsafe.Pointer(use)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func LookupAccountName(systemName *uint16, accountName *uint16, sid *SID, sidLen *uint32, refdDomainName *uint16, refdDomainNameLen *uint32, use *uint32) (err error) {
	r1, _, e1 := syscall.Syscall9(procLookupAccountNameW.Addr(), 7, uintptr(unsafe.Pointer(systemName)), uintptr(unsafe.Pointer(accountName)), uintptr(unsafe.Pointer(sid)), uintptr(unsafe.Pointer(sidLen)), uintptr(unsafe.Pointer(refdDomainName)), uintptr(unsafe.Pointer(refdDomainNameLen)), uintptr(unsafe.Pointer(use)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ConvertSidToStringSid(sid *SID, stringSid **uint16) (err error) {
	r1, _, e1 := syscall.Syscall(procConvertSidToStringSidW.Addr(), 2, uintptr(unsafe.Pointer(sid)), uintptr(unsafe.Pointer(stringSid)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ConvertStringSidToSid(stringSid *uint16, sid **SID) (err error) {
	r1, _, e1 := syscall.Syscall(procConvertStringSidToSidW.Addr(), 2, uintptr(unsafe.Pointer(stringSid)), uintptr(unsafe.Pointer(sid)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetLengthSid(sid *SID) (len uint32) {
	r0, _, _ := syscall.Syscall(procGetLengthSid.Addr(), 1, uintptr(unsafe.Pointer(sid)), 0, 0)
	len = uint32(r0)
	return
}

func CopySid(destSidLen uint32, destSid *SID, srcSid *SID) (err error) {
	r1, _, e1 := syscall.Syscall(procCopySid.Addr(), 3, uintptr(destSidLen), uintptr(unsafe.Pointer(destSid)), uintptr(unsafe.Pointer(srcSid)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func AllocateAndInitializeSid(identAuth *SidIdentifierAuthority, subAuth byte, subAuth0 uint32, subAuth1 uint32, subAuth2 uint32, subAuth3 uint32, subAuth4 uint32, subAuth5 uint32, subAuth6 uint32, subAuth7 uint32, sid **SID) (err error) {
	r1, _, e1 := syscall.Syscall12(procAllocateAndInitializeSid.Addr(), 11, uintptr(unsafe.Pointer(identAuth)), uintptr(subAuth), uintptr(subAuth0), uintptr(subAuth1), uintptr(subAuth2), uintptr(subAuth3), uintptr(subAuth4), uintptr(subAuth5), uintptr(subAuth6), uintptr(subAuth7), uintptr(unsafe.Pointer(sid)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FreeSid(sid *SID) (err error) {
	r1, _, e1 := syscall.Syscall(procFreeSid.Addr(), 1, uintptr(unsafe.Pointer(sid)), 0, 0)
	if r1 != 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func EqualSid(sid1 *SID, sid2 *SID) (isEqual bool) {
	r0, _, _ := syscall.Syscall(procEqualSid.Addr(), 2, uintptr(unsafe.Pointer(sid1)), uintptr(unsafe.Pointer(sid2)), 0)
	isEqual = r0 != 0
	return
}

func checkTokenMembership(tokenHandle Token, sidToCheck *SID, isMember *int32) (err error) {
	r1, _, e1 := syscall.Syscall(procCheckTokenMembership.Addr(), 3, uintptr(tokenHandle), uintptr(unsafe.Pointer(sidToCheck)), uintptr(unsafe.Pointer(isMember)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func OpenProcessToken(h Handle, access uint32, token *Token) (err error) {
	r1, _, e1 := syscall.Syscall(procOpenProcessToken.Addr(), 3, uintptr(h), uintptr(access), uintptr(unsafe.Pointer(token)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetTokenInformation(t Token, infoClass uint32, info *byte, infoLen uint32, returnedLen *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procGetTokenInformation.Addr(), 5, uintptr(t), uintptr(infoClass), uintptr(unsafe.Pointer(info)), uintptr(infoLen), uintptr(unsafe.Pointer(returnedLen)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetUserProfileDirectory(t Token, dir *uint16, dirLen *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetUserProfileDirectoryW.Addr(), 3, uintptr(t), uintptr(unsafe.Pointer(dir)), uintptr(unsafe.Pointer(dirLen)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
