





package windows

func itoa(val int) string { 
	if val < 0 {
		return "-" + itoa(-val)
	}
	var buf [32]byte 
	i := len(buf) - 1
	for val >= 10 {
		buf[i] = byte(val%10 + '0')
		i--
		val /= 10
	}
	buf[i] = byte(val + '0')
	return string(buf[i:])
}
