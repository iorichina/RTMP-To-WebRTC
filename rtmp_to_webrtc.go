// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

//go:build !js
// +build !js

// rtp-to-webrtc demonstrates how to consume a RTP stream video UDP, and then send to a WebRTC client.
package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/pion/webrtc/v4"
)

func main() {
	log.Printf("RTMP-To-WebRTC starting \n")
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlan,
	})
	if err != nil {
		panic(err)
	}

	// Open a UDP Listener for RTP Packets on port 5004
	streamPort := 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: streamPort})
	if err != nil {
		panic(err)
	}
	log.Printf("Stream Server starting on %v\n", streamPort)

	// Increase the UDP receive buffer size
	// Default UDP buffer sizes vary on different operating systems
	bufferSize := 300000 // 300KB
	err = listener.SetReadBuffer(bufferSize)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video_id1", "stream_id1")
	if err != nil {
		panic(err)
	}
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// sdp endpoint -----------------------------------------------
	http.HandleFunc("/sdp", handleSdp(peerConnection))

	port := ":8080"
	log.Printf("Http Server starting on %s\n", port)
	http.ListenAndServe(port, nil)

	// Read RTP packets forever and send them to the WebRTC Client
	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("error during read: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				// The peerConnection has been closed.
				return
			}

			panic(err)
		}
	}

}

func handleSdp(peerConnection *webrtc.PeerConnection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bs, err := io.ReadAll(r.Body)
		if nil != err {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}
		offer := webrtc.SessionDescription{}
		if err = json.Unmarshal(bs, &offer); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ss, _ := json.Marshal(&offer)
			w.Write([]byte(fmt.Sprintf("err=%v offer=%v body=%v", err, string(ss), string(bs))))
			return
		}

		// Set the remote SessionDescription
		if err = peerConnection.SetRemoteDescription(offer); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("err=%v  body=%v", err, string(bs))))
			return
		}

		// Create answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("err=%v  body=%v", err, string(bs))))
			return
		}

		// Sets the LocalDescription, and starts our UDP listeners
		if err = peerConnection.SetLocalDescription(answer); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("err=%v  body=%v", err, string(bs))))
			return
		}

		// Output the answer in base64 so we can paste it in browser
		fmt.Println(encode(peerConnection.LocalDescription()))
		b, err := json.Marshal(peerConnection.LocalDescription())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("err=%v  body=%v", err, string(bs))))
			return
		}
		w.Write(b)
	}
}

// JSON encode + base64 a SessionDescription
func encode(obj *webrtc.SessionDescription) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode a base64 and unmarshal JSON into a SessionDescription
func decode(in string, obj *webrtc.SessionDescription) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, obj); err != nil {
		panic(err)
	}
}
