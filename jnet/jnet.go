package jnet

import (
	"errors"
	"net"
)

var ErrLocalIPFailed = errors.New("can't get local IP")
var ErrLocalMacFailed = errors.New("can't get local mac")

// PickUnusedPort 防止指定端口已经被占用，可以用下面的函数来动态获取一个可用端口。
func PickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

// LocalIP 获取本地 ip
func LocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", ErrLocalIPFailed
}

// LocalMac 获取本地 mac地址
func LocalMac() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			return "", err
		}

		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return inter.HardwareAddr.String(), nil
				}
			}
		}
	}

	return "", ErrLocalMacFailed
}
