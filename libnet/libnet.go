package libnet

import (
    "os"
    "net"
    "sort"
    "strings"
)

func RemoteToLocalHostEquality(remoteHostOrIp string) bool {
    remoteIp := net.ParseIP(remoteHostOrIp)

    if remoteIp != nil {
        localAddresses, err := net.InterfaceAddrs()
        if err == nil {
            for _, localAddress := range localAddresses {
                if strings.HasPrefix(localAddress.String(), remoteHostOrIp) {
                    return true
                }
            }
        }
    } else {
        hostname, _ := os.Hostname()
        if hostname == remoteHostOrIp {
            return true
        } else {
            localAddresses, err  := net.LookupHost(hostname)
            remoteAddresses, err := net.LookupHost(remoteHostOrIp)

            if err != nil { return false }

            sort.Sort(sort.Reverse(sort.StringSlice(localAddresses)))
            sort.Sort(sort.Reverse(sort.StringSlice(remoteAddresses)))

            for _, remoteAddress := range localAddresses {
                for _, localAddress := range localAddresses {
                    if localAddress == remoteAddress {
                        return true
                    }
                }
            }
        }
    }
    return false
}