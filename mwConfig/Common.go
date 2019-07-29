package mwConfig

import (
	DevLogs "github.com/a2si/MiniWeb/DevLogs"
)

type TConfig struct {
	RandUserAgent bool
	EnableLogs    bool
}

func NewConfig() *TConfig {
	Obj := &TConfig{
		RandUserAgent: true,
		EnableLogs:    true,
	}
	return Obj
}

func (self *TConfig) ReConfig() {
	DevLogs.LogsEnable(self.EnableLogs)
}
