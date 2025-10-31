# FFmpeg Lite

A full-featured, statically-linked FFmpeg Docker image optimized for video transcoding with support for modern codecs including AV1, VP9, H.264, H.265, and audio codecs like Opus and MP3.

## Features

- **Modern codecs**: SVT-AV1, VP8/VP9, H.264 (x264), H.265 (x265)
- **Audio codecs**: MP3 (LAME), Opus
- **Static binaries**: No runtime dependencies required
- **Multi-architecture**: Supports both `linux/amd64` and `linux/arm64`
- **Built from scratch**: Minimal image size with maximum performance
- **Production-ready**: Optimized for high-quality video transcoding

## Use Cases

- High-quality video transcoding
- Multi-codec video conversion
- Adaptive bitrate streaming preparation
- Professional video processing
- Batch video encoding

## Image Details

- **Registry**: `ghcr.io/veloxpack/ffmpeg:8.0-lite`
- **Base**: `scratch` (no base image)
- **FFmpeg Version**: 8.0
- **Alpine Build Version**: 3.22.2

## Included Libraries

- **SVT-AV1** v1.7.0 - Next-gen AV1 encoder
- **libvpx** v1.13.0 - VP8/VP9 codecs
- **x264** - H.264 encoder
- **x265** v3.5 - H.265/HEVC encoder
- **LAME** v3.100 - MP3 encoder
- **Opus** v1.4 - Modern audio codec
- **MbedTLS** v3.4.1 - Secure communications

## Pull the Image

```bash
docker pull ghcr.io/veloxpack/ffmpeg:8.0-lite
```

## Variants

This is the full-featured FFmpeg image. For specialized use cases, consider these lightweight variants:

- **`ghcr.io/veloxpack/ffmpeg:8.0-thumbnail`** - Ultra-lightweight (2.39 MB) for thumbnail generation
- **`ghcr.io/veloxpack/ffmpeg:8.0-split`** - Optimized (3.92 MB) for video splitting and scene detection
- **`ghcr.io/veloxpack/ffmpeg:8.0-concat`** - Minimal (914 KB) for video concatenation

## Usage Examples

### Transcode to H.264

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-lite \
  -i /workspace/input.mp4 \
  -c:v libx264 -preset medium -crf 23 \
  -c:a aac -b:a 128k \
  /workspace/output.mp4
```

### Transcode to H.265 (HEVC)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-lite \
  -i /workspace/input.mp4 \
  -c:v libx265 -preset medium -crf 28 \
  -c:a aac -b:a 128k \
  /workspace/output.mp4
```

### Transcode to AV1 (modern, efficient)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-lite \
  -i /workspace/input.mp4 \
  -c:v libsvtav1 -preset 6 -crf 35 \
  -c:a libopus -b:a 128k \
  /workspace/output.mp4
```

### Transcode to VP9 (WebM)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-lite \
  -i /workspace/input.mp4 \
  -c:v libvpx-vp9 -crf 30 -b:v 0 \
  -c:a libopus \
  /workspace/output.webm
```

### Create HLS streaming variants

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-lite \
  -i /workspace/input.mp4 \
  -c:v libx264 -preset fast \
  -map 0:v -s 1920x1080 -b:v 5000k -maxrate 5000k -bufsize 10000k \
  -map 0:v -s 1280x720 -b:v 2800k -maxrate 2800k -bufsize 5600k \
  -map 0:v -s 854x480 -b:v 1400k -maxrate 1400k -bufsize 2800k \
  -map 0:a -c:a aac -b:a 128k \
  -f hls -hls_time 6 -hls_playlist_type vod \
  -master_pl_name master.m3u8 \
  -var_stream_map "v:0,a:0 v:1,a:0 v:2,a:0" \
  /workspace/stream_%v.m3u8
```

## Building Locally

```bash
docker build -t ghcr.io/veloxpack/ffmpeg:8.0-lite ./ffmpeg-lite
```
