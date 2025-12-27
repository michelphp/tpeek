package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	Blue   = "\033[34m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Yellow = "\033[33m"
)

/*
 * TPeek - TCP Proxy Debugger
 * Copyright (C) 2025  Michel.F
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by 
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
func main() {
	listenAddr := flag.String("l", "", "Local address and port to listen on (e.g. 0.0.0.0:8080)")
	targetAddr := flag.String("t", "", "Target address and port (e.g. MySQL, Redis : 127.0.0.1:3306)")
	dumpHex := flag.Bool("hex", false, "Enable full hexadecimal dump mode")
	flag.Parse()

	if *targetAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *listenAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("❌ Fatal: Could not listen on %s: %v", *listenAddr, err)
	}

	fmt.Println("🚀 TPeek TCP Proxy Initialized")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf(" 🔵 LISTEN  : %s\n", *listenAddr)
	fmt.Printf(" 🟢 TARGET  : %s\n", *targetAddr)
	fmt.Printf(" 🛠️  MODE    : %s\n", map[bool]string{true: "HEXADECIMAL", false: "PLAIN TEXT"}[*dumpHex])
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📡 Waiting for incoming connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("⚠️  Accept Error: %v", err)
			continue
		}

		log.Printf("✅ New connection from %s", conn.RemoteAddr())
		go handleProxy(conn, *targetAddr, *dumpHex)
	}
}

func handleProxy(clientConn net.Conn, target string, useHex bool) {
	defer clientConn.Close()

	remoteConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("❌ Target connection failed (%s): %v", target, err)
		return
	}
	defer remoteConn.Close()

	done := make(chan struct{})
	go func() {
		pipe(clientConn, remoteConn, "🔵 CLIENT >>", useHex, Blue)
		done <- struct{}{}
	}()
	go func() {
		pipe(remoteConn, clientConn, "🟢 SERVER <<", useHex, Green)
		done <- struct{}{}
	}()
	<-done
	log.Printf("🔌 Connection closed for %s", clientConn.RemoteAddr())
}

func pipe(src, dst net.Conn, label string, useHex bool, color string) {
	buffer := make([]byte, 16384)
	for {
		n, err := src.Read(buffer)
		if n > 0 {
			data := buffer[:n]
			timestamp := time.Now().Format("15:04:05.000")
			fmt.Printf("\n%s┌── %s ── %s ── %d bytes%s\n", color, label, timestamp, n, Reset)

			if useHex {
				fmt.Print(hex.Dump(data))
			} else {
				cleanText := strings.Map(func(r rune) rune {
					if r < 32 || r > 126 {
						return '.'
					}
					return r
				}, string(data))
				fmt.Printf("%s│%s %s\n", color, Reset, cleanText)
			}

			fmt.Printf("%s└%s\n", color, strings.Repeat("─", 40)+Reset)

			_, writeErr := dst.Write(data)
			if writeErr != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
}
