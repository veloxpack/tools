# FFMPEG Minimal

A minimal, statically-linked FFMPEG Docker image built from scratch with essential codecs. Optimized for lightweight operations like thumbnail generation and basic video processing.

## Features

- **Minimal size**: Built from scratch with only essential binaries
- **Static binaries**: No runtime dependencies required
- **Multi-architecture**: Supports both `linux/amd64` and `linux/arm64`
- **Includes**: Both `ffmpeg` and `ffprobe` binaries
- **Codecs**: Basic codec support with MbedTLS for secure communications

## Use Cases

- Thumbnail generation
- Basic video/audio conversion
- Quick video processing tasks
- Metadata extraction with ffprobe
- Lightweight video operations

## Image Details

- **Registry**: `ghcr.io/veloxpack/ffmpeg-minimal`
- **Base**: `scratch` (no base image)
- **FFMPEG Version**: 8.0
- **Alpine Build Version**: 3.22.2

## Pull the Image

```bash
docker pull ghcr.io/veloxpack/ffmpeg-minimal:latest
```

## Usage Examples

### Generate a thumbnail

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-minimal \
  -i /workspace/video.mp4 \
  -ss 00:00:10 \
  -vframes 1 \
  /workspace/thumbnail.jpg
```

### Convert video to different format

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-minimal \
  -i /workspace/input.mp4 \
  -c:v libx264 -preset fast \
  /workspace/output.mp4
```

### Use ffprobe for metadata

```bash
docker run --rm -v $(pwd):/workspace \
  --entrypoint /ffprobe \
  ghcr.io/veloxpack/ffmpeg-minimal \
  -v quiet -print_format json -show_format -show_streams \
  /workspace/video.mp4
```

## Building Locally

```bash
docker build -t ffmpeg-minimal ./ffmpeg-minimal
```
