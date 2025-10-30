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
        bash bsd-compat-headers build-base cmake curl git linux-headers ninja python3

# Set Shaka Packager version as an argument
ARG SHAKA_PACKAGER_VERSION=v3.4.2

WORKDIR /usr/src/

# Clone the repository with the specified tag (using --depth 1 for faster clone)
RUN git clone --recursive --depth 1 --branch "$SHAKA_PACKAGER_VERSION" https://github.com/shaka-project/shaka-packager.git

# Build Shaka Packager (now in the correct directory)
WORKDIR /usr/src/shaka-packager

RUN rm -rf build # Ensure clean build
RUN cmake -S . -B build -DCMAKE_BUILD_TYPE=Release -G Ninja
RUN cmake --build build/ --config Release --parallel

# Stage 2: Final container with tools available
# FROM alpine:${ALPINE_VERSION}

# Install runtime dependencies (libstdc++ and python3)
# RUN apk add --no-cache libstdc++ python3

# Copy Shaka Packager binaries to our final image.
# COPY --from=shaka-packager-builder /usr/src/shaka-packager/build/packager/packager /usr/local/bin/
# COPY --from=shaka-packager-builder /usr/src/shaka-packager/build/packager/mpd_generator /usr/local/bin/
# COPY --from=shaka-packager-builder /usr/src/shaka-packager/build/packager/pssh-box.py /usr/local/bin/

# Copy pyproto directory, which is needed by pssh-box.py script. This line
# cannot be combined with the line above as Docker's copy command skips the
# directory itself. See https://github.com/moby/moby/issues/15858 for details.
# COPY --from=shaka-packager-builder /usr/src/shaka-packager/build/packager/pssh-box-protos /usr/local/bin/pssh-box-protos

# Set the entrypoint
# ENTRYPOINT ["packager"]

# Stage 2: Final container with tools available
FROM scratch

# Copy Shaka Packager binaries to our final image.
COPY --from=shaka-packager-builder /usr/src/shaka-packager/build/packager/packager /

# Set the entrypoint
ENTRYPOINT ["/packager"]
