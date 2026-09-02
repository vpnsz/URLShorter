package config

import (
	"flag"
)

type hostPortFlag struct {
	host string
}

func (f *hostPortFlag) FlagPresent() bool {
	return len(f.host) != 0
}

func (f *hostPortFlag) String() string {
	return f.host
}

func (f *hostPortFlag) Set(arg string) error {
	f.host = arg
	return nil
}

func ParseFlags(c *Config) {
	paramA := new(hostPortFlag)
	paramB := new(hostPortFlag)
	var fileFlag = ""
	flag.Var(paramA, "a", "host:port")
	flag.Var(paramB, "b", "host:port")
	flag.StringVar(&fileFlag, "f", "./default_storage.txt", "-f file_path")
	flag.Parse()
	if paramA.FlagPresent() {
		c.ServerAddr = paramA.host
	}
	if paramB.FlagPresent() {
		c.BaseShorterAddr = paramB.host
	}
	if len(fileFlag) != 0 {
		c.StorageFilePath = fileFlag
	}
}
