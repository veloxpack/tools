# VeloxPack Tools

This repository contains production-ready, statically-linked Docker images for media processing tools. All images are built from scratch with no base OS, ensuring minimal size, maximum security, and no runtime dependencies.

## Available Tools

### [FFmpeg Lite](./ffmpeg-lite)
**Full-featured video transcoding powerhouse**

A complete FFmpeg build (16.84 MB) with modern codecs including AV1, VP9, H.264, H.265, Opus, and MP3. Perfect for high-quality video transcoding, multi-codec conversion, and adaptive bitrate streaming preparation.

```bash
docker pull ghcr.io/veloxpack/ffmpeg:8.0-lite
```

**Key Features:**
- Modern codecs: SVT-AV1, VP8/VP9, H.264 (x264), H.265 (x265)
- Audio codecs: MP3 (LAME), Opus, AAC
- Multi-architecture support (amd64/arm64)
- Production-ready for professional video processing

---

### [FFmpeg Thumbnail](./ffmpeg-thumbnail)
**Ultra-lightweight thumbnail & sprite generation**

An ultra-optimized FFmpeg build (2.39 MB) specifically for thumbnail and storyboard generation. Stripped down to only essential video decoders and image encoders.

```bash
docker pull ghcr.io/veloxpack/ffmpeg:8.0-thumbnail
```

**Key Features:**
- Ultra-minimal size (2.39 MB)
- Perfect for thumbnail/sprite generation
- PNG and JPEG output support
- H.264, VP8, VP9 video decoders

---

### [FFmpeg Split](./ffmpeg-split)
**Video splitting & scene detection**

Specialized FFmpeg build for splitting videos with scene detection support. Includes x264 encoder for re-encoding and stream copying for fast, lossless splits.

```bash
docker pull ghcr.io/veloxpack/ffmpeg:8.0-split
```

**Key Features:**
- Ultra-lightweight (3.92 MB)
- Scene detection with metadata export
- Fast stream copying
- H.264 encoding support

---

### [FFmpeg Concat](./ffmpeg-concat)
**Video concatenation & merging**

Specialized FFmpeg build for concatenating multiple videos. Outputs MP4 and WebM formats only. Supports advanced file list formats including duration metadata, trimming, and stream selection.

```bash
docker pull ghcr.io/veloxpack/ffmpeg:8.0-concat
```

**Key Features:**
- Ultra-lightweight (914 KB)
- MP4 and WebM output support
- Multiple concat methods (demuxer, protocol, filter)
- Advanced file list formats with trimming

---

### [FFprobe](./ffprobe)
**Media metadata extraction**

A standalone FFprobe image for fast and efficient media file analysis, metadata extraction, and format detection. Ideal for automated media validation and programmatic analysis.

```bash
docker pull ghcr.io/veloxpack/ffmpeg:8.0-probe
```

**Key Features:**
- Ultra-lightweight (1.21 MB)
- JSON output support
- Fast metadata extraction
- Stream analysis and codec detection

---

### [Shaka Packager](./shaka-packager)
**Adaptive streaming packaging**

Google's industry-standard tool (14.39 MB) for DASH/HLS manifest generation, DRM encryption, and adaptive streaming preparation. Supports Widevine, PlayReady, and FairPlay encryption.

```bash
docker pull ghcr.io/veloxpack/shaka-packager:latest
```

**Key Features:**
- DASH & HLS packaging
- DRM support (Widevine, PlayReady, FairPlay)
- Multi-bitrate preparation
- Live and VOD streaming support

---

## Why VeloxPack Tools?

### Performance
- **Static binaries** - No runtime dependencies or library conflicts
- **Built from scratch** - Minimal image size with maximum performance
- **Optimized builds** - Latest versions compiled with optimal flags

### Security
- **No base OS** - Reduced attack surface with scratch-based images
- **Minimal footprint** - Only essential binaries included
- **Static linking** - No external library vulnerabilities

### Multi-Architecture
- **AMD64** - Intel/AMD processors
- **ARM64** - Apple Silicon, AWS Graviton, and ARM servers
- **Unified images** - Same image works across platforms

### Production Ready
- **Industry standards** - Based on widely-used open-source tools
- **Tested workflows** - Proven in production environments
- **Regular updates** - Maintained with latest stable versions

---

## Quick Start

### Example: Complete Video Processing Workflow

```bash
# 1. Analyze source video
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-probe \
  -v quiet -print_format json -show_format -show_streams \
  /workspace/source.mp4

# 2. Generate thumbnail
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-thumbnail \
  -i /workspace/source.mp4 \
  -ss 00:00:10 -vframes 1 \
  /workspace/thumbnail.jpg

# 3. Transcode to multiple bitrates
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-lite \
  -i /workspace/source.mp4 \
  -map 0:v -s 1920x1080 -b:v 5000k -c:v libx264 -preset medium /workspace/video_1080p.mp4 \
  -map 0:v -s 1280x720 -b:v 2800k -c:v libx264 -preset medium /workspace/video_720p.mp4 \
  -map 0:v -s 640x480 -b:v 1400k -c:v libx264 -preset medium /workspace/video_480p.mp4 \
  -map 0:a -c:a aac -b:a 128k /workspace/audio.m4a

# 4. Package for adaptive streaming
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/video_1080p.mp4,stream=video,output=/workspace/1080p.mp4 \
  in=/workspace/video_720p.mp4,stream=video,output=/workspace/720p.mp4 \
  in=/workspace/video_480p.mp4,stream=video,output=/workspace/480p.mp4 \
  in=/workspace/audio.m4a,stream=audio,output=/workspace/audio_packaged.m4a \
  --mpd_output /workspace/manifest.mpd \
  --hls_master_playlist_output /workspace/master.m3u8
```

---

## Use Cases

### Video Streaming Platforms
- Transcode user uploads to multiple formats
- Generate adaptive streaming manifests
- Create thumbnails and previews
- Extract metadata for cataloging

### Social Media Applications
- Process video uploads efficiently
- Generate mobile-optimized versions
- Extract video information for UI display
- Create preview thumbnails

### Educational Platforms
- Convert lecture recordings to web formats
- Create multi-quality streaming options
- Generate video metadata for search
- Optimize for bandwidth-constrained environments

### Broadcast & Media
- Professional video transcoding workflows
- DRM-protected content distribution
- Live streaming preparation
- Archive format conversion

---

## Building Locally

Each tool can be built independently:

```bash
# Build FFmpeg Lite (full-featured)
docker build -t ghcr.io/veloxpack/ffmpeg:8.0-lite ./ffmpeg-lite

# Build FFmpeg Thumbnail variant
docker build -t ghcr.io/veloxpack/ffmpeg:8.0-thumbnail ./ffmpeg-thumbnail

# Build FFmpeg Split variant
docker build -t ghcr.io/veloxpack/ffmpeg:8.0-split ./ffmpeg-split

# Build FFmpeg Concat variant
docker build -t ghcr.io/veloxpack/ffmpeg:8.0-concat ./ffmpeg-concat

# Build FFprobe
docker build -t ghcr.io/veloxpack/ffmpeg:8.0-probe ./ffprobe

# Build Shaka Packager
docker build -t shaka-packager ./shaka-packager
```

---

## Testing

### Prerequisites

- Go 1.21+ installed
- Docker running locally
- Internet connection (for downloading test containers)

### Quick Start

```bash
# 1. Setup test environment (downloads sample video)
make test-setup

# 2. Install Go dependencies
go mod download

# 3. Run all tests
make test-all
```

---

## Documentation

For detailed documentation on each tool, see the individual README files:

- [FFmpeg Lite Documentation](./ffmpeg-lite/README.md)
- [FFmpeg Thumbnail Documentation](./ffmpeg-thumbnail/README.md)
- [FFmpeg Split Documentation](./ffmpeg-split/README.md)
- [FFmpeg Concat Documentation](./ffmpeg-concat/README.md)
- [FFprobe Documentation](./ffprobe/README.md)
- [Shaka Packager Documentation](./shaka-packager/README.md)

---

## Support & Resources

- **Website**: [veloxpack.io](https://veloxpack.io)
- **GitHub Organization**: [github.com/veloxpack](https://github.com/veloxpack)
- **Docker Registry**: `ghcr.io/veloxpack/*`

---

## License

Licensed under the Apache License, Version 2.0

See [LICENSE](./LICENSE) for full license text.
