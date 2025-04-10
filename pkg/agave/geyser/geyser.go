package geyser

import (
	"encoding/json"
	"fmt"
)

const (
	PathLibYellowstoneGRPC = "/usr/lib/libyellowstone_grpc_geyser.so"
)

// GeyserPlugin provides native support for YellowstoneGRPC geyser plugins through structured configuration.
// For other geyser plugins, use GenericPluginConfig to specify a JSON config string that must contain
// a top-level "libpath" field pointing to the plugin's shared library based on where it will be installed
// on the host machine.

type GeyserPlugin struct {
	YellowstoneGRPC     *YellowstoneGRPC `pulumi:"yellowstoneGRPC,optional"`
	GenericPluginConfig *string          `pulumi:"genericPluginConfig,optional"`
}

func (g *GeyserPlugin) Check() error {
	if g.YellowstoneGRPC != nil && g.GenericPluginConfig != nil {
		return fmt.Errorf("only one of YellowstoneGRPC or GenericPluginConfig can be specified")
	}

	if g.YellowstoneGRPC == nil && g.GenericPluginConfig == nil {
		return fmt.Errorf("either YellowstoneGRPC or GenericPluginConfig must be specified")
	}

	if g.YellowstoneGRPC != nil {
		if g.YellowstoneGRPC.Config == nil && g.YellowstoneGRPC.JSON == nil {
			return fmt.Errorf("either Config or JSON must be specified if YellowstoneGRPC is used")
		}
		if g.YellowstoneGRPC.Config != nil && g.YellowstoneGRPC.JSON != nil {
			return fmt.Errorf("only one of Config or JSON can be specified")
		}
	}

	return nil
}

func (g *GeyserPlugin) ToConfigString() (string, error) {
	if g.YellowstoneGRPC != nil {
		res, err := g.YellowstoneGRPC.MarshalJSON()
		if err != nil {
			return "", err
		} else {
			return string(res), nil
		}
	} else {
		// Return generic config string as-is
		return *g.GenericPluginConfig, nil
	}
}

type YellowstoneGRPC struct {
	// JSON field is mutually exclusive with Config
	// This is used in case the user wants to provide a JSON string directly
	JSON    *string `pulumi:"json,optional"`
	Config  *Config `pulumi:"config,optional"`
	Version string  `pulumi:"version"`
}

func (y *YellowstoneGRPC) MarshalJSON() ([]byte, error) {
	tempMap := make(map[string]any)
	if y.Config != nil {
		// Marshal the original config to get its JSON representation
		origBytes, err := json.Marshal(y.Config)
		if err != nil {
			return nil, err
		}

		// Marshall it into a temp structure so that we can add our own fields
		if err = json.Unmarshal(origBytes, &tempMap); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte(*y.JSON), &tempMap); err != nil {
			return nil, err
		}
	}
	// Add the libpath field
	tempMap["libpath"] = PathLibYellowstoneGRPC

	// Convert the combined data to pretty-printed JSON string
	configBytes, err := json.MarshalIndent(tempMap, "", "  ")
	if err != nil {
		return nil, err
	}

	return configBytes, nil
}

type Config struct {
	Log              *GrpcConfigLog        `json:"log,omitempty" pulumi:"log,optional"`
	Tokio            *GrpcConfigTokio      `json:"tokio,omitempty" pulumi:"tokio,optional"`
	Grpc             GrpcConfigGrpc        `json:"grpc" pulumi:"grpc"`
	Prometheus       *GrpcConfigPrometheus `json:"prometheus,omitempty" pulumi:"prometheus,optional"`
	DebugClientsHTTP *bool                 `json:"debug_clients_http,omitempty" pulumi:"debugClientsHttp,optional"`
}

type GrpcConfigLog struct {
	Level *string `json:"level,omitempty" pulumi:"level,optional"`
}

type GrpcConfigTokio struct {
	WorkerThreads *int   `json:"worker_threads,omitempty" pulumi:"workerThreads,optional"`
	Affinity      *[]int `json:"affinity" pulumi:"affinity,optional"`
}

type GrpcConfigGrpc struct {
	Address                           string                     `json:"address" pulumi:"address"`
	TLSConfig                         *GrpcConfigGrpcServerTLS   `json:"tls_config,omitempty" pulumi:"tlsConfig,optional"`
	Compression                       *GrpcConfigGrpcCompression `json:"compression,omitempty" pulumi:"compression,optional"`
	MaxDecodingMessageSize            *int                       `json:"max_decoding_message_size,omitempty" pulumi:"maxDecodingMessageSize,optional"`
	SnapshotPluginChannelCapacity     *int                       `json:"snapshot_plugin_channel_capacity,omitempty" pulumi:"snapshotPluginChannelCapacity,optional"`
	SnapshotClientChannelCapacity     *int                       `json:"snapshot_client_channel_capacity,omitempty" pulumi:"snapshotClientChannelCapacity,optional"`
	ChannelCapacity                   *int                       `json:"channel_capacity,omitempty" pulumi:"channelCapacity,optional"`
	UnaryConcurrencyLimit             *int                       `json:"unary_concurrency_limit,omitempty" pulumi:"unaryConcurrencyLimit,optional"`
	UnaryDisabled                     *bool                      `json:"unary_disabled,omitempty" pulumi:"unaryDisabled,optional"`
	FilterLimits                      *GrpcConfigFilterLimits    `json:"filter_limits,omitempty" pulumi:"filterLimits,optional"`
	XToken                            *string                    `json:"x_token,omitempty" pulumi:"xToken,optional"`
	FilterNameSizeLimit               *int                       `json:"filter_name_size_limit,omitempty" pulumi:"filterNameSizeLimit,optional"`
	FilterNamesSizeLimit              *int                       `json:"filter_names_size_limit,omitempty" pulumi:"filterNamesSizeLimit,optional"`
	FilterNamesCleanupInterval        *string                    `json:"filter_names_cleanup_interval,omitempty" pulumi:"filterNamesCleanupInterval,optional"`
	ReplayStoredSlots                 *int64                     `json:"replay_stored_slots,omitempty" pulumi:"replayStoredSlots,optional"`
	ServerHttp2AdaptiveWindow         *bool                      `json:"server_http2_adaptive_window,omitempty" pulumi:"serverHttp2AdaptiveWindow,optional"`
	ServerHttp2KeepaliveInterval      *string                    `json:"server_http2_keepalive_interval,omitempty" pulumi:"serverHttp2KeepaliveInterval,optional"`
	ServerHttp2KeepaliveTimeout       *string                    `json:"server_http2_keepalive_timeout,omitempty" pulumi:"serverHttp2KeepaliveTimeout,optional"`
	ServerInitialConnectionWindowSize *int32                     `json:"server_initial_connection_window_size,omitempty" pulumi:"serverInitialConnectionWindowSize,optional"`
	ServerInitialStreamWindowSize     *int32                     `json:"server_initial_stream_window_size,omitempty" pulumi:"serverInitialStreamWindowSize,optional"`
}

type GrpcConfigGrpcServerTLS struct {
	CertPath string `json:"cert_path" pulumi:"certPath"`
	KeyPath  string `json:"key_path" pulumi:"keyPath"`
}

type GrpcConfigGrpcCompression struct {
	Accept []string `json:"accept,omitempty" pulumi:"accept,optional"`
	Send   []string `json:"send,omitempty" pulumi:"send,optional"`
}

type GrpcConfigPrometheus struct {
	Address string `json:"address" pulumi:"address"`
}

type GrpcConfigFilterLimits struct {
	Accounts           *GrpcConfigFilterLimitsAccounts     `json:"accounts,omitempty" pulumi:"accounts,optional"`
	Slots              *GrpcConfigFilterLimitsSlots        `json:"slots,omitempty" pulumi:"slots,optional"`
	Transactions       *GrpcConfigFilterLimitsTransactions `json:"transactions,omitempty" pulumi:"transactions,optional"`
	TransactionsStatus *GrpcConfigFilterLimitsTransactions `json:"transactions_status,omitempty" pulumi:"transactionsStatus,optional"`
	Blocks             *GrpcConfigFilterLimitsBlocks       `json:"blocks,omitempty" pulumi:"blocks,optional"`
	BlocksMeta         *GrpcConfigFilterLimitsBlocksMeta   `json:"blocks_meta,omitempty" pulumi:"blocksMeta,optional"`
	Entries            *GrpcConfigFilterLimitsEntries      `json:"entries,omitempty" pulumi:"entries,optional"`
}

type GrpcConfigFilterLimitsAccounts struct {
	Max           *int     `json:"max,omitempty" pulumi:"max,optional"`
	Any           *bool    `json:"any,omitempty" pulumi:"any,optional"`
	AccountMax    *int     `json:"account_max,omitempty" pulumi:"accountMax,optional"`
	AccountReject []string `json:"account_reject,omitempty" pulumi:"accountReject,optional"`
	OwnerMax      *int     `json:"owner_max,omitempty" pulumi:"ownerMax,optional"`
	OwnerReject   []string `json:"owner_reject,omitempty" pulumi:"ownerReject,optional"`
	DataSliceMax  *int     `json:"data_slice_max,omitempty" pulumi:"dataSliceMax,optional"`
}

type GrpcConfigFilterLimitsSlots struct {
	Max *int `json:"max,omitempty" pulumi:"max,optional"`
}

type GrpcConfigFilterLimitsTransactions struct {
	Max                  *int     `json:"max,omitempty" pulumi:"max,optional"`
	Any                  *bool    `json:"any,omitempty" pulumi:"any,optional"`
	AccountIncludeMax    *int     `json:"account_include_max,omitempty" pulumi:"accountIncludeMax,optional"`
	AccountIncludeReject []string `json:"account_include_reject,omitempty" pulumi:"accountIncludeReject,optional"`
	AccountExcludeMax    *int     `json:"account_exclude_max,omitempty" pulumi:"accountExcludeMax,optional"`
	AccountRequiredMax   *int     `json:"account_required_max,omitempty" pulumi:"accountRequiredMax,optional"`
}

type GrpcConfigFilterLimitsBlocks struct {
	Max                  *int     `json:"max,omitempty" pulumi:"max,optional"`
	AccountIncludeMax    *int     `json:"account_include_max,omitempty" pulumi:"accountIncludeMax,optional"`
	AccountIncludeAny    *bool    `json:"account_include_any,omitempty" pulumi:"accountIncludeAny,optional"`
	AccountIncludeReject []string `json:"account_include_reject,omitempty" pulumi:"accountIncludeReject,optional"`
	IncludeTransactions  *bool    `json:"include_transactions,omitempty" pulumi:"includeTransactions,optional"`
	IncludeAccounts      *bool    `json:"include_accounts,omitempty" pulumi:"includeAccounts,optional"`
	IncludeEntries       *bool    `json:"include_entries,omitempty" pulumi:"includeEntries,optional"`
}

type GrpcConfigFilterLimitsBlocksMeta struct {
	Max *int `json:"max,omitempty" pulumi:"max,optional"`
}

type GrpcConfigFilterLimitsEntries struct {
	Max *int `json:"max,omitempty" pulumi:"max,optional"`
}
