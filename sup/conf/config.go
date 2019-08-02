package conf

// LoggingConfig specifies all the parameters needed for logging
type LoggingConfig struct {
	Level string // default: INFO
	File  string // default stdout
}

type AnsibleCfg struct {
	BinPath string // default /usr/local/bin/ansible-playbook
	Workdir string // default /tmp
}

// Config encapsulates he configuration
type Config struct {
	Backend    string // default ansible
	LogCfg     LoggingConfig
	AnsibleCfg AnsibleCfg
}
