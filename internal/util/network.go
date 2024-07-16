package util

import "net"

func GetLocalIpAddress() (string, error) {
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, networkInterface := range networkInterfaces {
		addresses, err := networkInterface.Addrs()
		if err != nil {
			return "", err
		}
		for _, address := range addresses {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}
