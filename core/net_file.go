package core

import (
	"encoding/json"
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lock       sync.Mutex
	cipher     *utils.Cipher
	jsonParser fastjson.Parser
)

func init() {
	cipherNew, err := utils.NewCipher()
	if err != nil {
		config.Logger.Fatal(err.Error())
	}
	cipher = cipherNew
}

// 115的文件对象，这个对象指向的可能是文件，也可能是文件夹
type NetFile struct {
	// 有fid就是文件
	Fid string `json:"fid"`
	// 有cid但没有Fid，就是文件夹
	Cid string `json:"cid"`
	// 文件大小
	Size int    `json:"s"`
	Name string `json:"n"`
	Sha  string `json:"sha"`
	Pc   string `json:"pc"`
}

type downloadBody struct {
	State bool   `json:"state"`
	Msg   string `json:"msg"`
	Data  string `json:"data"`
	Errno int    `json:"errno"`
}

type importBody struct {
	Status     int `json:"status"`
	StatusCode int `json:"statuscode"`
}

// 开启一定量的worker，通过channel接收任务，channel有一定的缓冲区
// worker在接收到任务后执行任务，当遇到需要人机验证的时候，改变全局
// 变量，然后进入循环等待模式，期间一直检测，直到人机验证完成；
// Note: 只要用到Lock的地方，都要考虑超时问题
// 将当前网络文件的内容导出到目录中
func (file *NetFile) Export() string {
	// 保证worker不会panic
	defer func() {
		if err := recover(); err != nil {
			config.Logger.Error("export link error",
				zap.String("name", file.Name),
				zap.String("reason", fmt.Sprintf("%v", err)))

			// 在报错的情况下，如果依然处于人机验证的阻塞状态，就解除状态。不
			// 希望因为一个处于检测状态的任务panic而导致整个任务卡死
			if config.SpiderVerification {
				config.SpiderVerification = false
			}
		}
	}()

	url, cookie := file.extractDownloadInfo()
	if cookie == "" || url == "" {
		return ""
	}

	fileSha1 := file.extractFileSha1(url, cookie)
	if fileSha1 == "" {
		return ""
	}

	joinStrings := []string{file.Name, strconv.Itoa(file.Size), file.Sha, fileSha1}
	result := strings.Join(joinStrings, config.LinkSep)

	var formatSize string
	sizeM := file.Size >> 20
	if sizeM == 0 {
		formatSize = "小于1MB"
	} else {
		formatSize = fmt.Sprintf("%dMB", sizeM)
	}

	nowFormat := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s 导出成功，大小: %s\t文件: %s\n", nowFormat, formatSize, file.Name)
	config.Logger.Info("export success", zap.String("name", file.Name), zap.Int("size", file.Size))
	return result
}

func (file *NetFile) extractDownloadInfo() (downloadUrl, cookie string) {
	downUrl := "https://proapi.115.com/app/chrome/downurl"
	headers := config.GetFakeHeaders(true)
	for {
		// 先检查是否在等待人机验证状态
		headOff := config.SpiderVerification
		if headOff {
			// 太长时间的停滞之后，确保真的有worker去查询，设置
			// 一个超时时间
			if int(time.Now().Unix())-config.SpiderStatWaitAliveTime > config.SpiderStatWaitTimeout {
				config.SpiderStatWaitAliveTime = int(time.Now().Unix())
				goto Work
			}
			config.Logger.Info(fmt.Sprintf("waiting Man-machine verification: %s", file.Name))
			time.Sleep(config.SpiderCheckInterval / 2)
			continue
		}

	Work:
		text, err := cipher.Encrypt([]byte(fmt.Sprintf(`{"pickcode":"%s"}`, file.Pc)))
		if err != nil {
			config.Logger.Warn("encrypt data error",
				zap.String("name", file.Name))
			return
		}
		postData := map[string]string{"data": string(text)}
		body, response, err := utils.PostFormWithResponse(downUrl, headers, postData)
		if err != nil {
			config.Logger.Warn("export file network error",
				zap.String("name", file.Name))
			return
		}

		parsedDownloadBody := new(downloadBody)
		err = json.Unmarshal(body, parsedDownloadBody)
		if err != nil {
			config.Logger.Warn("parse download body fail",
				zap.String("content", string(body)),
				zap.String("name", file.Name))
			return
		}

		// ==================验证文件下载地址是否正常==================
		// 文件状态异常
		if !parsedDownloadBody.State {
			config.Logger.Warn("download file state odd",
				zap.String("content", parsedDownloadBody.Msg),
				zap.String("name", file.Name))
			return
		}

		// 有多个worker因为时间差，都进入人机检测验证状态，也无所谓
		// 进入人机验证之后，反复检测状态
		if parsedDownloadBody.Errno == 911 {
			fmt.Println("发现人机验证，请到115浏览器中播放任意一个视频，完成人机检测...")
			config.Logger.Warn("found Man-machine verification， waiting...")
			config.SpiderVerification = true
			config.SpiderStatWaitAliveTime = int(time.Now().Unix())
			time.Sleep(config.SpiderCheckInterval)
			goto Work
		}
		// 如果有人机验证状态,取消人机验证状态
		if config.SpiderVerification {
			config.SpiderVerification = false
		}

		// 返回的下载信息中不包含下载地址
		if parsedDownloadBody.Data == "" {
			config.Logger.Warn("download file body not contain download url",
				zap.String("content", fmt.Sprintf("%v", parsedDownloadBody)),
				zap.String("name", file.Name))
			return
		}
		// 下载的时候有自己单独的cookie,提取下载cookie
		downloadCookie := extractDownloadCookie(response)
		if downloadCookie == "" {
			config.Logger.Warn("get download cookie fail", zap.String("name", file.Name))
			return
		}

		downloadUrlContent, err := cipher.Decrypt([]byte(parsedDownloadBody.Data))
		if err != nil {
			config.Logger.Warn("decrypt download url error",
				zap.String("name", file.Name))
			return
		}

		extractedDownloadUrl, err := extractDownloadUrl(downloadUrlContent, file.Fid)
		if err != nil {
			config.Logger.Warn("parse download url error",
				zap.String("name", file.Name))
			return
		}
		return extractedDownloadUrl, downloadCookie
	}
}

func extractDownloadCookie(response *http.Response) string {
	newCookie, ok := response.Header["Set-Cookie"]
	if ok && len(newCookie) >= 1 {
		cookies := strings.SplitN(newCookie[0], ";", 2)
		if len(cookies) >= 2 {
			return cookies[0]
		}
	}
	return ""
}

func extractDownloadUrl(downloadContent []byte, fileId string) (string, error) {
	parsedValue, err := jsonParser.ParseBytes(downloadContent)
	if err != nil {
		return "", err
	}
	return string(parsedValue.GetStringBytes(fileId, "url", "url")), nil
}

func (file *NetFile) extractFileSha1(downloadUrl, cookie string) string {
	downloadHeader := config.GetFakeRangeHeaders()
	downloadHeader["Cookie"] = cookie

	body, err := utils.Get(downloadUrl, downloadHeader, nil)
	if err != nil {
		config.Logger.Warn("get file header to calculate file sha1 fail", zap.String("name", file.Name))
		return ""
	}

	sha1 := utils.Sha1(body)
	return strings.ToUpper(sha1)
}

// 导入时，要指定文件所属的Cid（文件夹），文件夹不存在就需要创建；
// 创建文件夹的方法是，指定这个文件夹的父文件夹，填入文件夹的名字
// 之后创建，返回Cid，在做这个任务的时候，Cid需要是创建好的文件夹；
// 创建文件夹的工作在调用这个函数的地方提前准备好，这里不涉及创建文
// 件夹
func (file *NetFile) Import() bool {
	// 保证worker不会panic
	defer func() {
		if err := recover(); err != nil {
			config.Logger.Error("import link error",
				zap.String("name", file.Name),
				zap.String("reason", fmt.Sprintf("%v", err)))
		}
	}()

	if file.Cid == "" {
		config.Logger.Warn("empty target dir")
		return false
	}

	target := config.DirTargetPrefix + file.Cid
	shaFirstJoinStrings := []string{config.UserId, file.Sha, file.Sha, target, "0"}
	shaFirstRaw := strings.Join(shaFirstJoinStrings, "")
	shaFirst := utils.Sha1([]byte(shaFirstRaw))

	shaSecondRaw := config.UserKey + shaFirst + config.EndString
	sig := strings.ToUpper(utils.Sha1([]byte(shaSecondRaw)))

	url := fmt.Sprintf("http://uplb.115.com/3.0/initupload.php?isp=0&appid=0&appversion=%s&format=json&sig=%s",
		config.AppVer, sig)
	postData := map[string]string{
		"preid":    file.Pc,
		"filename": file.Name,
		"quickid":  file.Sha,
		"app_ver":  config.AppVer,
		"filesize": strconv.Itoa(file.Size),
		"userid":   config.UserId,
		"exif":     "",
		"target":   target,
		"fileid":   file.Sha,
	}

	headers := config.GetFakeHeaders(true)
	body, err := utils.PostForm(url, headers, postData)
	if err != nil {
		config.Logger.Warn("import file network error",
			zap.String("name", file.Name))
		return false
	}

	parsedImportBody := new(importBody)
	err = json.Unmarshal(body, parsedImportBody)
	if err != nil {
		config.Logger.Warn("parse import body fail",
			zap.String("content", string(body)),
			zap.String("name", file.Name))
		return false
	}

	if parsedImportBody.Status == 2 && parsedImportBody.StatusCode == 0 {
		var formatSize string
		sizeM := file.Size >> 20
		if sizeM == 0 {
			formatSize = "小于1MB"
		} else {
			formatSize = fmt.Sprintf("%dMB", sizeM)
		}
		nowFormat := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("%s 导入成功，大小: %s\t文件: %s\n", nowFormat, formatSize, file.Name)
		return true
	} else {
		config.Logger.Warn("import info not expect",
			zap.String("content", string(body)),
			zap.String("name", file.Name))
		return false
	}
}

// 格式不对，创建NetFile不一定成功
func CreateNetFile(fileInfo string) *NetFile {
	splitStrings := strings.Split(fileInfo, config.LinkSep)
	if len(splitStrings) != 4 {
		return nil
	}
	size, err := strconv.Atoi(splitStrings[1])
	if err != nil {
		return nil
	}

	return &NetFile{
		Name: splitStrings[0],
		Size: size,
		Sha:  splitStrings[2],
		Pc:   splitStrings[3],
	}
}
