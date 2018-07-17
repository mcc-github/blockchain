package system 


var lcowSupported = false


func InitLCOW(experimental bool) {
	v := GetOSVersion()
	if experimental && v.Build >= 16299 {
		lcowSupported = true
	}
}
