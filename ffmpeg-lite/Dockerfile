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

# Stage 1: Build ffmpeg (using a specific version)
FROM alpine:${ALPINE_VERSION} AS ffmpeg-builder

# Install dependencies (including upx for compression and zlib for PNG encoding)
RUN apk add --no-cache build-base pkgconfig nasm yasm ca-certificates git cmake upx zlib-dev zlib-static

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

# FFmpeg - A complete, cross-platform solution to record, convert, and stream audio and video.
ARG FFMPEG_VERSION=8.0
ADD "http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz" "ffmpeg-${FFMPEG_VERSION}.tar.gz"
RUN tar -xzf "ffmpeg-${FFMPEG_VERSION}.tar.gz"

# Go into the extracted directory
WORKDIR /usr/src/ffmpeg-$FFMPEG_VERSION

# Configure and build ffmpeg for thumbnail generation only
RUN ./configure \
    --disable-ffplay \
    --disable-ffprobe \
    --disable-doc \
    --disable-htmlpages \
    --disable-manpages \
    --disable-podpages \
    --disable-txtpages \
    --disable-debug \
    --disable-encoders \
    --enable-encoder=png \
    --enable-encoder=mjpeg \
    --disable-decoders \
    --enable-decoder=h264 \
    --enable-decoder=hevc \
    --enable-decoder=vp8 \
    --enable-decoder=vp9 \
    --enable-decoder=av1 \
    --enable-decoder=mpeg4 \
    --enable-decoder=mpeg2video \
    --enable-decoder=mpeg1video \
    --enable-decoder=aac \
    --enable-decoder=aac_latm \
    --enable-decoder=mp3 \
    --enable-decoder=mjpeg \
    --disable-hwaccels \
    --disable-muxers \
    --enable-muxer=image2 \
    --enable-muxer=image2pipe \
    --enable-muxer=mjpeg \
    --disable-demuxers \
    --enable-demuxer=mov \
    --enable-demuxer=mp4 \
    --enable-demuxer=m4a \
    --enable-demuxer=3gp \
    --enable-demuxer=matroska \
    --enable-demuxer=mpegts \
    --enable-demuxer=mpegvideo \
    --enable-demuxer=flv \
    --enable-demuxer=avi \
    --disable-parsers \
    --enable-parser=h264 \
    --enable-parser=hevc \
    --enable-parser=mpeg4video \
    --enable-parser=mpegvideo \
    --enable-parser=vp8 \
    --enable-parser=vp9 \
    --enable-parser=aac \
    --enable-parser=aac_latm \
    --enable-parser=mpegaudio \
    --enable-parser=mjpeg \
    --disable-bsfs \
    --enable-bsf=h264_mp4toannexb \
    --enable-bsf=hevc_mp4toannexb \
    --enable-bsf=extract_extradata \
    --disable-protocols \
    --enable-protocol=file \
    --enable-protocol=pipe \
    --enable-protocol=http \
    --enable-protocol=https \
    --enable-protocol=rtmp \
    --enable-protocol=rtsp \
    --enable-protocol=udp \
    --disable-devices \
    --disable-filters \
    --enable-filter=scale \
    --enable-filter=thumbnail \
    --enable-filter=setpts \
    --enable-filter=fps \
    --enable-filter=select \
    --enable-filter=null \
    --enable-filter=format \
    --enable-filter=tile \
    --enable-filter=pad \
    --enable-filter=crop \
    --enable-filter=hflip \
    --enable-filter=vflip \
    --enable-filter=transpose \
    --disable-swresample \
    --enable-swscale \
    --enable-avfilter \
    --disable-avdevice \
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
    --enable-mbedtls \
    --enable-gpl \
    --enable-small \
    --enable-version3 \
    --enable-nonfree \
    --enable-lto \
    --enable-static \
    --disable-shared \
    --extra-cflags="-static -Oz -flto -ffunction-sections -fdata-sections" \
    --extra-ldflags="-static -flto -Wl,--gc-sections -Wl,-s" \
    --prefix=/usr/local && \
    make -j$(nproc) && \
    make install && \
    strip --strip-all /usr/local/bin/ffmpeg && \
    upx --best --lzma /usr/local/bin/ffmpeg

# Stage 2: Final container with tools available
FROM scratch

# Copy ffmpeg
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /

# Set the entrypoint
ENTRYPOINT ["/ffmpeg"]
