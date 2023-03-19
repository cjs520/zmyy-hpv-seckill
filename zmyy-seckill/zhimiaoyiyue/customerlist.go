package zhimiaoyiyue

import (
	"errors"
	"fmt"
	"log"
	"zmyy_seckill/consts"
	"zmyy_seckill/fetcher"
	"zmyy_seckill/model"
	"zmyy_seckill/utils"
)

//获取指定接种地ID
func (e *ZMYYEngine) GetCustomerList(ip ...string) (int, error) {
	params := "[\"" + e.Conf.Province + "\",\"" + e.Conf.City + "\",\"" + e.Conf.District + "\"]"
	newUrl := consts.CustomerListUrl + "&city=" + utils.UrlEncode(params)  + "&lat=31.23037&lng=121.4737" + "&id=0&cityCode=" + e.Conf.CityCode + "&product=0"
	headers := make(map[string]string)
	headers["User-Agent"] = consts.UserAgent
	headers["Referer"] = consts.Refer
	zftsl := utils.GetZFTSL()
	headers["zftsl"] = zftsl
	bytes, err2 := fetcher.FetchWithRatelimter(newUrl, headers, ip...)
	if err2 != nil {
		return -1, err2
	}
	customers := model.CustomerList{}
	err2 = utils.Transfer2CustomerListModel(bytes, &customers)
	if err2 != nil {
		return -1, err2
	}
	fmt.Printf("正在查找接种地点：\n")
	for k, v := range customers.Customers {
		log.Printf("第 %d个接种地：%s\n", k+1, v.Cname)
		if v.Cname == e.Conf.CustomerName {
			log.Printf("====选中第 %d个接种地：%s，其customerId为 %d====\n", k+1, v.Cname, v.Id)
			return v.Id, nil
		}
	}
	log.Printf("未找到指定接种地，请对比配置文件接种地是否正确！\n")
	return -1, errors.New("未找到指定接种地，请对比配置文件接种地是否正确！")
}
