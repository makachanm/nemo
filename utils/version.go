package utils

type VersionInfo struct {
	BuildDate string
	Arch      string
}

func MakeVersionInfo(BuildDate string, Arch string) VersionInfo {
	return VersionInfo{}
}

func (v *VersionInfo) GetInfo() string {
	return "nemo-" + v.BuildDate + "-" + v.Arch
}
