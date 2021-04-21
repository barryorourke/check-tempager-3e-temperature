package main

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	"net"
)

type Config struct {
	sensu.PluginConfig
	Target    string
	Community string
	Warning   float64
	Critical  float64
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-tempager-3e-temperature",
			Short:    "A simple temperature monitor check written in Go",
			Keyspace: "sensu.io/plugins/check-tempager-3e-temperature/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "target",
			Argument:  "target",
			Shorthand: "t",
			Default:   "",
			Usage:     "IP address of the target unit.",
			Value:     &plugin.Target,
		},
		{
			Path:      "community",
			Argument:  "community",
			Shorthand: "C",
			Default:   "public",
			Usage:     "SNMP community.",
			Value:     &plugin.Community,
		},
		{
			Path:      "warning",
			Argument:  "warning",
			Shorthand: "w",
			Default:   35.0,
			Usage:     "warning threshold.",
			Value:     &plugin.Warning,
		},
		{
			Path:      "critical",
			Argument:  "critical",
			Shorthand: "c",
			Default:   40.0,
			Usage:     "critical threshold.",
			Value:     &plugin.Critical,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {

	// target is a required argument
	if plugin.Target == "" {
		return sensu.CheckStateCritical, fmt.Errorf("target unit must be specified.")
	}

	// target must be an IP address
	if net.ParseIP(plugin.Target) == nil {
		return sensu.CheckStateCritical, fmt.Errorf("target must be an IP address.")
	}

	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {

	// configure the SNMP connection
	gosnmp.Default.Target = plugin.Target
	gosnmp.Default.Community = plugin.Community
	gosnmp.Default.Version = gosnmp.Version1

	// make the connection
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to connect to tempager.\n", plugin.PluginConfig.Name)
		return sensu.CheckStateCritical, nil
	}
	defer gosnmp.Default.Conn.Close()

	// gather the required values (internal sensor / external sensor)
	oids := []string{".1.3.6.1.2.1.1.6.0", ".1.3.6.1.4.1.20916.1.7.1.1.1.1.0", ".1.3.6.1.4.1.20916.1.7.1.2.1.1.0"}
	result, err := gosnmp.Default.Get(oids)
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to gather oids.\n", plugin.PluginConfig.Name)
		return sensu.CheckStateCritical, nil
	}

	// validate the location oid
	location_oid, ok := result.Variables[0].Value.([]uint8)
	if !ok {
		fmt.Printf("%s CRITICAL: failed to read location.\n", plugin.PluginConfig.Name)
		return sensu.CheckStateCritical, nil
	}

	// validate the internal temperature oid
	inttemp_oid, ok := result.Variables[1].Value.(int)
	if !ok {
		fmt.Printf("%s CRITICAL: failed to read internal temperature.\n", plugin.PluginConfig.Name)
		return sensu.CheckStateCritical, nil
	}

	// validate the external temperature oid
	exttemp_oid, ok := result.Variables[2].Value.(int)
	if !ok {
		fmt.Printf("%s CRITICAL: failed to read external temperature.\n", plugin.PluginConfig.Name)
		return sensu.CheckStateCritical, nil
	}

	// convert oid values into something usable
	location := string(location_oid)
	internal_temperature := float64(inttemp_oid) / 100.0
	external_temperature := float64(exttemp_oid) / 100.0

	// construct the performance data
	perfData := fmt.Sprintf("tempager_internal=%.2f, tempager_external=%.2f", internal_temperature, external_temperature)
	t := fmt.Sprintf("%s temperature is %.2fc | %s\n", location, external_temperature, perfData)

	if external_temperature > plugin.Critical {
		fmt.Printf("%s CRITICAL: %s", plugin.PluginConfig.Name, t)
		return sensu.CheckStateCritical, nil
	}

	if external_temperature > plugin.Warning {
		fmt.Printf("%s WARNING: %s", plugin.PluginConfig.Name, t)
		return sensu.CheckStateWarning, nil
	}

	fmt.Printf("%s OK: %s", plugin.PluginConfig.Name, t)
	return sensu.CheckStateOK, nil
}
