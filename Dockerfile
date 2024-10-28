# to build this docker image:
#   docker build -f Dockerfile -t schisandra-cloud-album-server .
#   docker build --build-arg OPENCV_VERSION="4.x" --build-arg OPENCV_FILE="https://github.com/opencv/opencv/archive/refs/heads/4.x.zip" --build-arg OPENCV_CONTRIB_FILE="https://github.com/opencv/opencv_contrib/archive/refs/heads/4.x.zip" -f Dockerfile -t schisandra-cloud-album-server .

FROM ubuntu:20.04 AS opencv-builder

LABEL maintainer="landaiqing <<landaiqing@126.com>>"

RUN sed -i 's|http://archive.ubuntu.com/ubuntu/|http://mirrors.aliyun.com/ubuntu/|g' /etc/apt/sources.list

ENV TZ=Europe/Madrid

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install -y --no-install-recommends --fix-missing \
      tzdata git build-essential cmake pkg-config wget unzip libgtk2.0-dev \
      curl ca-certificates libcurl4-openssl-dev libssl-dev \
      libavcodec-dev libavformat-dev libswscale-dev libtbb2 libtbb-dev \
      libharfbuzz-dev libfreetype6-dev \
      libjpeg-turbo8-dev libpng-dev libtiff-dev libdc1394-22-dev nasm && \
      rm -rf /var/lib/apt/lists/*

ARG OPENCV_VERSION="4.10.0"

ENV OPENCV_VERSION=$OPENCV_VERSION

ARG OPENCV_FILE="https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip"

ENV OPENCV_FILE=$OPENCV_FILE

ARG OPENCV_CONTRIB_FILE="https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip"

ENV OPENCV_CONTRIB_FILE=$OPENCV_CONTRIB_FILE

RUN curl -Lo opencv.zip ${OPENCV_FILE} && \
      unzip -q opencv.zip && \
      curl -Lo opencv_contrib.zip ${OPENCV_CONTRIB_FILE} && \
      unzip -q opencv_contrib.zip && \
      rm opencv.zip opencv_contrib.zip && \
      cd opencv-${OPENCV_VERSION} && \
      mkdir build && cd build && \
      cmake -D CMAKE_BUILD_TYPE=RELEASE \
      -D WITH_IPP=OFF \
      -D WITH_OPENGL=OFF \
      -D WITH_QT=OFF \
      -D WITH_FREETYPE=ON \
      -D CMAKE_INSTALL_PREFIX=/usr/local \
      -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-${OPENCV_VERSION}/modules \
      -D OPENCV_ENABLE_NONFREE=ON \
      -D WITH_JASPER=OFF \
      -D WITH_TBB=ON \
      -D BUILD_JPEG=ON \
      -D WITH_SIMD=ON \
      -D ENABLE_LIBJPEG_TURBO_SIMD=ON \
      -D BUILD_DOCS=OFF \
      -D BUILD_EXAMPLES=OFF \
      -D BUILD_TESTS=OFF \
      -D BUILD_PERF_TESTS=ON \
      -D BUILD_opencv_java=NO \
      -D BUILD_opencv_python=NO \
      -D BUILD_opencv_python2=NO \
      -D BUILD_opencv_python3=NO \
      -D OPENCV_GENERATE_PKGCONFIG=ON .. && \
      make -j $(nproc --all) && \
      make preinstall && make install && ldconfig && \
      cd / && rm -rf opencv*


FROM golang:1.23.1-alpine AS go-builder

RUN apk add --no-cache gcc musl-dev libgcc libstdc++ cmake

WORKDIR /app

COPY . .

ENV CGO_ENABLED=1

ENV CGO_CFLAGS="-I/usr/local/include/opencv4"

ENV CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_imgproc -lopencv_highgui"

ENV GOOS=linux

ENV GOARCH=amd64

ENV GOPROXY=https://goproxy.cn,direct

COPY --from=opencv-builder /usr/local/lib /usr/local/lib
COPY --from=opencv-builder /usr/local/include/opencv4 /usr/local/include/opencv4

RUN go mod download

RUN go build -o schisandra-cloud-album-server

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --no-cache tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=go-builder /app/schisandra-cloud-album-server .

COPY --from=go-builder /app/config.yaml .

COPY --from=go-builder /app/resource ./resource

COPY --from=go-builder /app/config/rbac_model.conf ./config/rbac_model.conf

COPY --from=opencv-builder /usr/local/lib /usr/local/lib

COPY --from=opencv-builder /usr/local/include/opencv4 /usr/local/include/opencv4

EXPOSE 80

CMD ["./schisandra-cloud-album-server"]
