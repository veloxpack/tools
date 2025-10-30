# FFprobe

A minimal, statically-linked FFprobe Docker image built from scratch for fast and efficient media file analysis and metadata extraction.

## Features

- **Ultra-Lightweight**: Only 1.21 MB compressed image size
- **Static binary**: No runtime dependencies required
- **Multi-architecture**: Supports both `linux/amd64` and `linux/arm64`
- **Fast**: Minimal overhead for quick metadata extraction
- **JSON output**: Perfect for programmatic media analysis
- **Streaming Support**: HTTPS, HTTP, RTMP, RTSP, UDP protocols
- **Format Support**: MP4, MOV, MPEGTS, Matroska (MKV), FLV
- **Codec Parsing**: H.264, H.265/HEVC

## Use Cases

- Media metadata extraction
- Video/audio format detection
- Duration and bitrate analysis
- Codec information retrieval
- Stream analysis
- Automated media validation

## Image Details

- **Registry**: `ghcr.io/veloxpack/ffprobe`
- **Base**: `scratch` (no base image)
- **Image Size**: ~1.21 MB (compressed)
- **FFMPEG Version**: 8.0 (ffprobe only)
- **Alpine Build Version**: 3.22.2
- **SSL/TLS**: mbedTLS (lightweight alternative to OpenSSL)
- **Compression**: UPX with LZMA
- **Build Optimizations**: LTO, -Oz, aggressive stripping

## Pull the Image

```bash
docker pull ghcr.io/veloxpack/ffprobe:latest
```

## Usage Examples

### Get media information in JSON format

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffprobe \
  -v quiet -print_format json -show_format -show_streams \
  /workspace/video.mp4
```

### Get video duration

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffprobe \
  -v error -show_entries format=duration \
  -of default=noprint_wrappers=1:nokey=1 \
  /workspace/video.mp4
```

### Get video resolution

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffprobe \
  -v error -select_streams v:0 \
  -show_entries stream=width,height \
  -of csv=s=x:p=0 \
  /workspace/video.mp4
```

### Get codec information

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffprobe \
  -v error -select_streams v:0 \
  -show_entries stream=codec_name,codec_type \
  -of default=noprint_wrappers=1 \
  /workspace/video.mp4
```

### Get all stream details

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffprobe \
  -show_streams -show_format \
  /workspace/video.mp4
```

### Check if file is valid video

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffprobe \
  -v error -select_streams v:0 \
  -show_entries stream=codec_type \
  -of default=noprint_wrappers=1:nokey=1 \
  /workspace/video.mp4
```

### Probe remote HTTP/HTTPS stream

```bash
docker run --rm ghcr.io/veloxpack/ffprobe \
  -v quiet -print_format json -show_format -show_streams \
  https://example.com/video.mp4
```

### Probe RTMP stream

```bash
docker run --rm ghcr.io/veloxpack/ffprobe \
  -v quiet -print_format json -show_format -show_streams \
  rtmp://server/live/stream
```

### Probe RTSP stream

```bash
docker run --rm ghcr.io/veloxpack/ffprobe \
  -v quiet -print_format json -show_format -show_streams \
  rtsp://camera:554/stream1
```

## Integration Example (Shell Script)

```bash
#!/bin/bash

# Extract video metadata
metadata=$(docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffprobe \
  -v quiet -print_format json -show_format -show_streams \
  /workspace/video.mp4)

# Parse with jq
duration=$(echo "$metadata" | jq -r '.format.duration')
width=$(echo "$metadata" | jq -r '.streams[0].width')
height=$(echo "$metadata" | jq -r '.streams[0].height')

echo "Duration: ${duration}s"
echo "Resolution: ${width}x${height}"
```

## Building Locally

```bash
docker build -t ffprobe ./ffprobe
```
