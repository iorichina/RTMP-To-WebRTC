RTMP to WebRTC Streaming Application
This project re-streams media from an RTMP source to WebRTC using a Go server and FFmpeg for RTP streaming. It provides real-time audio and video from a specified media file to a WebRTC client in the browser, using WebSockets for automated SDP and ICE exchange.


go run main.go

ffmpeg -re -i input.mp4 -map 0:v -c:v libvpx -payload_type 96 -ssrc 1 -f rtp rtp://127.0.0.1:5004 -map 0:a -c:a libopus -payload_type 111 -ssrc 2 -f rtp rtp://127.0.0.1:5006



