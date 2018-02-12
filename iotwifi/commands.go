package iotwifi

import (
	"os/exec"

	"github.com/bhoriuchi/go-bunyan/bunyan"
)

// Command for device network commands
type Command struct {
	Log      bunyan.Logger
	Runner   CmdRunner
}

// CheckInterface
func (c *Command) CheckInterface() {
	cmd := exec.Command("ifconfig", "uap0")
	go c.Runner.ProcessCmd("ifconfig_uap0", cmd)
}

// StartWpaSupplicant
func (c *Command) StartWpaSupplicant() {
	args := []string{
		"-d",
		"-Dnl80211",
		"-iwlan0",
		"-c/etc/wpa_supplicant/wpa_supplicant.conf",
	}
	
	cmd := exec.Command("wpa_supplicant", args...)
	go c.Runner.ProcessCmd("wpa_supplicant", cmd)
}

// StartDnsmasq
func (c *Command) StartDnsmasq() {
	// hostapd is enabled, fire up dnsmasq
	args := []string{
		"--no-hosts", // Don't read the hostnames in /etc/hosts.
		"--keep-in-foreground",
		"--log-queries",
		"--no-resolv",
		"--address=/#/192.168.27.1",
		"--dhcp-range=192.168.27.100,192.168.27.150,1h",
		"--dhcp-vendorclass=set:device,IoT",
		"--dhcp-authoritative",
		"--log-facility=-",
	}
	
	cmd := exec.Command("dnsmasq", args...)
	go c.Runner.ProcessCmd("dnsmasq", cmd)
}

// StartHostapd
func (c *Command) StartHostapd() {

	c.Runner.Log.Info("Starting hostapd.");
	
	cmd := exec.Command("hostapd", "-d", "/dev/stdin")
	hostapdPipe, _ := cmd.StdinPipe()
	c.Runner.ProcessCmd("hostapd", cmd)
	
	cfg := `interface=uap0
ssid=iotwifi2
hw_mode=g
channel=6
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase=iotwifipass
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP`
	
	hostapdPipe.Write([]byte(cfg))
	hostapdPipe.Close()
	
}
