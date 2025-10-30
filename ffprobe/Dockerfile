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

# Install dependencies (including upx for compression)
RUN apk add --no-cache build-base pkgconfig nasm yasm ca-certificates git cmake upx
WORKDIR /usr/src

# Build mbedTLS (much smaller than OpenSSL)
ARG MBEDTLS_VERSION=v3.6.2
RUN git clone --depth 1 --recurse-submodules --branch "$MBEDTLS_VERSION" https://github.com/ARMmbed/mbedtls.git && \
    cd mbedtls && \
    git submodule update --init && \
    cmake . \
      -DCMAKE_INSTALL_PREFIX=/usr/local \
      -DENABLE_PROGRAMS=OFF \
      -DUNSAFE_BUILD=OFF \
      -DGEN_FILES=OFF \
      -DENABLE_TESTING=OFF \
      -DBUILD_SHARED_LIBS=OFF && \
    make -j$(nproc) && \
    make install

# Define the FFmpeg version
ARG FFMPEG_VERSION=8.0

# Download and extract FFmpeg source
ADD "http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz" "ffmpeg-${FFMPEG_VERSION}.tar.gz"
RUN tar -xzf "ffmpeg-${FFMPEG_VERSION}.tar.gz"

# Go into the extracted directory
WORKDIR /usr/src/ffmpeg-$FFMPEG_VERSION

# Configure and build ffprobe only
RUN ./configure \
    --disable-ffplay \
    --disable-ffmpeg \
    --disable-doc \
    --disable-htmlpages \
    --disable-manpages \
    --disable-podpages \
    --disable-txtpages \
    --disable-encoders \
    --disable-decoders \
    --disable-hwaccels \
    --disable-muxers \
    --disable-swscale \
    --disable-swresample \
    --disable-avfilter \
    --disable-avdevice \
    --disable-debug \
    --disable-bsfs \
    --disable-parsers \
    --disable-demuxers \
    --disable-protocols \
    --disable-devices \
    --disable-filters \
    --disable-autodetect \
    --disable-runtime-cpudetect \
    --disable-swscale-alpha \
    --disable-w32threads \
    --disable-os2threads \
    --disable-dwt \
    --disable-error-resilience \
    --disable-lsp \
    --disable-faan \
    --disable-iamf \
    --disable-pixelutils \
    --disable-zlib \
    --disable-bzlib \
    --disable-lzma \
    --disable-iconv \
    --disable-xlib \
    --disable-alsa \
    --disable-sndio \
    --disable-sdl2 \
    --disable-metal \
    --disable-appkit \
    --disable-avfoundation \
    --disable-coreimage \
    --disable-audiotoolbox \
    --disable-videotoolbox \
    --disable-securetransport \
    --disable-schannel \
    --disable-d3d11va \
    --disable-d3d12va \
    --disable-dxva2 \
    --disable-vaapi \
    --disable-vdpau \
    --disable-vulkan \
    --disable-cuda-llvm \
    --disable-cuvid \
    --disable-nvdec \
    --disable-nvenc \
    --disable-v4l2-m2m \
    --disable-libdrm \
    --disable-amf \
    --disable-ffnvcodec \
    --disable-symver \
    --disable-version-tracking \
    --disable-safe-bitstream-reader \
    --disable-logging \
    --enable-demuxer=mov,mp4,mpegts,matroska,flv \
    --enable-parser=h264,hevc \
    --enable-protocol=file,http,https,rtmp,rtsp,udp \
    --enable-mbedtls \
    --enable-gpl \
    --enable-small \
    --enable-version3 \
    --enable-nonfree \
    --enable-lto \
    --enable-static \
    --disable-shared \
    --extra-cflags="-static -Oz -flto -ffunction-sections -fdata-sections -fno-unwind-tables -fno-asynchronous-unwind-tables" \
    --extra-ldflags="-static -flto -Wl,--gc-sections -Wl,-s -Wl,--strip-all" \
    --prefix=/usr/local && \
    make -j$(nproc) && \
    make install && \
    strip --strip-all --remove-section=.comment --remove-section=.note /usr/local/bin/ffprobe && \
    upx --best --ultra-brute --lzma /usr/local/bin/ffprobe

# Stage 2: Final container with tools available
FROM scratch

# Copy ffprobe to the final image
COPY --from=ffprobe-builder /usr/local/bin/ffprobe /

# Set the entrypoint
ENTRYPOINT ["/ffprobe"]
