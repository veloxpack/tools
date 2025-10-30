# Copyright 2025 Veloxpack.io
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Define build arguments
ARG ALPINE_VERSION=3.22.2

# Stage 1: Build ffprobe (using a specific version)
FROM alpine:${ALPINE_VERSION} AS ffprobe-builder

# Install dependencies
RUN apk add --no-cache build-base openssl-dev pkgconfig yasm ca-certificates openssl-libs-static
WORKDIR /usr/src/ffmpeg

# Define the FFmpeg version
ARG FFMPEG_VERSION=8.0

# Download and extract FFmpeg source
ADD "http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz" "ffmpeg-${FFMPEG_VERSION}.tar.gz"
RUN tar -xzf "ffmpeg-${FFMPEG_VERSION}.tar.gz"

# Go into the extracted directory
WORKDIR /usr/src/ffmpeg/ffmpeg-$FFMPEG_VERSION

# Configure and build ffprobe only
RUN ./configure \
    --disable-ffplay \
    --disable-ffmpeg \
    --disable-doc \
    --disable-muxers \
    --disable-swscale \
    --disable-avfilter \
    --disable-avdevice \
    --disable-debug \
    --disable-bsfs \
    --disable-parsers \
    --disable-demuxers \
    --enable-demuxer=flv,mov,mp4,m4a,3gp,3g2,mj2,mpegts,mpegvideo,matroska,ogg,avi \
    --enable-parser=h264,hevc,aac,mpeg4video,mpegaudio \
    --enable-protocol=https \
    --enable-openssl \
    --enable-gpl \
    --enable-small \
    --enable-version3 \
    --enable-nonfree \
    --enable-static \
    --disable-shared \
    --extra-cflags="-static" \
    --extra-ldflags="-static" \
    --prefix=/usr/local && \
    make -j$(nproc) && \
    make install

# Stage 2: Final container with tools available
FROM scratch

# Copy ffprobe to the final image
COPY --from=ffprobe-builder /usr/local/bin/ffprobe /

# Set the entrypoint
ENTRYPOINT ["/ffprobe"]
