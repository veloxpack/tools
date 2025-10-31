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

# Stage 1: Build ffmpeg for video splitting
FROM alpine:${ALPINE_VERSION} AS ffmpeg-builder

# Install dependencies (including upx for compression)
RUN apk add --no-cache build-base pkgconfig nasm yasm ca-certificates git cmake upx

WORKDIR /usr/src

# Build mbedTLS (for HTTPS support)
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

# x264 - H.264 encoder (needed for scene detection re-encoding)
ARG X264_VERSION="0480cb05fa188d37ae87e8f4fd8f1aea3711f7ee"
RUN git clone --depth 1 https://code.videolan.org/videolan/x264.git && \
    cd x264 && \
    git checkout $X264_VERSION && \
    sh ./configure --disable-opencl --enable-static && \
    make -j$(nproc) && \
    make install

# FFmpeg - Video splitting optimized build
ARG FFMPEG_VERSION=8.0
ADD "http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz" "ffmpeg-${FFMPEG_VERSION}.tar.gz"
RUN tar -xzf "ffmpeg-${FFMPEG_VERSION}.tar.gz"

WORKDIR /usr/src/ffmpeg-$FFMPEG_VERSION

# Configure and build ffmpeg for video splitting (stream copy, no encoding)
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
    --enable-encoder=copy \
    --enable-encoder=libx264 \
    --enable-encoder=aac \
    --disable-decoders \
    --enable-decoder=h264 \
    --enable-decoder=hevc \
    --enable-decoder=vp8 \
    --enable-decoder=vp9 \
    --enable-decoder=aac \
    --enable-decoder=mp3 \
    --disable-hwaccels \
    --disable-muxers \
    --enable-muxer=mp4 \
    --enable-muxer=mov \
    --enable-muxer=matroska \
    --enable-muxer=webm \
    --enable-muxer=mpegts \
    --enable-muxer=segment \
    --disable-demuxers \
    --enable-demuxer=mov \
    --enable-demuxer=mp4 \
    --enable-demuxer=matroska \
    --enable-demuxer=mpegts \
    --disable-parsers \
    --enable-parser=h264 \
    --enable-parser=hevc \
    --enable-parser=vp8 \
    --enable-parser=vp9 \
    --enable-parser=aac \
    --enable-parser=mpegaudio \
    --disable-bsfs \
    --enable-bsf=h264_mp4toannexb \
    --enable-bsf=hevc_mp4toannexb \
    --enable-bsf=aac_adtstoasc \
    --disable-protocols \
    --enable-protocol=file \
    --enable-protocol=http \
    --enable-protocol=https \
    --disable-devices \
    --disable-filters \
    --enable-filter=select \
    --enable-filter=showinfo \
    --enable-filter=setpts \
    --enable-filter=scale \
    --enable-filter=metadata \
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
    --enable-libx264 \
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
    strip --strip-all --remove-section=.comment --remove-section=.note /usr/local/bin/ffmpeg && \
    upx --best --ultra-brute --lzma /usr/local/bin/ffmpeg

# Stage 2: Final container
FROM scratch

# Copy ffmpeg
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /

# Set the entrypoint
ENTRYPOINT ["/ffmpeg"]

