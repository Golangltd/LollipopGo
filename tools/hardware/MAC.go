package hardware


import (
	"fmt"
	"net"
)

// get mac addr
func GetMacAddr() string {
	var macaddr = ""
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul,here is what you got: " + err.Error())
	}
	for _, inter := range interfaces {
		fmt.Println(inter.Name)
		mac := inter.HardwareAddr //获取本机MAC地址
		fmt.Println("MAC = ", mac)
		macaddr = mac.String()
		if len(macaddr)!=0{
			break
		}
	}
	return macaddr
}
