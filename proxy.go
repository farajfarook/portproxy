package main

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func proxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portproxy",
		Short: "PortProxy is a tool to proxy TCP/UDP port via aother",
		Long: `A Simple TCP/UDP tunnel the service via a different port.
					Written in GO by enbiso.
					Complete documentation is available at https://www.enbiso.com/portproxy`,
		Run: func(cmd *cobra.Command, args []string) {
			proxy(
				cmd.Flags().Lookup("source").Value.String(),
				cmd.Flags().Lookup("target").Value.String(),
				strings.ToLower(cmd.Flags().Lookup("protocol").Value.String()),
			)
		},
	}
	cmd.Flags().StringP("source", "s", "127.0.0.1:80", "Source Address (127.0.0.1:80)")
	cmd.MarkFlagRequired("source")
	cmd.Flags().StringP("target", "t", ":8080", "Target Address (:8080)")
	cmd.MarkFlagRequired("target")
	cmd.Flags().StringP("protocol", "p", "tcp", "Protocol (TCP/UDP)")
	viper.BindPFlag("source", cmd.Flags().Lookup("source"))
	viper.BindPFlag("target", cmd.Flags().Lookup("target"))
	viper.BindPFlag("protocol", cmd.Flags().Lookup("protocol"))
	return cmd
}

func proxy(source string, target string, proto string) {
	log.Info("Port Proxy")
	log.Info("Source   - " + source)
	log.Info("Target   - " + target)
	log.Info("Protocol - " + proto)

	if proto == "udp" || proto == "*" {
		protocolUDPProxy("udp", target, source)
	}

	if proto == "tcp" || proto == "*" {
		protocolTCPProxy("tcp", target, source)
	}

}

func protocolUDPProxy(proto string, target string, source string) {
	ln, err := net.ListenPacket(proto, target)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		go func() {
			buf := make([]byte, 1024)
			clientConn, err := net.Dial(proto, source)
			if err != nil {
				log.Error(err)
				clientConn.Close()
				return
			}
			_, err = clientConn.Read(buf)
			if err != nil {
				ln.Close()
				log.Error(err)
				return
			}

			ln.WriteTo(buf, &net.IPAddr{IP: net.IP(target)})

		}()
	}

}

func protocolTCPProxy(proto string, target string, source string) {
	ln, err := net.Listen(proto, target)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, _ := ln.Accept()
		go func() {
			clientConn, err := net.Dial(proto, source)
			if err != nil {
				log.Error(err)
				conn.Close()
				return
			}
			pipe(conn, clientConn)
		}()
	}
}

func pipe(conn1 net.Conn, conn2 net.Conn) {
	chan1 := getChannel(conn1)
	chan2 := getChannel(conn2)
	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return
			}
			conn2.Write(b1)
		case b2 := <-chan2:
			if b2 == nil {
				return
			}
			conn1.Write(b2)
		}
	}
}

func getChannel(conn net.Conn) chan []byte {
	c := make(chan []byte)
	go func() {
		b := make([]byte, 1024)
		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()
	return c
}
