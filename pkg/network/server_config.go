package network

import (
	"fmt"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/config/netmode"
	"go.uber.org/zap/zapcore"
)

type (
	// ServerConfig holds the server configuration.
	ServerConfig struct {
		// MinPeers is the minimum number of peers for normal operation.
		// When a node has less than this number of peers, it tries to
		// connect with some new ones.
		MinPeers int

		// AttemptConnPeers is the number of connection to try to
		// establish when the connection count drops below the MinPeers
		// value.
		AttemptConnPeers int

		// MaxPeers is the maximum number of peers that can
		// be connected to the server.
		MaxPeers int

		// The user agent of the server.
		UserAgent string

		// Addresses stores the list of bind addresses for the node.
		Addresses []config.AnnounceableAddress

		// The network mode the server will operate on.
		// ModePrivNet docker private network.
		// ModeTestNet NEO test network.
		// ModeMainNet NEO main network.
		Net netmode.Magic

		// Relay determines whether the server is forwarding its inventory.
		Relay bool

		// Seeds is a list of initial nodes used to establish connectivity.
		Seeds []string

		// Maximum duration a single dial may take.
		DialTimeout time.Duration

		// The duration between protocol ticks with each connected peer.
		// When this is 0, the default interval of 5 seconds will be used.
		ProtoTickInterval time.Duration

		// Interval used in pinging mechanism for syncing blocks.
		PingInterval time.Duration
		// Time to wait for pong(response for sent ping request).
		PingTimeout time.Duration

		// Level of the internal logger.
		LogLevel zapcore.Level

		// TimePerBlock is an interval which should pass between two successive blocks.
		TimePerBlock time.Duration

		// OracleCfg is oracle module configuration.
		OracleCfg config.OracleConfiguration

		// P2PNotaryCfg is notary module configuration.
		P2PNotaryCfg config.P2PNotary

		// StateRootCfg is stateroot module configuration.
		StateRootCfg config.StateRoot

		// ExtensiblePoolSize is the size of the pool for extensible payloads from a single sender.
		ExtensiblePoolSize int

		// BroadcastFactor is the factor (0-100) for fan-out optimization.
		BroadcastFactor int
	}
)

// NewServerConfig creates a new ServerConfig struct
// using the main applications config.
func NewServerConfig(cfg config.Config) (ServerConfig, error) {
	appConfig := cfg.ApplicationConfiguration
	protoConfig := cfg.ProtocolConfiguration
	timePerBlock := protoConfig.TimePerBlock
	if timePerBlock == 0 && protoConfig.SecondsPerBlock > 0 { //nolint:staticcheck // SA1019: protoConfig.SecondsPerBlock is deprecated
		timePerBlock = time.Duration(protoConfig.SecondsPerBlock) * time.Second //nolint:staticcheck // SA1019: protoConfig.SecondsPerBlock is deprecated
	}

	addrs, err := appConfig.GetAddresses()
	if err != nil {
		return ServerConfig{}, fmt.Errorf("failed to parse addresses: %w", err)
	}
	return ServerConfig{
		UserAgent:          cfg.GenerateUserAgent(),
		Addresses:          addrs,
		Net:                protoConfig.Magic,
		Relay:              appConfig.Relay,
		Seeds:              protoConfig.SeedList,
		DialTimeout:        time.Duration(appConfig.DialTimeout) * time.Second,
		ProtoTickInterval:  time.Duration(appConfig.ProtoTickInterval) * time.Second,
		PingInterval:       time.Duration(appConfig.PingInterval) * time.Second,
		PingTimeout:        time.Duration(appConfig.PingTimeout) * time.Second,
		MaxPeers:           appConfig.MaxPeers,
		AttemptConnPeers:   appConfig.AttemptConnPeers,
		MinPeers:           appConfig.MinPeers,
		TimePerBlock:       timePerBlock,
		OracleCfg:          appConfig.Oracle,
		P2PNotaryCfg:       appConfig.P2PNotary,
		StateRootCfg:       appConfig.StateRoot,
		ExtensiblePoolSize: appConfig.ExtensiblePoolSize,
		BroadcastFactor:    appConfig.BroadcastFactor,
	}, nil
}
