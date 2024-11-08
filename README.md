
# RTMP to WebRTC Streaming Application

This project re-streams media from an RTMP source to WebRTC using a Go server and FFmpeg for RTP streaming. It provides real-time audio and video from a specified media file to a WebRTC client in the browser, using WebSockets for automated SDP and ICE exchange.
Features

Streams audio and video from a media file (e.g., .mp4) to a WebRTC client.
WebSocket-based signaling for automated SDP and ICE candidate exchange.
Configurable RTP ports for audio and video ingestion.
Improved error handling and logging for production-ready robustness.

## Project Structure

#### 1) main.go
Go server handling WebRTC peer connections, RTP packet ingestion, and WebSocket signaling.

#### 2) main.go
app.js Client-side JavaScript for WebRTC connection and signaling.

#### 3) index.html
Frontend UI for initiating and viewing the WebRTC stream.

Replace input.mp4 with the path to your media file.

## Start server

```bash
go run main.go
```
    
## ffmpeg

```bash
ffmpeg -re -i input.mp4 -map 0:v -c:v libvpx -payload_type 96 -ssrc 1 -f rtp rtp://127.0.0.1:5004 -map 0:a -c:a libopus -payload_type 111 -ssrc 2 -f rtp rtp://127.0.0.1:5004
```

#### Note:

Replace input.mp4 with the path to your media file.

