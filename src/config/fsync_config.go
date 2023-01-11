package config

const (
	BAIDU_YUN = iota
)

const APP_KEY = "o7K4jkz12yuNzkYNo7KYKbpTff1S7juM"

type FsyncConfigBase interface {
}

type FsyncConfig struct {
	FsyncConfigBase
	WebConfig        *FsyncWebConfig
	ListenFileConfig []*FsyncListenFileConfig
}

type FsyncWebConfig struct {
	Port *int32
}

type FsyncListenFileConfig struct {
	ListenFileName *string
	SyncWayConfigs []*FsyncSyncWayConfig
}

type FsyncSyncWayConfig struct {
	SyncType      *int32
	RemotePathUrl *string
}
