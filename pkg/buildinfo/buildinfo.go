package buildinfo

func Version() string {
	return version
}

func Build() string {
	return build
}

func DevMode() bool {
	return devMode
}

func BuildTime() string {
	return buildTime
}
