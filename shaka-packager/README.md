# Shaka Packager

A minimal, statically-linked Shaka Packager Docker image for media packaging, encryption, and DASH/HLS manifest generation. Built from scratch for optimal performance in media streaming workflows.

## Features

- **DASH & HLS**: Generate manifests for adaptive streaming
- **DRM Support**: Common encryption (CENC) for Widevine, PlayReady, FairPlay
- **Static binary**: No runtime dependencies required
- **Multi-architecture**: Supports both `linux/amd64` and `linux/arm64`
- **Production-ready**: Industry-standard packaging tool from Google

## Use Cases

- DASH/HLS packaging for adaptive streaming
- DRM encryption and key management
- Multi-bitrate video preparation
- Live streaming packaging
- VOD content preparation
- Protected content distribution

## Image Details

- **Registry**: `ghcr.io/veloxpack/shaka-packager`
- **Base**: `scratch` (no base image)
- **Shaka Packager Version**: v3.4.2
- **Alpine Build Version**: 3.22.2

## Pull the Image

```bash
docker pull ghcr.io/veloxpack/shaka-packager:latest
```

## Usage Examples

### Basic DASH packaging

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/input.mp4,stream=audio,output=/workspace/audio.mp4 \
  in=/workspace/input.mp4,stream=video,output=/workspace/video.mp4 \
  --mpd_output /workspace/manifest.mpd
```

### HLS packaging

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/input.mp4,stream=audio,output=/workspace/audio.m4a,playlist_name=audio.m3u8 \
  in=/workspace/input.mp4,stream=video,output=/workspace/video.mp4,playlist_name=video.m3u8 \
  --hls_master_playlist_output /workspace/master.m3u8
```

### Multi-bitrate DASH packaging

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/video_1080p.mp4,stream=video,output=/workspace/video_1080p.mp4 \
  in=/workspace/video_720p.mp4,stream=video,output=/workspace/video_720p.mp4 \
  in=/workspace/video_480p.mp4,stream=video,output=/workspace/video_480p.mp4 \
  in=/workspace/audio.mp4,stream=audio,output=/workspace/audio.mp4 \
  --mpd_output /workspace/manifest.mpd
```

### Encryption with Widevine

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/input.mp4,stream=audio,output=/workspace/audio.mp4 \
  in=/workspace/input.mp4,stream=video,output=/workspace/video.mp4 \
  --enable_widevine_encryption \
  --key_server_url "https://license.usercontent.irdeto.com/..." \
  --content_id "your-content-id" \
  --signer "widevine_test" \
  --aes_signing_key "1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9" \
  --aes_signing_iv "d58ce954203b7c9a9a9d467f59839249" \
  --mpd_output /workspace/manifest.mpd
```

### HLS with AES-128 encryption

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/input.mp4,stream=audio,output=/workspace/audio.m4a,playlist_name=audio.m3u8 \
  in=/workspace/input.mp4,stream=video,output=/workspace/video.mp4,playlist_name=video.m3u8 \
  --hls_master_playlist_output /workspace/master.m3u8 \
  --enable_raw_key_encryption \
  --keys label=AUDIO:key_id=<key_id>:key=<key>,label=SD:key_id=<key_id>:key=<key> \
  --hls_key_uri "https://example.com/keys"
```

### Live profile DASH

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/input.mp4,stream=audio,output=/workspace/audio.mp4 \
  in=/workspace/input.mp4,stream=video,output=/workspace/video.mp4 \
  --mpd_output /workspace/manifest.mpd \
  --generate_static_live_mpd
```

### Fragmented MP4 for streaming

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/input.mp4,stream=audio,output=/workspace/audio_fragmented.mp4 \
  in=/workspace/input.mp4,stream=video,output=/workspace/video_fragmented.mp4 \
  --fragment_duration 2
```

## Common Options

| Option | Description |
|--------|-------------|
| `--mpd_output` | Output path for DASH manifest |
| `--hls_master_playlist_output` | Output path for HLS master playlist |
| `--fragment_duration` | Fragment duration in seconds (default: 2) |
| `--segment_duration` | Segment duration in seconds (default: same as fragment) |
| `--enable_widevine_encryption` | Enable Widevine DRM encryption |
| `--enable_raw_key_encryption` | Enable raw key encryption (for HLS AES-128) |
| `--generate_static_live_mpd` | Generate static DASH manifest for live profile |

## Complete Workflow Example

```bash
#!/bin/bash

# 1. Transcode to multiple bitrates using ffmpeg
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg \
  -i /workspace/source.mp4 \
  -map 0:v -s 1920x1080 -b:v 5000k -c:v libx264 /workspace/video_1080p.mp4 \
  -map 0:v -s 1280x720 -b:v 2800k -c:v libx264 /workspace/video_720p.mp4 \
  -map 0:v -s 640x480 -b:v 1400k -c:v libx264 /workspace/video_480p.mp4 \
  -map 0:a -c:a aac -b:a 128k /workspace/audio.mp4

# 2. Package for DASH with Shaka Packager
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/shaka-packager \
  in=/workspace/video_1080p.mp4,stream=video,output=/workspace/video_1080p_packaged.mp4 \
  in=/workspace/video_720p.mp4,stream=video,output=/workspace/video_720p_packaged.mp4 \
  in=/workspace/video_480p.mp4,stream=video,output=/workspace/video_480p_packaged.mp4 \
  in=/workspace/audio.mp4,stream=audio,output=/workspace/audio_packaged.mp4 \
  --mpd_output /workspace/manifest.mpd
```

## Building Locally

```bash
docker build -t shaka-packager ./shaka-packager
```

## Documentation

For more information, visit the [Shaka Packager documentation](https://shaka-project.github.io/shaka-packager/html/).
