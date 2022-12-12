package utils

type VersionInfo struct {
	BuildDate string
	Arch      string
	Version   string
}

func MakeVersionInfo(BuildDate string, Arch string, Version string) VersionInfo {
	return VersionInfo{BuildDate: BuildDate, Arch: Arch, Version: Version}
}

func (v *VersionInfo) GetInfo() string {
	return "nemo-" + v.Version + "-" + v.BuildDate + "-" + v.Arch
}
