// servers_response.go
package main

import (
	"fmt"
	"github.com/blacknight2018/GoOut/api"
	"github.com/blacknight2018/GoOut/utils"
	"log"
	"net"
	"strconv"
)

/*
	Implements servers response of SOCKS5 for non Linux systems
*/
func server_response(local_conn net.Conn, remote_address string) {
	load_balancer := get_load_balancer()

	local_tcpaddr, _ := net.ResolveTCPAddr("tcp4", load_balancer.address)
	remote_tcpaddr, _ := net.ResolveTCPAddr("tcp4", remote_address)

	log.Println("[DEBUG]", remote_address, "->", load_balancer.address)
	local_conn.Write([]byte{5, SUCCESS, 0, 1, 0, 0, 0, 0, 0, 0})

	remoteIP := (remote_tcpaddr.IP.String())
	if utils.IsChinaIP(remoteIP) || GoOutServer == nil || len(*GoOutServer) == 0 {
		//Direct
		remote_conn, err := net.DialTCP("tcp4", local_tcpaddr, remote_tcpaddr)

		if err != nil {
			log.Println("[WARN]", remote_address, "->", load_balancer.address, fmt.Sprintf("{%s}", err))
			local_conn.Write([]byte{5, NETWORK_UNREACHABLE, 0, 1, 0, 0, 0, 0, 0, 0})
			local_conn.Close()
			return
		}
		pipe_connections(local_conn, remote_conn)

	} else {
		//Through By GoOut
		GoOut := *GoOutServer
		port := strconv.Itoa(remote_tcpaddr.Port)
		api.TcpOnProxy(local_conn, local_tcpaddr, remote_tcpaddr.IP.String(), port, &GoOut)
	}

}
