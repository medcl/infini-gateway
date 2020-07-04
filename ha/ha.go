package ha

import (
    "fmt"
    "net"
    "os"
    "os/exec"
)

func SetupVIP() {
    alias:="en0:1"

    ///sbin/ifconfig eth0:1 192.168.3.198 netmask 255.255.255.0
    setupVIP := exec.Command("/sbin/ifconfig", alias, "192.168.3.213", "netmask", "255.255.255.0")
    _, err := setupVIP.CombinedOutput()
    if err != nil {
        fmt.Printf("failed to invoke /sbin/ifconfig to set floating IP on interface: %s", err)
        os.Exit(1)
    }

    intf, err := net.InterfaceByName(alias)
    if intf == nil {
        fmt.Printf("failed to locate interface by alias %s: %s\n", alias, err)
        os.Exit(1)
    }

    fmt.Println("success!!")
}
