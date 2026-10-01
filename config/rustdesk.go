package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	DefaultIdServerPort    = 21116
	DefaultRelayServerPort = 21117
)

type Rustdesk struct {
	IdServer        string `mapstructure:"id-server"`
	IdServerPort    int    `mapstructure:"-"`
	RelayServer     string `mapstructure:"relay-server"`
	RelayServerPort int    `mapstructure:"-"`
	ApiServer       string `mapstructure:"api-server"`
	Key             string `mapstructure:"key"`
	KeyFile         string `mapstructure:"key-file"`
	Personal        int    `mapstructure:"personal"`
	//webclient-magic-queryonline
	WebclientMagicQueryonline int    `mapstructure:"webclient-magic-queryonline"`
	WsHost                    string `mapstructure:"ws-host"`
	// XC: address-book auto sync (mirrors every peer into a shared collection)
	AutoAbSync       bool   `mapstructure:"auto-ab-sync"`
	AutoAbCollection string `mapstructure:"auto-ab-collection"`
}

// XC: Init fills the auto-sync defaults so a fresh deployment works without
// extra tuning: enabled unless explicitly disabled, default collection name.
func (rd *Rustdesk) Init(v *viper.Viper) {
	if v != nil && !v.IsSet("rustdesk.auto-ab-sync") {
		rd.AutoAbSync = true
	}
	if rd.AutoAbCollection == "" {
		rd.AutoAbCollection = "全部设备"
	}
}

func (rd *Rustdesk) LoadKeyFile() {
	// Load key file
	if rd.Key != "" {
		return
	}
	if rd.KeyFile != "" {
		// Load key from file
		b, err := os.ReadFile(rd.KeyFile)
		if err != nil {
			return
		}
		rd.Key = string(b)
		return
	}
}
