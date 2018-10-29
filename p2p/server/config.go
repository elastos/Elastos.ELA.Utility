package server

import (
	"net"
	"strconv"
	"time"

	"github.com/elastos/Elastos.ELA.Utility/p2p"
)

const (
	defaultMaxPeers              = 125
	defaultBanThreshold   uint32 = 100
	defaultBanDuration           = time.Hour * 24
	defaultConnectTimeout        = time.Second * 30
)

// Config is a descriptor which specifies the server instance configuration.
type Config struct {
	MagicNumber      uint32
	ProtocolVersion  uint32
	Services         uint64
	SeedPeers        []string
	ListenAddrs      []string
	ExternalIPs      []string
	Upnp             bool
	DefaultPort      uint16
	DisableListen    bool
	DisableRelayTx   bool
	MaxPeers         int
	DisableBanning   bool
	BanThreshold     uint32
	BanDuration      time.Duration
	Whitelists       []*net.IPNet
	TargetOutbound   int
	OnNewPeer        func(IPeer)
	OnDonePeer       func(IPeer)
	MakeEmptyMessage func(string) (p2p.Message, error)
	BestHeight       func() uint64
}

func (cfg *Config) normalize() {
	defaultPort := strconv.FormatUint(uint64(cfg.DefaultPort), 10)

	// Add default port to all seed peer addresses if needed and remove
	// duplicate addresses.
	cfg.SeedPeers = normalizeAddresses(cfg.SeedPeers, defaultPort)

	// Add default port to all listener addresses if needed and remove
	// duplicate addresses.
	cfg.ListenAddrs = normalizeAddresses(cfg.ListenAddrs, defaultPort)
}

// inWhitelist returns whether the IP address is included in the whitelisted
// networks and IPs.
func (cfg *Config) inWhitelist(addr net.Addr) bool {
	if len(cfg.Whitelists) == 0 {
		return false
	}

	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		log.Warnf("Unable to SplitHostPort on '%s': %v", addr, err)
		return false
	}
	ip := net.ParseIP(host)
	if ip == nil {
		log.Warnf("Unable to parse IP '%s'", addr)
		return false
	}

	for _, ipnet := range cfg.Whitelists {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

func dialTimeout(addr net.Addr) (net.Conn, error) {
	return net.DialTimeout(addr.Network(), addr.String(), defaultConnectTimeout)
}

// removeDuplicateAddresses returns a new slice with all duplicate entries in
// addrs removed.
func removeDuplicateAddresses(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	seen := map[string]struct{}{}
	for _, val := range addrs {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = struct{}{}
		}
	}
	return result
}

// normalizeAddress returns addr with the passed default port appended if
// there is not already a port specified.
func normalizeAddress(addr, defaultPort string) string {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return net.JoinHostPort(addr, defaultPort)
	}
	return addr
}

// normalizeAddresses returns a new slice with all the passed peer addresses
// normalized with the given default port, and all duplicates removed.
func normalizeAddresses(addrs []string, defaultPort string) []string {
	for i, addr := range addrs {
		addrs[i] = normalizeAddress(addr, defaultPort)
	}

	return removeDuplicateAddresses(addrs)
}

// NewDefaultConfig returns a new config instance filled by default settings
// for the server.
func NewDefaultConfig(
	magic, pver uint32,
	services uint64,
	defaultPort uint16,
	seeds, listenAddrs []string,
	onNewPeer func(IPeer),
	onDonePeer func(IPeer),
	makeEmptyMessage func(string) (p2p.Message, error),
	bestHeight func() uint64) *Config {
	return &Config{
		MagicNumber:      magic,
		ProtocolVersion:  pver,
		Services:         services,
		SeedPeers:        seeds,
		ListenAddrs:      listenAddrs,
		ExternalIPs:      nil,
		Upnp:             false,
		DefaultPort:      defaultPort,
		DisableListen:    false,
		DisableRelayTx:   false,
		MaxPeers:         defaultMaxPeers,
		DisableBanning:   false,
		BanThreshold:     defaultBanThreshold,
		BanDuration:      defaultBanDuration,
		Whitelists:       nil,
		TargetOutbound:   defaultTargetOutbound,
		OnNewPeer:        onNewPeer,
		OnDonePeer:       onDonePeer,
		MakeEmptyMessage: makeEmptyMessage,
		BestHeight:       bestHeight,
	}
}
