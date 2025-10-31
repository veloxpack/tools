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

RUN apk update && \
    apk upgrade && \
    apk add \
    cmake \
    curl \
    diffutils \
    g++ \
    git \
    libvdpau-dev \
    linux-headers \
    make \
    nasm \
    patch \
    perl \
    pkgconfig \
    yasm \
    upx

WORKDIR /usr/src

# LAME - High-quality MP3 audio encoder.
ARG LAME_VERSION=3.100
ADD "https://sourceforge.net/projects/lame/files/lame/$LAME_VERSION/lame-$LAME_VERSION.tar.gz/download" "lame-$LAME_VERSION.tar.gz"
RUN tar -xzf "lame-$LAME_VERSION.tar.gz" && \
    cd "lame-$LAME_VERSION" && \
    ./configure \
        --prefix=/usr/local \
        --disable-frontend \
        --enable-static \
        --disable-shared && \
    make -j$(nproc) && \
    make install

# Opus - A versatile audio codec for interactive speech and music transmission.
ARG OPUS_VERSION=v1.4
RUN git clone --depth 1 --branch "$OPUS_VERSION" https://github.com/xiph/opus.git && \
    cd opus && \
    cmake . \
        -DCMAKE_INSTALL_PREFIX=/usr/local \
        -DOPUS_BUILD_SHARED_LIBRARY=OFF \
        -DOPUS_BUILD_FRAMEWORK=OFF \
        -DBUILD_SHARED_LIBS=OFF \
        -DOPUS_BUILD_TESTING=OFF \
        -DBUILD_TESTING=OFF \
        -DOPUS_BUILD_PROGRAMS=OFF && \
    make -j$(nproc) && \
    make install

# SVT-AV1 - Scalable Video Technology for AV1, an AV1 encoder from Intel.
ENV SVT_AV1_VERSION=v1.7.0
RUN git clone --depth 1 --branch "$SVT_AV1_VERSION" https://gitlab.com/AOMediaCodec/SVT-AV1.git && \
    mkdir SVT-AV1-build && \
    cd SVT-AV1-build && \
    cmake ../SVT-AV1 \
      -DCMAKE_INSTALL_PREFIX=/usr/local \
      -DBUILD_SHARED_LIBS=OFF \
      -DBUILD_TESTING=OFF \
      -DCOVERAGE=OFF \
      -DBUILD_APPS=OFF \
      -DREPRODUCIBLE_BUILDS=ON && \
    make -j$(nproc) && \
    make install

# libvpx - VP8/VP9 video codec library from Google.
ARG LIBVPX_VERSION=v1.13.0
RUN git clone --depth 1 --branch "$LIBVPX_VERSION" https://chromium.googlesource.com/webm/libvpx && \
    cd libvpx && \
    ./configure \
      --enable-vp8 \
      --enable-vp9 \
      --disable-unit-tests \
      --disable-examples \
      --enable-static \
      --disable-shared && \
    make -j$(nproc) && \
    make install

# x264 - Open-source H.264 encoder.
ARG X264_VERSION="0480cb05fa188d37ae87e8f4fd8f1aea3711f7ee"
RUN git clone --depth 1 https://code.videolan.org/videolan/x264.git && \
    cd x264 && \
    git checkout $X264_VERSION && \
    sh ./configure --disable-opencl --enable-static && \
    make -j$(nproc) && \
    make install

# x265 - Open-source H.265/HEVC encoder.
ARG X265_VERSION=3.5
RUN git clone --depth 1 --branch "$X265_VERSION" https://bitbucket.org/multicoreware/x265_git.git && \
    cd x265_git/build && \
    cmake ../source \
      -DCMAKE_INSTALL_PREFIX=/usr/local \
      -DENABLE_SHARED=OFF \
      -DENABLE_CLI=OFF && \
    make -j$(nproc) && \
    make install && \
    # This adjustment to the x265 linker flags is needed, at least on
    # arm, to successfully link against it statically.  (-lgcc_s not
    # found (or needed), and -lpthread missing)
    sed -e 's/-lgcc_s -lgcc -lgcc_s -lgcc/-lpthread -lgcc/' -i.bk /usr/local/lib/pkgconfig/x265.pc

# Define the FFmpeg version
ARG FFMPEG_VERSION=8.0

# Download and extract FFmpeg source
ADD "http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz" "ffmpeg-${FFMPEG_VERSION}.tar.gz"
RUN tar -xzf "ffmpeg-${FFMPEG_VERSION}.tar.gz"

# Go into the extracted directory
WORKDIR /usr/src/ffmpeg-$FFMPEG_VERSION

# Configure and build ffmpeg for video encoding/transcoding (no image generation)
RUN ./configure \
    --pkg-config-flags="--static" \
    --disable-ffplay \
    --disable-ffprobe \
    --disable-doc \
    --disable-htmlpages \
    --disable-manpages \
    --disable-podpages \
    --disable-txtpages \
    --disable-debug \
    --disable-encoders \
    --enable-encoder=libx264 \
    --enable-encoder=libx265 \
    --enable-encoder=libvpx_vp8 \
    --enable-encoder=libvpx_vp9 \
    --enable-encoder=libsvtav1 \
    --enable-encoder=libmp3lame \
    --enable-encoder=libopus \
    --enable-encoder=aac \
    --disable-hwaccels \
    --disable-vaapi \
    --disable-muxers \
    --enable-muxer=mp4 \
    --enable-muxer=mov \
    --enable-muxer=matroska \
    --enable-muxer=webm \
    --enable-muxer=mpegts \
    --enable-muxer=mp3 \
    --enable-muxer=ogg \
    --enable-muxer=flv \
    --disable-protocols \
    --enable-protocol=file \
    --enable-protocol=rtmp \
    --enable-protocol=rtsp \
    --enable-protocol=udp \
    --disable-autodetect \
    --disable-swscale-alpha \
    --disable-iamf \
    --disable-pixelutils \
    --enable-libvpx \
    --enable-libsvtav1 \
    --enable-libx264 \
    --enable-libx265 \
    --enable-libmp3lame \
    --enable-libopus \
    --enable-gpl \
    --enable-small \
    --enable-version3 \
    --enable-lto \
    --enable-static \
    --disable-shared \
    --extra-cflags="-static -Oz -flto -ffunction-sections -fdata-sections" \
    --extra-ldflags="-static -flto -Wl,--gc-sections -Wl,-s" \
    --prefix=/usr/local && \
    make -j$(nproc) && \
    make install && \
    strip --strip-all --remove-section=.comment --remove-section=.note /usr/local/bin/ffmpeg && \
    upx --best --lzma /usr/local/bin/ffmpeg

# Stage 2: Final container with tools available
FROM scratch

# Copy ffmpeg to the final image
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /

# Set the entrypoint
ENTRYPOINT ["/ffmpeg"]
