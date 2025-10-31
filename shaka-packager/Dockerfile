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

# Stage 1: Build Shaka Packager
FROM alpine:${ALPINE_VERSION} AS shaka-packager-builder

# Install utilities, libraries, and dev tools.
RUN apk add --no-cache \
        bash bsd-compat-headers build-base cmake curl git linux-headers ninja python3 upx

# Set Shaka Packager version as an argument
ARG SHAKA_PACKAGER_VERSION=v3.4.2

WORKDIR /usr/src/

# Clone the repository with the specified tag (using --depth 1 for faster clone)
RUN git clone --recursive --depth 1 --branch "$SHAKA_PACKAGER_VERSION" https://github.com/shaka-project/shaka-packager.git

# Build Shaka Packager (now in the correct directory)
WORKDIR /usr/src/shaka-packager

# Build with fully static linking for scratch container
ENV PACKAGER_LOW_MEMORY_BUILD=yes
RUN rm -rf build
RUN cmake -S . -B build \
    -DCMAKE_BUILD_TYPE=Release \
    -DBUILD_SHARED_LIBS=OFF \
    -DFULLY_STATIC=ON \
    -G Ninja
RUN cmake --build build/ --config Release --parallel

# Strip debug symbols to reduce binary size
RUN strip /usr/src/shaka-packager/build/packager/packager

# Compress binary with UPX for maximum size reduction
RUN upx --best --lzma /usr/src/shaka-packager/build/packager/packager

# Stage 2: Minimal scratch container with static binary
FROM scratch

# Copy Shaka Packager binaries to our final image.
COPY --from=shaka-packager-builder /usr/src/shaka-packager/build/packager/packager /

# Create /tmp directory for Shaka Packager temporary files
COPY --from=shaka-packager-builder --chown=0:0 --chmod=1777 /tmp /tmp

# Set the entrypoint
ENTRYPOINT ["/packager"]

