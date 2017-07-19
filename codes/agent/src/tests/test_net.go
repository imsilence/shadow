package main

import (
    "net"
    "fmt"
    "encoding/json"
    "os"
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
            if ipnet, ok := addr.(*net.IPNet); ok {
                fmt.Println(ipnet.IP.DefaultMask())
                interfaces[mac] = append(interfaces[mac], ipnet.IP.String())
            }
        }
    }
    return interfaces, nil
}

func main() {
    addrs, _ := net.InterfaceAddrs()
    for _, addr := range addrs {
        fmt.Println(addr)
        fmt.Println(addr.Network(), addr.String())
    }
    var a []int = make([]int, 1)

    a = append(a, 1)
    fmt.Println(a)

    fmt.Println("==========")
    interfaces, _ := GetInterfaces()
    j, _ := json.Marshal(interfaces)
    fmt.Println(string(j))
    os.Exit(0)

}
