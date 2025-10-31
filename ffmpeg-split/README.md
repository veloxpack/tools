# FFmpeg Split

Ultra-lightweight FFmpeg Docker image optimized specifically for splitting videos without re-encoding. Uses stream copying for fast, lossless video splitting.

## Features

- **Ultra-Lightweight**: Only 3.92 MB compressed image size
- **Stream Copy**: No re-encoding for fast splits
- **Scene Detection**: Automatically split on scene changes
- **Metadata Export**: Export scene timestamps for automation
- **H.264 Encoding**: Re-encode when needed (scene detection)
- **Fast**: Split videos quickly with or without transcoding
- **Format Support**: MP4, MOV, Matroska (MKV/WebM), MPEGTS
- **Protocol Support**: File (local files only)
- **Static binary**: No runtime dependencies

## Image Details

- **Registry**: `ghcr.io/veloxpack/ffmpeg:8.0-split`
- **Base**: `scratch` (no base image)
- **Image Size**: 3.92 MB (compressed)
- **FFmpeg Version**: 8.0
- **Alpine Build Version**: 3.22.2
- **Compression**: UPX with LZMA
- **Build Optimizations**: LTO, -Oz, aggressive stripping

## Included Components

### Video Encoders
- Copy (stream copy, no re-encoding)
- libx264 (H.264 encoder for scene detection)
- AAC (audio encoder)

### Video Decoders
- H.264, HEVC, VP8, VP9
- AAC, MP3

### Filters
- `select` - Scene detection and frame selection
- `metadata` - Export scene metadata to files (for automation)
- `showinfo` - Display frame information
- `setpts` - Set presentation timestamps
- `scale` - Resize frames

### Protocols
- File (local files only)

## Use Cases

- Split videos into segments
- Extract portions of videos
- Automatic scene detection and splitting
- Create clips without re-encoding
- Time-based video trimming
- Fast video editing workflows
- Content-aware video segmentation

## Pull the Image

```bash
docker pull ghcr.io/veloxpack/ffmpeg:8.0-split
```

## Usage Examples

### Split video into 60-second segments

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -c copy \
  -f segment \
  -segment_time 60 \
  -reset_timestamps 1 \
  /workspace/output-%03d.mp4
```

### Extract a specific time range (no re-encoding)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -ss 00:01:30 \
  -to 00:05:00 \
  -c copy \
  /workspace/clip.mp4
```

### Split video at exact timestamps

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -c copy \
  -f segment \
  -segment_times 30,60,120,180 \
  /workspace/segment-%d.mp4
```

### Split by file size (e.g., 10MB chunks)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -c copy \
  -f segment \
  -segment_size 10485760 \
  /workspace/chunk-%03d.mp4
```

### Extract first 30 seconds

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -t 30 \
  -c copy \
  /workspace/first-30s.mp4
```

### Split with segment list output

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -c copy \
  -f segment \
  -segment_time 300 \
  -segment_list /workspace/segments.txt \
  /workspace/part-%03d.mp4
```

### Split on scene changes (automatic scene detection)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -vf "select='gt(scene,0.4)',showinfo" \
  -vsync vfr \
  /workspace/scene_%03d.mp4
```

### Scene detection with custom threshold

```bash
# Lower threshold (0.3) = more sensitive, more scenes detected
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -vf "select='gt(scene,0.3)'" \
  -vsync vfr \
  /workspace/scene_%03d.mp4

# Higher threshold (0.5) = less sensitive, fewer scenes
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -vf "select='gt(scene,0.5)'" \
  -vsync vfr \
  /workspace/scene_%03d.mp4
```

### Scene detection with H.264 encoding

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -vf "select='gt(scene,0.4)'" \
  -vsync vfr \
  -c:v libx264 -preset fast -crf 23 \
  -c:a aac \
  /workspace/scene_%03d.mp4
```

### Scene detection with metadata output

```bash
# Export scene metadata to a file for later processing
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -vf "select='gt(scene,0.4)',metadata=print:file=/workspace/scenes.txt" \
  -vsync vfr \
  /workspace/scene_%03d.mp4
```

### Scene detection with custom threshold and metadata

```bash
# Programmatic scene detection (useful for automation)
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg:8.0-split \
  -i /workspace/input.mp4 \
  -vf "select='gt(scene,0.35)',metadata=print:file=/workspace/scene_timestamps.txt" \
  -vsync vfr \
  -c:v libx264 -preset fast -crf 23 \
  -c:a aac \
  /workspace/scene_%03d.mp4

# The scene_timestamps.txt will contain frame metadata for analysis
```

## Advanced Options

### Reset timestamps for each segment
```bash
-reset_timestamps 1
```

### Split at keyframes only (faster, may not be exact)
```bash
-f segment -segment_time 60 -segment_format mp4
```

### Split with frame-accurate cuts (slower)
```bash
-ss 00:01:30 -to 00:05:00 -c copy -avoid_negative_ts make_zero
```

## Limitations

### Not Included
- Advanced video encoders (only libx264 and copy)
- Advanced audio encoders (only AAC and copy)
- Complex filters (only select, showinfo, setpts, scale, metadata)
- Thumbnail generation
- Hardware acceleration
- Network protocols (HTTP, HTTPS, RTMP) - local files only

### What This Image IS For
Fast video splitting (stream copy)
Scene-based video splitting
Extracting video clips
Segmenting videos
Lossless trimming
Content-aware splitting

### What This Image IS NOT For
❌ Advanced video transcoding
❌ Complex audio processing
❌ Live streaming
❌ Multi-format encoding
❌ Advanced video effects

## Performance Tips

- Always use `-c copy` for stream copying (no re-encoding)
- Use `-avoid_negative_ts make_zero` if you encounter timestamp issues
- For exact cuts at specific times, input seeking (`-ss` before `-i`) is faster
- For frame-accurate cuts, output seeking (`-ss` after `-i`) is more precise

## Building Locally

```bash
docker build -t ghcr.io/veloxpack/ffmpeg:8.0-split ./ffmpeg-split
```

