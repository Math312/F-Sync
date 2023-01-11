package common

import (
	"context"
	"encoding/json"
	"fmt"
	openapiclient "github.com/fsync/common/openxpanapi"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type PreCreateReturnType struct {
	Path       string        `json:"path"`
	Uploadid   string        `json:"uploadid"`
	ReturnType int           `json:"return_type"`
	BlockList  []interface{} `json:"block_list"`
	Errno      int           `json:"errno"`
	RequestID  int64         `json:"request_id"`
}

type DepartUploadResponse struct {
	Errno int    `json:"errno"`
	Md5   string `json:"md5"`
}

func preUploadFile(accessToken string, fileParts []string, fileName string, baiduYunPath string, md5Map map[string]string) (*PreCreateReturnType, error) {
	path := baiduYunPath               // string
	isdir := int32(0)                  // int32
	size, err := GetFileSize(fileName) // int32
	if err != nil {
		log.Fatalf("File Error: %s, %s", fileName, err)
		return nil, err
	}
	autoinit := int32(1) // int32
	blockList := "["     // string
	for idx, filePart := range fileParts {
		blockList += "\"" + md5Map[filePart] + "\""
		if idx != len(md5Map)-1 {
			blockList += ","
		}
	}
	blockList += "]"
	rtype := int32(3) // int32 | rtype (optional)

	configuration := openapiclient.NewConfiguration()
	api_client := openapiclient.NewAPIClient(configuration)
	resp, r, err := api_client.FileuploadApi.Xpanfileprecreate(context.Background()).AccessToken(accessToken).Path(path).Isdir(isdir).Size(size).Autoinit(autoinit).BlockList(blockList).Rtype(rtype).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FileuploadApi.Xpanfileprecreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Xpanfileprecreate`: Fileprecreateresponse
	fmt.Fprintf(os.Stdout, "Response from `FileuploadApi.Xpanfileprecreate`: %v\n", resp)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", r)
	}

	result := &PreCreateReturnType{}
	if err := json.Unmarshal(bodyBytes, result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func departUpload(accessToken string, preUploadResult PreCreateReturnType, fileName string, fileParts []string, baiduYunPath string, md5Map map[string]string) map[string]*DepartUploadResponse {
	uploadid := preUploadResult.Uploadid // string
	type_ := "tmpfile"
	partIndex := 0
	result := make(map[string]*DepartUploadResponse, len(fileParts))
	for _, filePart := range fileParts {
		partseq := strconv.Itoa(partIndex) // string
		response, err := doDepartUpload(accessToken, filePart, partseq, baiduYunPath, uploadid, type_)
		if err != nil {
			log.Fatalf("Depart upload error, File name: %s, Error: %s", filePart, err)
		} else {
			result[filePart] = response
		}

		partIndex++
	}
	return result
}

func createUploadFile(accessToken string, fileName string, baiduYunPath string, preCreateReturnResponse PreCreateReturnType, fileParts []string, md5Map map[string]string) {
	path := baiduYunPath
	isdir := int32(0)                  // int32
	size, err := GetFileSize(fileName) // int32
	if err != nil {
		log.Fatalf("File Error: %s, %s", fileName, err)
		return
	}
	uploadid := preCreateReturnResponse.Uploadid // string
	blockList := "["                             // string
	for idx, filePart := range fileParts {
		blockList += "\"" + md5Map[filePart] + "\""
		if idx != len(md5Map)-1 {
			blockList += ","
		}
	}
	blockList += "]"
	rtype := int32(3) // int32 | rtype (optional)

	configuration := openapiclient.NewConfiguration()
	api_client := openapiclient.NewAPIClient(configuration)
	resp, r, err := api_client.FileuploadApi.Xpanfilecreate(context.Background()).AccessToken(accessToken).Path(path).Isdir(isdir).Size(size).Uploadid(uploadid).BlockList(blockList).Rtype(rtype).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FileuploadApi.Xpanfilecreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Xpanfilecreate`: Filecreateresponse
	fmt.Fprintf(os.Stdout, "Response from `FileuploadApi.Xpanfilecreate`: %v\n", resp)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", r)
	}

	fmt.Println(string(bodyBytes))
}

func doDepartUpload(accessToken string, partFileName string, partseq string, path string, uploadid string, type_ string) (*DepartUploadResponse, error) {
	file, err := os.Open(partFileName)
	if err != nil {
		log.Fatalf("File Open Error: %s, %s", partFileName, err)
		return nil, err
	}
	defer file.Close()
	configuration := openapiclient.NewConfiguration()
	//configuration.Debug = true
	api_client := openapiclient.NewAPIClient(configuration)
	resp, r, err := api_client.FileuploadApi.Pcssuperfile2(context.Background()).AccessToken(accessToken).Partseq(partseq).Path(path).Uploadid(uploadid).Type_(type_).File(file).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FileuploadApi.Pcssuperfile2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Pcssuperfile2`: string
	fmt.Fprintf(os.Stdout, "Response from `FileuploadApi.Pcssuperfile2`: %v\n", resp)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", r)
	}

	departUploadResponse := &DepartUploadResponse{}
	if err := json.Unmarshal(bodyBytes, departUploadResponse); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return departUploadResponse, nil
}

// todo 处理AccessToken超时
func UploadFileToBaiduYun(accessToken string, fileName string, baiduYunPath string) {
	var fileParts []string
	size, err := GetFileSize(fileName)
	if err != nil {
		log.Fatalf("File Error: %s, %s", fileName, err)
		return
	}
	if size < 4*1024*1024 {
		fileParts = []string{fileName}
	} else {
		fileParts = FileDepart(fileName, 4*1024*1024)
	}
	md5Map := make(map[string]string, len(fileParts))
	for _, filePart := range fileParts {
		buf, err := os.ReadFile(filePart)
		if err != nil {
			panic(err)
		}
		md5 := GenMd5(buf)
		md5Map[filePart] = md5
	}
	preUploadResult, err := preUploadFile(accessToken, fileParts, fileName, baiduYunPath, md5Map)
	if err != nil {
		log.Fatalf("PreUpload Error, File name: %s, Error: %s", fileName, err)
	}
	departUpload(accessToken, *preUploadResult, fileName, fileParts, baiduYunPath, md5Map)
	createUploadFile(accessToken, fileName, baiduYunPath, *preUploadResult, fileParts, md5Map)
}
