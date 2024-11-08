
# RTMP to WebRTC Streaming Application

This project re-streams media from a source to WebRTC using a Go server and FFmpeg for RTP streaming. It provides real-time audio and video from a specified media file to a WebRTC client in the browser, using WebSockets for automated SDP and ICE exchange.
Features

Streams audio and video from a media file (e.g., .mp4) to a WebRTC client.
WebSocket-based signaling for automated SDP and ICE candidate exchange.
Configurable RTP ports for audio and video ingestion.
Improved error handling and logging for production-ready robustness.

## Prerequisites

Before getting started, make sure you have the following tools installed:

- Go Programming Language, , download it from [the Go website](https://golang.org/dl/).
- ffmpeg, [Get ffmpeg](https://www.ffmpeg.org/download.html).

## Installation

To install the required dependencies, follow these steps:

1. **Clone the repository**
2. **Add these modules**
   ```bash
   go get github.com/pion/webrtc/v4
   go get github.com/pion/rtp
   go get github.com/gorilla/websocket

## Start server

```bash
go run main.go
```
    
## ffmpeg

```bash
ffmpeg -re -i input.mp4 -map 0:v -c:v libvpx -payload_type 96 -ssrc 1 -f rtp rtp://127.0.0.1:5004 -map 0:a -c:a libopus -payload_type 111 -ssrc 2 -f rtp rtp://127.0.0.1:5004
```

or

```bash
ffmpeg -re -i rtmp://your-rtmp-server/stream-key -map 0:v -c:v libvpx -payload_type 96 -ssrc 1 -f rtp rtp://127.0.0.1:5004 -map 0:a -c:a libopus -payload_type 111 -ssrc 2 -f rtp rtp://127.0.0.1:5004
```

#### Note:

Replace input.mp4/rtmp_url with the path to your media file.



### Packages used

This project uses the following Go packages:

- [`github.com/pion/webrtc/v4`](https://github.com/pion/webrtc) – A WebRTC API implementation for Go.
- [`github.com/pion/rtp`](https://github.com/pion/rtp) – RTP (Real-Time Protocol) handling for Go.
- [`github.com/gorilla/websocket`](https://github.com/gorilla/websocket) – A WebSocket implementation for Go.

