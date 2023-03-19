package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zmyy_seckill/consts"
	"zmyy_seckill/model"
)

func Transfer2Model(jsonCont []byte, m interface{}) (interface{}, error) {
	switch m.(type) {
	case model.CustomerList:
		res := &model.CustomerList{}
		err := json.Unmarshal(jsonCont, res)
		if err != nil {
			fmt.Printf("Transfer2Model err:%v\n", err)
			return nil, err
		}
		return *res, nil
	}
	return nil, nil
}
func Transfer2CustomerListModel(jsonCont []byte, cumtomers *model.CustomerList) error {
	err := json.Unmarshal(jsonCont, &cumtomers)
	if err != nil {
		fmt.Printf("Transfer2CustomerListModel err:%v\n", err)
		return err
	}
	return nil
}
func Transfer2CustomerProductListModel(jsonCont []byte, m *model.RootSource) error {
	err := json.Unmarshal(jsonCont, &m)
	if err != nil {
		return err
	}
	return nil
}

func Transfer2SubscribeDateModel(jsonCont []byte, m *model.SubscribeDate) error {
	err := json.Unmarshal(jsonCont, &m)
	if err != nil {
		return err
	}
	return nil
}
func Transfer2VerifyResultModel(jsonCont []byte, m *model.VerifyResultModel) error {
	err := json.Unmarshal(jsonCont, &m)
	if err != nil {
		return err
	}
	return nil
}
func Transfer2SubscribeDateDetailModel(jsonCont []byte, m *model.SubscribeDateDetail) error {
	err := json.Unmarshal(jsonCont, &m)
	if err != nil {
		return err
	}
	return nil
}

// 将Base64文件（../imgs/veryfiPics）转成图片
func Base64ToPics(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Base64ToPics() can not load file err : %v\n", err)
		return err
	}
	m := &model.VerifyPicModel{}
	err = json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	d, _ := base64.StdEncoding.DecodeString(m.Dragon)
	t, _ := base64.StdEncoding.DecodeString(m.Tiger)
	fd, _ := os.OpenFile(filePath+"-dragon.png", os.O_RDWR|os.O_CREATE, os.ModePerm)
	ft, _ := os.OpenFile(filePath+"-tiger.png", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer fd.Close()
	defer ft.Close()
	_, err = fd.Write(d)
	if err != nil {
		fmt.Printf("Base64文件转图片失败！ err : %v", err)
		return err
	}
	_, err = ft.Write(t)
	if err != nil {
		fmt.Printf("Base64文件转图片失败！ err : %v", err)
		return err
	}
	return nil
}
func GetCookie() (string, error) {
	path := GetCurrentPath()
	exePath := path + "/getCookie.py"
	args := []string{exePath}
	cmd := exec.Command("python", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err1 := cmd.Run()
	if err1 != nil {
		panic(err1)
	}
	// 从输出中获取最后一行
	lines := strings.Split(out.String(), "\n")
	lastLine := lines[len(lines)-2] // 去除最后一个空行
	//fmt.Println(lastLine)
	return string(lastLine), nil

}

/*
*

	调用Python脚本破解验证码
*/
func CallPythonScript(tigerPath, dragonPath, procssPath string) (string, error) {
	path := GetCurrentPath()
	exePath := path + "/captcha.py"
	args := []string{exePath, tigerPath, dragonPath, procssPath}
	out, err := exec.Command("python", args...).Output()
	if err != nil {
		fmt.Printf("滑块验证码识别失败！ 图片为： %s,  err: %v\n", dragonPath, err.Error())
		return "", err
	}
	str := strings.Replace(string(out), "\r", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	//fmt.Printf("滑块坐标为： %s\n", str)
	return str, nil
}
func GetZFTSL() string {
	d := []byte("zfsw_" + strconv.FormatInt(time.Now().Unix()/10, 10))
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func UrlEncode(orgUrl string) string {
	encodeUrl := url.QueryEscape(orgUrl)
	return encodeUrl
}
func ParseSessionId(s string) string {
	const sessionIdRe = `ASP.NET_SessionId=([^;]+)`
	compile := regexp.MustCompile(sessionIdRe)
	match := compile.FindSubmatch([]byte(s))
	return string(match[1])
}
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Get current process path failed . err : %v \n", err)
		return ""
	}
	//dir = strings.Replace(dir, "\\", "/", -1)
	//const pathRe = `([0-9a-zA-z:]*[0-9a-zA-Z/]+/zmyy_seckill)`
	//compile := regexp.MustCompile(pathRe)
	//match := compile.FindSubmatch([]byte(dir))
	//if len(match) > 1 {
	//	dir = string(match[1])
	//}
	return dir
}

func DeleteFile(path ...string) {
	for _, v := range path {
		err := os.Remove(v)
		if err != nil {
			fmt.Printf("删除文件%s失败：%v\n", v, err)
		}
	}
}

func SetRandomIP() string {
	selectedIp := ""
	if consts.ProxyIpArr == nil || len(consts.ProxyIpArr) == 0 {
		//如果IP池IP用尽，则使用本机IP
		return selectedIp
	}
	//随机从ip池中获取一个ip
	index := rand.Intn(len(consts.ProxyIpArr))
	selectedIp = consts.ProxyIpArr[index]
	//删除该ip
	consts.ProxyIpArr = append(consts.ProxyIpArr[:index], consts.ProxyIpArr[index+1:]...)
	return selectedIp
}
