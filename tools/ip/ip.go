package ip

import (
	"errors"
	"fmt"
	"net"

	"LollipopGo/tools/sample"
	"github.com/jinzhu/gorm"
)

type IP struct {
	Start       string `gorm:"primary_key"`
	End         string `gorm:"primary_key"`
	StartNum    int64 //ip对应的长整型
	EndNum      int64
	Continent   string //大洲
	Country     string
	Province    string
	City        string
	District    string
	ISP         string
	AreaCode    string
	CountryEN   string
	CountryCode string
	Longitude   string //经度
	Latitude    string //维度
}

const DefaultLocation = "未知地区"

func GetIPInfo(db *gorm.DB, ip string) (*IP, error) {
	var (
		info IP
		err  error
	)
	if err = db.Where("inet_aton(?) between start_num and end_num", ip).First(&info).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func GetIPLocation(ip *IP, err error) (string, error) {
	if err != nil || ip == nil {
		return DefaultLocation, nil
	}
	if ip.Province == ip.City || ip.City == "" {
		return ip.Province, nil
	}
	return fmt.Sprintf("%s.%s", ip.Province, ip.City), nil
}

func GetRandomLocation(db *gorm.DB) string {
	offset := sample.RandInt(0, 10000)
	var (
		info IP
		err  error
	)
	//只取大陆ip
	if err = db.Where("country_code = ?", "CN").Offset(offset).First(&info).Error; err != nil {
		return DefaultLocation
	}
	if info.Province == info.City || info.City == "" {
		return info.Province
	}
	return fmt.Sprintf("%s.%s", info.Province, info.City)
}

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func IsPrivateIP(ip string) (bool, error) {
	addr := net.ParseIP(ip)
	if addr == nil {
		return false, errors.New("invalid ip")
	}
	for _, block := range privateIPBlocks {
		if block.Contains(addr) {
			return true, nil
		}
	}
	return false, nil
}
