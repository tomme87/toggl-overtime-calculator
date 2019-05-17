package config

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

type toggl struct {
	ApiToken    string `mapstructure:"api_token"`
	WorkspaceId int    `mapstructure:"workspace_id"`
}

type options struct {
	Since string `mapstructure:"since"`
	Until string `mapstructure:"until"`
}

type config struct {
	Toggl   toggl   `mapstructure:"toggl"`
	Options options `mapstructure:"options"`
}

var C config

func (c *config) Init() error {
	viper.SetConfigName("toggl_ot_calc")
	viper.AddConfigPath("/etc/toggl_ot_calc/")
	viper.AddConfigPath("$HOME/.toggl_ot_calc")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	flag.String("api-token", "", "Toggle API token (find it in toggl settings)")
	flag.Int("workspace-id", 0, "Your toggl Worksapace ID")
	flag.String("since", "", "Date since (yyyy-mm-dd)")
	flag.String("until", "", "Date until (yyyy-mm-dd)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	c.bindCmdFlag("toggl.api_token", "api-token")
	c.bindCmdFlag("toggl.workspace_id", "workspace-id")
	c.bindCmdFlag("options.since", "since")
	c.bindCmdFlag("options.until", "until")

	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	return nil
}

func (c *config) bindCmdFlag(cfgkey string, cmdkey string) {
	f := pflag.Lookup(cmdkey)
	if f == nil || f.Changed == false {
		return
	}

	err := viper.BindPFlag(cfgkey, f)
	if err != nil {
		log.Fatal(err)
	}
}
