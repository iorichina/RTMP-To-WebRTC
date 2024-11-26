// turnserver.go
package main

import (
	"log"
	"net"
	"regexp"
	"strconv"

	"github.com/pion/turn/v4"
)

// Run TURN server
func runTurnServer() {
	// Hardcoded values
	publicIP := "10.227.141.116"       // Hardcoded public IP address
	users := "username=password"      // Hardcoded users list (can add more like "user1=pass1,user2=pass2" if needed)
	port := 3478                       // Listening port
	realm := "pion.ly"                 // Realm (defaults to "pion.ly")

	// Create a UDP listener to pass into pion/turn
	// pion/turn itself doesn't allocate any UDP sockets, but lets the user pass them in
	// this allows us to add logging, storage or modify inbound/outbound traffic
	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	// Cache -users flag for easy lookup later
	// If passwords are stored they should be saved to your DB hashed using turn.GenerateAuthKey
	usersMap := map[string][]byte{}
	for _, kv := range regexp.MustCompile(`(\w+)=(\w+)`).FindAllStringSubmatch(users, -1) {
		usersMap[kv[1]] = turn.GenerateAuthKey(kv[1], realm, kv[2])
	}

	// Create and start the TURN server
	log.Printf("Starting TURN server...")
	_, err = turn.NewServer(turn.ServerConfig{
		Realm: realm,
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			if key, ok := usersMap[username]; ok {
				return key, true
			}
			return nil, false
		},
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(publicIP), // Hardcoded public IP address
					Address:      "0.0.0.0",              // Listening on every interface
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to start TURN server: %v", err)
		return
	}

	// Log that the TURN server has started
	log.Printf("TURN server started successfully on %s:%d\n", publicIP, port)

	// Block until user sends SIGINT or SIGTERM
	select {}
}
