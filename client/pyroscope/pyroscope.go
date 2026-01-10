package pyroscope

import "github.com/grafana/pyroscope-go"

func InitPyroscope(name, address string) {
	pyroscope.Start(pyroscope.Config{
		ApplicationName: name,
		ServerAddress:   address,
		Logger:          pyroscope.StandardLogger,

		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})
}
