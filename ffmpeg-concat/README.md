# FFmpeg Concat

Ultra-lightweight FFmpeg Docker image optimized specifically for concatenating videos without re-encoding. Uses stream copying for fast, lossless video concatenation.

## Features

- **Ultra-Lightweight**: Minimal binary size (estimated ~1.5-2 MB)
- **Stream Copy**: No re-encoding, preserves original quality
- **Fast**: Concatenate videos in seconds without transcoding
- **Multiple Methods**: Concat demuxer, concat protocol, and concat filter
- **Advanced File Lists**: Support for duration, inpoint/outpoint, trimming, stream selection
- **Format Support**: MP4, MOV, Matroska (MKV/WebM), MPEGTS
- **Protocol Support**: File, HTTP, HTTPS, concat
- **Static binary**: No runtime dependencies

## Image Details

- **Registry**: `ghcr.io/veloxpack/ffmpeg-concat`
- **Base**: `scratch` (no base image)
- **FFmpeg Version**: 8.0
- **Alpine Build Version**: 3.22.2
- **SSL/TLS**: mbedTLS
- **Compression**: UPX with LZMA
- **Build Optimizations**: LTO, -Oz, aggressive stripping

## Use Cases

- Merge multiple video files
- Join video segments
- Combine video clips with trimming (inpoint/outpoint)
- Stitch videos together
- Create compilations without re-encoding
- Concatenate with duration metadata
- Stream-selective concatenation

## Pull the Image

```bash
docker pull ghcr.io/veloxpack/ffmpeg-concat:latest
```

## Usage Examples

### Method 1: Concat Demuxer (Recommended - Most Reliable)

Create a text file listing videos to concatenate:

```bash
# Create concat list file with absolute paths
cat > /workspace/filelist.txt << EOF
file '/workspace/video1.mp4'
file '/workspace/video2.mp4'
file '/workspace/video3.mp4'
EOF

# Concatenate
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -f concat \
  -safe 0 \
  -i /workspace/filelist.txt \
  -c copy \
  /workspace/output.mp4
```

**Or with relative paths:**

```bash
# Create concat list file with relative paths
cat > /workspace/list.txt << EOF
file 'video1.mp4'
file 'video2.mp4'
file 'video3.mp4'
EOF

docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -f concat \
  -safe 0 \
  -i /workspace/list.txt \
  -c copy \
  /workspace/output.mp4
```

### Method 2: Concat Protocol (Simple, Same Format Only)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -i "concat:/workspace/video1.mp4|/workspace/video2.mp4|/workspace/video3.mp4" \
  -c copy \
  /workspace/output.mp4
```

### Method 3: Concat Filter (Most Flexible)

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -i /workspace/video1.mp4 \
  -i /workspace/video2.mp4 \
  -i /workspace/video3.mp4 \
  -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a]concat=n=3:v=1:a=1[outv][outa]" \
  -map "[outv]" \
  -map "[outa]" \
  /workspace/output.mp4
```

### Concatenate with absolute paths

```bash
# list.txt with absolute paths
cat > /workspace/list.txt << EOF
file '/workspace/clip1.mp4'
file '/workspace/clip2.mp4'
file '/workspace/clip3.mp4'
EOF

docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -f concat \
  -safe 0 \
  -i /workspace/list.txt \
  -c copy \
  /workspace/merged.mp4
```

### Concatenate MPEGTS files directly

```bash
# MPEGTS can be concatenated with simple cat
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -i "concat:/workspace/part1.ts|/workspace/part2.ts|/workspace/part3.ts" \
  -c copy \
  -bsf:a aac_adtstoasc \
  /workspace/output.mp4
```

### Concatenate with timestamps preserved

```bash
docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -f concat \
  -safe 0 \
  -i /workspace/list.txt \
  -c copy \
  -copyts \
  /workspace/output.mp4
```

### Concatenate WebM files

```bash
# list.txt
cat > /workspace/list.txt << EOF
file 'video1.webm'
file 'video2.webm'
file 'video3.webm'
EOF

docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -f concat \
  -safe 0 \
  -i /workspace/list.txt \
  -c copy \
  /workspace/output.webm
```

### Concatenate remote videos

```bash
cat > /workspace/list.txt << EOF
file 'https://example.com/video1.mp4'
file 'https://example.com/video2.mp4'
file 'https://example.com/video3.mp4'
EOF

docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -f concat \
  -safe 0 \
  -protocol_whitelist file,http,https,tcp,tls \
  -i /workspace/list.txt \
  -c copy \
  /workspace/output.mp4
```

## FFmpeg File List Formats

The concat demuxer supports multiple file list formats for advanced use cases:

### 1. Standard Format (Most Common)
```bash
cat > /workspace/list.txt << EOF
file '/workspace/video1.mp4'
file '/workspace/video2.mp4'
file '/workspace/video3.mp4'
EOF
```

### 2. With Duration Metadata
```bash
cat > /workspace/list.txt << EOF
file '/workspace/video1.mp4'
duration 10.5
file '/workspace/video2.mp4'
duration 15.2
file '/workspace/video3.mp4'
duration 8.7
EOF
```

### 3. With In/Out Points (Trimming)
```bash
cat > /workspace/list.txt << EOF
file '/workspace/video1.mp4'
inpoint 5.0
outpoint 15.0

file '/workspace/video2.mp4'
inpoint 0.0
outpoint 10.5

file '/workspace/video3.mp4'
EOF

# inpoint - Start time in seconds
# outpoint - End time in seconds
```

**Example: Trim and concatenate**
```bash
cat > /workspace/trimlist.txt << EOF
file '/workspace/long_video1.mp4'
inpoint 30.0
outpoint 60.0

file '/workspace/long_video2.mp4'
inpoint 10.0
outpoint 45.0
EOF

docker run --rm -v $(pwd):/workspace \
  ghcr.io/veloxpack/ffmpeg-concat \
  -f concat \
  -safe 0 \
  -i /workspace/trimlist.txt \
  -c copy \
  /workspace/trimmed_concat.mp4
```

### 4. Relative Paths
```bash
# When files are in the same directory as the list file
cat > /workspace/list.txt << EOF
file 'video1.mp4'
file 'video2.mp4'
file 'video3.mp4'
EOF
```

### 5. With File Options
```bash
cat > /workspace/list.txt << EOF
file '/workspace/video1.mp4'
file_packet_metadata key=value

file '/workspace/video2.mp4'
file_packet_metadata title=MyVideo
EOF
```

### 6. Stream Selection
```bash
cat > /workspace/list.txt << EOF
file '/workspace/video1.mp4'
stream
stream 0

file '/workspace/video2.mp4'
stream 1
EOF
```

## Important Notes

### Concat Demuxer Requirements
- All videos must have the same codec
- All videos must have the same resolution
- All videos must have the same frame rate
- For best results, videos should have the same encoding parameters

### When Videos Don't Match
If videos have different properties, you'll need to re-encode (use main ffmpeg image):
```bash
# This image won't work for different formats - use ghcr.io/veloxpack/ffmpeg instead
```

### Format-Specific Tips

**MP4/MOV Files:**
```bash
-c copy -movflags +faststart
```

**MPEGTS Files:**
```bash
-c copy -bsf:a aac_adtstoasc
```

**WebM/Matroska:**
```bash
-c copy
```

## Troubleshooting

### "Non-monotonous DTS in output stream"
Add `-fflags +genpts`:
```bash
-f concat -safe 0 -i list.txt -c copy -fflags +genpts output.mp4
```

### "Timestamps are unset in a packet"
Add `-avoid_negative_ts make_zero`:
```bash
-f concat -safe 0 -i list.txt -c copy -avoid_negative_ts make_zero output.mp4
```

### "Unsafe file name"
Use `-safe 0` or provide only filenames (not paths):
```bash
-f concat -safe 0 -i list.txt
```

## Limitations

### Not Included
- Video encoding/transcoding
- Audio encoding/transcoding
- Video filters (except concat filter)
- Format conversion with re-encoding
- Hardware acceleration

### What This Image IS For
✅ Fast video concatenation
✅ Merging same-format videos
✅ Joining video segments
✅ Lossless video stitching

### What This Image IS NOT For
❌ Concatenating different formats (requires re-encoding)
❌ Video transcoding
❌ Format conversion
❌ Video filtering/effects

## Performance Tips

- Always use `-c copy` for stream copying (no re-encoding)
- Use the concat demuxer method for most reliable results
- Ensure all input videos have matching properties
- Use `-safe 0` when working with absolute paths
- Consider using MPEGTS intermediate format for better compatibility

## Building Locally

```bash
docker build -t ffmpeg-concat ./ffmpeg-concat
```

