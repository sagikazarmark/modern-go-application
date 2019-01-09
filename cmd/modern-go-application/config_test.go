package main

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfigure(t *testing.T) {
	var config Config

	v := viper.New()
	p := pflag.NewFlagSet("test", pflag.ContinueOnError)

	Configure(v, p)

	file, err := os.Open("../../config.toml.dist")
	require.NoError(t, err)

	v.SetConfigType("toml")

	err = v.ReadConfig(file)
	require.NoError(t, err)

	err = v.Unmarshal(&config)
	require.NoError(t, err)

	err = config.Validate()
	require.NoError(t, err)
}
