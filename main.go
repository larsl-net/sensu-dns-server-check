package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"crypto/tls"


	"github.com/miekg/dns"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	// udp, tcp, dot
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Port       int
	Protocol   string
	Record	   string
	Server     string
	ServerName string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-dns-server-check",
			Short:    "Check DNS Server functionality",
			Keyspace: "sensu.io/plugins/check-dns-server/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "port",
			Env:       "CHECK_PORT",
			Argument:  "port",
			Shorthand: "p",
			Default:   53,
			Usage:     "Port to check",
			Value:     &plugin.Port,
		},
		&sensu.PluginConfigOption{
			Path:      "protocol",
			Env:       "CHECK_PROTOCOL",
			Argument:  "protocol",
			Shorthand: "P",
			Default:   "udp",
			Usage:     "DNS Protocol to check (udp, tcp, dot)",
			Value:     &plugin.Protocol,
		},
		
		&sensu.PluginConfigOption{
			Path:      "record",
			Env:       "CHECK_RECORD",
			Argument:  "record",
			Shorthand: "r",
			Default:   "sensu.io",
			Usage:     "DNS Record to check",
			Value:     &plugin.Record,
		},
		&sensu.PluginConfigOption{
			Path:      "server",
			Env:       "CHECK_SERVER",
			Argument:  "server",
			Shorthand: "s",
			Default:   "",
			Usage:     "DNS Server to check",
			Value:     &plugin.Server,
		},
		&sensu.PluginConfigOption{
			Path:      "server-name",
			Env:       "CHECK_SERVER_NAME",
			Argument:  "server-name",
			Shorthand: "n",
			Default:   "",
			Usage:     "Hostname for DoT",
			Value:     &plugin.ServerName,
		},
	}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		log.Println("using stdin")
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.Server) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--server or CHECK_SERVER environment variable is required")
	}
	if net.ParseIP(plugin.Server) == nil {
		return sensu.CheckStateWarning, fmt.Errorf("no valid server IP")
	}
	if !isValidProtocol(plugin.Protocol) {
		return sensu.CheckStateWarning, fmt.Errorf("unknown Protocol")
	}
	if plugin.Protocol == "dot" && len(plugin.ServerName) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--server-name or CHECK_SERVER_NAME environment variable is required for protocol DoT")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	state := 0
	var err error
	if plugin.Protocol == "udp" {
		c := dns.Client{}
		state, err = checkDNS(plugin.Server, &c)
	} else if plugin.Protocol == "tcp" {
		c := dns.Client{
			Net: "tcp",
		}
		state, err = checkDNS(plugin.Server, &c)
	} else if plugin.Protocol == "dot" {
		c := dns.Client{
			Net: "tcp-tls",
			TLSConfig: &tls.Config{
				ServerName: plugin.ServerName,
			},
		}
		state, err = checkDNS(plugin.Server, &c)
	}
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to run check, error: %v\n", plugin.PluginConfig.Name, err)
		return sensu.CheckStateCritical, nil
	}
	return state, nil
}

func isValidProtocol(protocol string) bool {
	switch protocol {
	case
		"udp",
		"tcp",
		"dot":
		return true
	}
	return false
}

func checkDNS(server string, c *dns.Client) (int, error) {
	ip := net.ParseIP(server)
	if ip.To16() != nil {
		server = "[" + server + "]"
	}
	Records := []net.IP{}
	
	m := dns.Msg{}
	// A
	m.SetQuestion(plugin.Record+".", dns.TypeA)
	r, _, err := c.Exchange(&m, server+":"+strconv.Itoa(plugin.Port))

	if err != nil {
		return sensu.CheckStateCritical, err
	}
	if len(r.Answer) != 0 {
		for _, ans := range r.Answer {
			record := ans.(*dns.A)
			Records = append(Records, record.A)
		}
	}

	// AAAA
	m.SetQuestion(plugin.Record+".", dns.TypeAAAA)
	r, _, err = c.Exchange(&m, server+":"+strconv.Itoa(plugin.Port))

	if err != nil {
		return sensu.CheckStateCritical, err
	}
	if len(r.Answer) != 0 {
		for _, ans := range r.Answer {
			record := ans.(*dns.AAAA)
			Records = append(Records, record.AAAA)
		}
	}
	if len(Records) == 0 {
		fmt.Printf("%s WARNING: no Records found for %s\n", plugin.PluginConfig.Name, plugin.Record)
		return sensu.CheckStateWarning, nil
	}
	
	fmt.Printf("%s OK: %s returns %s\n", plugin.PluginConfig.Name, plugin.Record, Records)
	return sensu.CheckStateOK, nil
}
