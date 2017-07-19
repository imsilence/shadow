package utils

import (
    "net"
)


func GetInterfaces() (interfaces map[string][]string, err error) {
    inters, err := net.Interfaces()
    if nil != err {
        return nil, err
    }
    interfaces = make(map[string][]string)

    for _, inter := range inters {
        mac := inter.HardwareAddr.String()
        if "" == mac {
            continue
        }
        interfaces[mac] = make([]string, 0)
        addrs, _ := inter.Addrs()
        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && nil != ipnet.IP.To4() {
                interfaces[mac] = append(interfaces[mac], ipnet.IP.String())
            }
        }
    }
    return interfaces, nil
}
