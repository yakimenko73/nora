package common

type LoadType string

const (
	LoadTypeUnknown   LoadType = ""
	LoadTypeRamp      LoadType = "ramp"
	LoadTypeDisturbed LoadType = "disturbed"
)
