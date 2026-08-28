package config

import (
	"flag"
	"fmt"
)

type hostPortFlag struct {
	host string
}

func (f *hostPortFlag) FlagPresent() bool {
	if len(f.host) != 0 {
		return true
	}
	return false
}

func (f *hostPortFlag) String() string {
	return fmt.Sprintf("%s", f.host)
}

func (f *hostPortFlag) Set(arg string) error {
	f.host = arg
	return nil
}

func ParseFlags(c *Config) {
	paramA := new(hostPortFlag)
	paramB := new(hostPortFlag)
	var fileFlag string = ""
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
