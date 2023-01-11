package common

type SyncTerminalConfig struct {
	configType int
}

type BaiduYunConfig struct {
	AccessToken string
	ExpiresIn   int
}

type ListenedEvent struct {
	listenedFileName    string
	syncTerminalConfigs []SyncTerminalConfig
}
