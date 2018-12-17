package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var proxyCmd = &cobra.Command{
	Use:   "portproxy",
	Short: "PortProxy is a tool to proxy TCP/UDP port via aother",
	Long: `A Simple TCP/UDP tunnel the service via a different port.
				Written in GO by enbiso.
				Complete documentation is available at https://www.enbiso.com/portproxy`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy(
			cmd.PersistentFlags().Lookup("source").Value.String(),
			cmd.PersistentFlags().Lookup("dest").Value.String(),
			strings.ToLower(cmd.PersistentFlags().Lookup("protocol").Value.String()),
		)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of PortProxy",
	Long:  `All software has versions. This is PortProxy's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PortProxy v0.1.0")
	},
}

func main() {
	proxyCmd.PersistentFlags().StringP("source", "s", "", "Source Address (127.0.0.1:80)")
	proxyCmd.MarkPersistentFlagRequired("source")
	proxyCmd.PersistentFlags().StringP("dest", "d", "", "Destination Address (0.0.0.0:8080)")
	proxyCmd.MarkPersistentFlagRequired("dest")
	proxyCmd.PersistentFlags().StringP("protocol", "p", "tcp", "Protocol (TCP/UDP)")
	viper.AutomaticEnv()
	viper.BindPFlag("source", proxyCmd.PersistentFlags().Lookup("source"))
	viper.BindPFlag("dest", proxyCmd.PersistentFlags().Lookup("dest"))
	viper.BindPFlag("protocol", proxyCmd.PersistentFlags().Lookup("protocol"))

	proxyCmd.AddCommand(versionCmd)

	if err := proxyCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func proxy(source string, target string, proto string) {
	log.Info("Port Proxy")
	ln, err := net.Listen(proto, target)
	if err != nil {
		log.Fatal(err)
	}
	// accept connection on port
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
