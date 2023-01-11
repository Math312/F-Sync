package config

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)
import "github.com/fsync/common"

const DEFAULT_BAIDUYUN_CONFIG_FILE = "fsync_config_baidu_yun.fsync"

type FsyncBaiduYunConfigContainer struct {
	Config         *FsyncBaiduYunConfig
	ConfigFileName string
}

var containerInitialized uint32

var container FsyncBaiduYunConfigContainer

var containerMu sync.Mutex

func GetBaiduYunConfigContainer(configFile string) FsyncBaiduYunConfigContainer {
	if atomic.LoadUint32(&containerInitialized) == 1 {
		return container
	}
	containerMu.Lock()
	defer containerMu.Unlock()

	if containerInitialized == 0 {
		container = FsyncBaiduYunConfigContainer{
			ConfigFileName: configFile,
		}
		var data []byte
		config := FsyncBaiduYunConfig{}
		if common.FileIsExist(configFile) {
			file, err := os.Open(configFile)
			if err != nil {
				log.Fatalf("Open BaiduYun Config File Failure, File name: %s, Error: %s", configFile, err)
			}
			data, err = io.ReadAll(file)
			if err != nil {
				log.Fatalf("Read BaiduYun Config File Failure, File name: %s, Error: %s", configFile, err)
			}
		} else {
			if common.FileIsExist(DEFAULT_BAIDUYUN_CONFIG_FILE) {
				file, err := os.Open(DEFAULT_BAIDUYUN_CONFIG_FILE)
				if err != nil {
					log.Fatalf("Open BaiduYun Config File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
				}
				data, err = io.ReadAll(file)
				if err != nil {
					log.Fatalf("Read BaiduYun Config File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
				}
			}
		}
		if data != nil {
			readErr := common.ReadJson(data, &config)
			if readErr != nil {
				log.Fatalf("Parse BaiduYun Config File Failure, Error: %s ", readErr)
			}
			container.Config = &config
		} else {
			_, err := os.Create(DEFAULT_BAIDUYUN_CONFIG_FILE)
			if err != nil {
				log.Fatalf("Create BaiduYun Config File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
			}
			json, err := common.WriteJson(&config)
			if err != nil {
				log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
			}
			err = os.WriteFile(DEFAULT_BAIDUYUN_CONFIG_FILE, json, 0666)
			if err != nil {
				log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
			}
			container.Config = &config
		}
		if container.Config.OAuthConfigContainer == nil {
			container.Config.OAuthConfigContainer = make(map[string]*FsyncBaiduYunOAuthConfig, 10)
		}
		atomic.StoreUint32(&containerInitialized, 1)
	}
	return container
}

type FsyncBaiduYunConfig struct {
	OAuthConfigContainer map[string]*FsyncBaiduYunOAuthConfig `json:"oAuthConfigContainer"`
}

type FsyncBaiduYunOAuthConfig struct {
	AccessToken *string `json:"accessToken"`
	ExpireTime  int64   `json:"expireTime"`
	UniqueName  *string `json:"uniqueName"`
}

func WriteBaiduYunOAuthDataByUniqueName(uniqueName string, accessToken string, expireTime int64) error {
	if atomic.LoadUint32(&containerInitialized) != 1 {
		return errors.New("BaiduYun Config ")
	}
	containerMu.Lock()
	defer containerMu.Unlock()
	container.Config.OAuthConfigContainer[uniqueName] = &FsyncBaiduYunOAuthConfig{
		AccessToken: &accessToken,
		ExpireTime:  expireTime,
		UniqueName:  &uniqueName,
	}
	json, err := common.WriteJson(container.Config)
	if err != nil {
		log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
	}
	err = os.WriteFile(DEFAULT_BAIDUYUN_CONFIG_FILE, json, 0666)
	if err != nil {
		log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
	}
	return nil
}

func ListAllBaiduYunOAuthData() ([]FsyncBaiduYunOAuthConfig, error) {
	if atomic.LoadUint32(&containerInitialized) != 1 {
		return nil, errors.New("BaiduYun Config ")
	}
	containerMu.Lock()
	defer containerMu.Unlock()
	result := make([]FsyncBaiduYunOAuthConfig, len(container.Config.OAuthConfigContainer))
	idx := 0
	for _, value := range container.Config.OAuthConfigContainer {
		result[idx] = *value
		idx++
	}
	return result, nil
}

func ReadBaiduYunOAuthDataByUniqueName(uniqueName string, accessToken string) (string, error) {
	if atomic.LoadUint32(&containerInitialized) != 1 {
		return "", errors.New("BaiduYun Config ")
	}
	containerMu.Lock()
	defer containerMu.Unlock()
	value, ok := container.Config.OAuthConfigContainer[uniqueName]
	if !ok {
		return "", errors.New("OAuth Config Not Found, Config Name: " + uniqueName)
	} else {
		if value.ExpireTime > time.Now().UnixMilli() {
			return "", errors.New("OAuth AccessToken Expired")
		} else {
			return *(value.AccessToken), nil
		}
	}
}

func DeleteBaiduYunOAuthDataByUniqueName(uniqueName string) error {
	if atomic.LoadUint32(&containerInitialized) != 1 {
		return errors.New("BaiduYun Config ")
	}
	containerMu.Lock()
	defer containerMu.Unlock()
	delete(container.Config.OAuthConfigContainer, uniqueName)
	json, err := common.WriteJson(container.Config)
	if err != nil {
		log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
	}
	err = os.WriteFile(DEFAULT_BAIDUYUN_CONFIG_FILE, json, 0666)
	if err != nil {
		log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
	}
	return nil
}

func WriteBaiduYunConfigToFile(configFileName string) error {
	if atomic.LoadUint32(&containerInitialized) != 1 {
		return errors.New("BaiduYun Config ")
	}
	containerMu.Lock()
	defer containerMu.Unlock()
	json, err := common.WriteJson(container.Config)
	if err != nil {
		log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
	}
	err = os.WriteFile(DEFAULT_BAIDUYUN_CONFIG_FILE, json, 0666)
	if err != nil {
		log.Fatalf("Write Default BaiduYun Config To File Failure, File name: %s, Error: %s", DEFAULT_BAIDUYUN_CONFIG_FILE, err)
	}
	return nil
}

type BaiduOAuthData struct {
	AccessToken string
	ExpiresIn   int
	Scope       string
}

func ParseBaiduOAuthUrl(url string) (*BaiduOAuthData, error) {
	m, err := common.ParseUrlFragmentParameters(url)
	if err != nil {
		return nil, err
	}
	expiresIn, _ := strconv.Atoi(m["expires_in"])
	return &BaiduOAuthData{
		AccessToken: m["access_token"],
		ExpiresIn:   expiresIn,
		Scope:       m["scope"],
	}, nil
}
