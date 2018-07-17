package container 



func (n NetworkMode) IsBridge() bool {
	return n == "nat"
}



func (n NetworkMode) IsHost() bool {
	return false
}


func (n NetworkMode) IsUserDefined() bool {
	return !n.IsDefault() && !n.IsNone() && !n.IsBridge() && !n.IsContainer()
}


func (i Isolation) IsValid() bool {
	return i.IsDefault() || i.IsHyperV() || i.IsProcess()
}


func (n NetworkMode) NetworkName() string {
	if n.IsDefault() {
		return "default"
	} else if n.IsBridge() {
		return "nat"
	} else if n.IsNone() {
		return "none"
	} else if n.IsContainer() {
		return "container"
	} else if n.IsUserDefined() {
		return n.UserDefined()
	}

	return ""
}
