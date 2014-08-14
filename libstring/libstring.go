package libstring

import (
    "os"
    "time"
    "net/url"
    "strconv"
    "strings"
    "code.google.com/p/go-uuid/uuid"
)

func Guid() string {
    return strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + uuid.New()
}

func NormalizeLocalhost(hostAndPort string) string {
    hostname, _ := os.Hostname()

    hostAndPort = strings.Replace(hostAndPort, "127.0.0.1:", hostname + ":", -1)
    hostAndPort = strings.Replace(hostAndPort, "localhost:", hostname + ":", -1)
    return hostAndPort
}

func HostWithoutPort(uri string) string {
    if !strings.Contains(uri, ":") { return uri }

    parsed, err := url.Parse(uri)
    if err != nil { return uri }

    host  := parsed.Host
    parts := strings.Split(host, ":")
    return parts[0]
}

func SplitDockerPorts(ports string) (string, string, string) {
    var (
        hostIp string
        hostPort string
        containerPort string
    )
    parts := strings.Split(ports, ":")

    for i, part := range parts {
        parts[i] = strings.TrimSpace(part)
    }

    if len(parts) == 1 {
        hostPort      = parts[0]
        containerPort = parts[0]

    } else if len(parts) == 2 {
        hostPort      = parts[0]
        containerPort = parts[1]

    } else if len(parts) == 3 {
        hostIp        = parts[0]
        hostPort      = parts[1]
        containerPort = parts[2]
    }

    return hostIp, hostPort, containerPort
}
