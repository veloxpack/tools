# FFmpeg Thumbnail

A highly optimized, ultra-lightweight FFmpeg Docker image built specifically for thumbnail and sprite/storyboard generation. Stripped down to only essential video processing capabilities.

## Features

- **Ultra-Lightweight**: Only 2.39 MB compressed image size
- **Static binary**: No runtime dependencies required
- **Multi-architecture**: Supports both `linux/amd64` and `linux/arm64`
- **Optimized for Thumbnails**: Specifically configured for thumbnail and sprite generation
- **Video Decoders**: H.264, VP8, VP9 support
- **Image Encoders**: PNG and JPEG/MJPEG output
- **Streaming Support**: HTTP/HTTPS protocols
- **Format Support**: MP4, MOV, Matroska (WebM/MKV)

## Image Details

- **Registry**: `ghcr.io/veloxpack/ffmpeg-thumbnail`
- **Base**: `scratch` (no base image)
- **Image Size**: 2.39 MB (compressed)
- **FFmpeg Version**: 8.0
- **Alpine Build Version**: 3.22.2
- **SSL/TLS**: mbedTLS (lightweight alternative to OpenSSL)
- **Compression**: UPX with LZMA
- **Build Optimizations**: LTO, -Oz, aggressive stripping

## Included Components

### Video Decoders
- H.264 (AVC)
- VP8 (WebM)
- VP9 (WebM)

### Image Encoders
- PNG (with zlib compression)
- MJPEG/JPEG

### Filters
- `scale` - Resize frames
- `thumbnail` - Select best representative frames
- `fps` - Frame rate control
- `tile` - Create sprite sheets/storyboards
- `select` - Advanced frame selection
- `format` - Pixel format conversion

### Protocols
- File (local files)
- HTTP/HTTPS (remote streams)

## Use Cases

- Generate video thumbnails (PNG/JPEG)
- Create sprite sheets / storyboards
- Extract best frames from videos
- Process remote video streams
- Lightweight video frame extraction

## Pull the Image

```bash
docker pull ghcr.io/veloxpack/ffmpeg-thumbnail:latest
```

## Usage Examples

### Generate a thumbnail at specific time

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -ss 5 \
  -vframes 1 \
  /output/thumbnail.jpg
```

### Extract best frame using thumbnail filter

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -vf thumbnail \
  -frames:v 1 \
  /output/best-frame.jpg
```

### Generate PNG thumbnail

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -ss 10 \
  -vframes 1 \
  /output/thumbnail.png
```

### Process remote HTTP/HTTPS stream

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i https://example.com/video.mp4 \
  -ss 5 \
  -vframes 1 \
  /output/thumbnail.jpg
```

### Generate storyboard / sprite sheet (5x5 grid)

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -vf "fps=1/10,scale=160:90,tile=5x5" \
  /output/storyboard.jpg
```

### Create storyboard from remote WebM video

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i https://example.com/video.webm \
  -loglevel error \
  -vf "fps=1/10,scale=160:90,tile=5x5" \
  -threads 4 \
  /output/storyboard.jpg
```

### Multiple thumbnails at intervals

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -vf "fps=1/60" \
  /output/thumb-%04d.jpg
```

### Generate scaled thumbnail with custom resolution

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -ss 5 \
  -vf "scale=320:180" \
  -vframes 1 \
  /output/thumbnail.jpg
```

## Advanced Examples

### Storyboard with custom grid and scaling

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -vf "fps=1/30,scale=200:112,tile=8x8" \
  /output/storyboard-8x8.png
```

### Select best frames at specific intervals

```bash
docker run --rm -v $(pwd):/output \
  ghcr.io/veloxpack/ffmpeg-thumbnail \
  -i /output/video.mp4 \
  -vf "select='not(mod(n\,300))',scale=320:180" \
  -vsync 0 \
  /output/frame-%04d.jpg
```

## Limitations

### Not Included
- Audio processing (all audio codecs/filters disabled)
- Video encoding (only image encoding: PNG/JPEG)
- Advanced video codecs (HEVC, AV1, MPEG-4 removed for size)
- Hardware acceleration
- Complex filters (blur, overlay, etc.)
- Streaming protocols (RTMP, RTSP, UDP)

### What This Image IS For
✅ Thumbnail generation
✅ Sprite/storyboard creation
✅ Frame extraction
✅ Basic image scaling

### What This Image IS NOT For
❌ Full video encoding/transcoding
❌ Audio processing
❌ Live streaming
❌ Complex video editing

## Performance Tips

- Use `-threads 4` for faster processing on multi-core systems
- Use `-loglevel error` to suppress unnecessary output
- Use JPEG for smaller file sizes, PNG for quality
- Pre-scale videos before tiling for faster storyboard generation

## Building Locally

```bash
docker build -t ffmpeg-thumbnail ./ffmpeg-thumbnail
```

