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

RUN apk update && apk upgrade && \
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
    yasm

WORKDIR /usr/src

# MbedTLS - Lightweight cryptographic library for secure communications.
ARG MBEDTLS_VERSION=v3.4.1
RUN git clone --depth 1 --branch "$MBEDTLS_VERSION" https://github.com/ARMmbed/mbedtls.git && \
    cd mbedtls && \
    sed -e 's/-Wdocumentation//' -e 's/-Wno-documentation-deprecated-sync//' -i.bk library/CMakeLists.txt && \
    cmake . \
      -DCMAKE_INSTALL_PREFIX=/usr/local \
      -DENABLE_PROGRAMS=OFF \
      -DUNSAFE_BUILD=OFF \
      -DGEN_FILES=OFF \
      -DENABLE_TESTING=OFF && \
    make -j$(nproc) && \
    make install

# FFmpeg - A complete, cross-platform solution to record, convert, and stream audio and video.
ARG FFMPEG_VERSION=8.0
ADD "http://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz" "ffmpeg-${FFMPEG_VERSION}.tar.gz"
RUN tar -xzf "ffmpeg-${FFMPEG_VERSION}.tar.gz"

# Go into the extracted directory
WORKDIR /usr/src/ffmpeg-$FFMPEG_VERSION

# Configure and build ffmpeg only
RUN ./configure \
    --pkg-config-flags="--static" \
    --disable-ffplay \
    --disable-debug \
    --disable-asm \
    --disable-doc \
    --disable-hwaccels \
    --disable-vaapi \
    --enable-mbedtls \
    --enable-runtime-cpudetect \
    --enable-gpl \
    --enable-version3 \
    --enable-static \
    --disable-shared \
    --extra-cflags="-static" \
    --extra-ldflags="-static" \
    --prefix=/usr/local && \
    make -j$(nproc) && \
    make install

# Stage 2: Final container with tools available
FROM scratch

# Copy ffmpeg
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /

# Copy ffprobe
COPY --from=ffmpeg-builder /usr/local/bin/ffprobe /

# Set the entrypoint
ENTRYPOINT ["/ffmpeg"]
# ENTRYPOINT ["/ffprobe"]
