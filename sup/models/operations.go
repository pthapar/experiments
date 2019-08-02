package models

// reload:
//   srcTarballPath:
type DockerImageReloadOp struct {
	Reload struct {
		SrcTarballPath string `yaml:"srcTarballPath"`
	} `yaml:"reload"`
}

// restart:
//   serviceName: sherlock_configserver
type SystemDServiceRestartOp struct {
	Restart struct {
		ServiceName string `yaml:"serviceName"`
	} `yaml:"restart"`
}

// copy:
//   src:
//   dst:
type FileSystemCopyOp struct {
	Copy struct {
		Src string `yaml:"src"`
		Dst string `yaml:"dst"`
	} `yaml:"copy"`
}

// Operation ...
type Operation struct {
	Operation interface{} `yaml:"operation"`
	// TODO: add an execute method to each op
}
