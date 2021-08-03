#FROM ubuntu:latest
FROM debian:stable-slim
RUN apt-get update -y
RUN apt-get install gcc -y
RUN apt-get install nano -y
RUN apt-get install tar -y
RUN apt-get install bash -y
RUN apt-get install sudo -y
RUN apt-get install openssl -y
RUN apt-get install git -y
RUN apt-get install make -y
RUN apt-get install cmake -y
RUN apt-get install gfortran -y
RUN apt-get install pkg-config -y
RUN apt-get install wget -y
RUN apt-get install curl -y
RUN apt-get install unzip -y
RUN apt-get install net-tools -y
RUN apt-get install iproute2 -y
RUN apt-get install iputils-ping -y
RUN apt-get install tzdata -y
RUN apt-get install bc -y
# Programming Languages , Python Comes Pre-Packaged Now
# RUN apt-get install golang-go -y
# Apparently, gocv.io requires go 1.11.6
# and 11JUN2021 , apt-get installs
# we installed a specific version of go somewhere in a docker
# https://golang.org/dl/

RUN apt-get install python-pip -y
RUN apt-get install python3-pip -y
RUN apt-get install python3-venv -y
RUN apt-get install build-essential -y
RUN apt-get install python3-dev -y
RUN apt-get install python3-setuptools -y
RUN apt-get install python3-smbus -y
RUN apt-get install python3-numpy -y
RUN apt-get install python3-scipy -y
RUN apt-get install libncursesw5-dev -y
RUN apt-get install libgdbm-dev -y
RUN apt-get install libc6-dev -y
RUN apt-get install zlib1g-dev -y
RUN apt-get install libsqlite3-dev -y
RUN apt-get install tk-dev -y
RUN apt-get install libssl-dev -y
RUN apt-get install openssl -y
RUN apt-get install libffi-dev -y

# OpenCV Stuff
RUN apt-get install libsm6 -y
RUN apt-get install libxrender1 -y
RUN apt-get install libfontconfig1 -y
RUN apt-get install libopencv-dev -y
RUN apt-get install python3-opencv -y
RUN apt-get install python3-h5py -y
RUN apt-get install yasm -y
RUN apt-get install ffmpeg -y
RUN apt-get install libswscale-dev -y
RUN apt-get install libtbb2 -y
RUN apt-get install libtbb-dev -y
RUN apt-get install libjpeg-dev -y
RUN apt-get install libpng-dev -y
RUN apt-get install libtiff-dev -y
RUN apt-get install libavformat-dev -y
RUN apt-get install libpq-dev -y
RUN apt-get install libxvidcore-dev -y
RUN apt-get install libx264-dev -y
RUN apt-get install libavcodec-dev -y
RUN apt-get install libv4l-dev -y
RUN apt-get install libgtk-3-dev -y
RUN apt-get install libdc1394-22-dev -y
RUN apt-get install libjpeg62 -y
RUN apt-get install libopenjp2-7 -y
RUN apt-get install libilmbase-dev -y
# RUN apt-get install libilmbase24 -y
RUN apt-get install libatlas-base-dev -y
RUN apt-get install libgstreamer1.0-dev -y
RUN apt-get install openexr -y
RUN apt-get install libopenexr-dev -y

ENV TZ="US/Eastern"
RUN echo "US/Eastern" > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata
ARG USERNAME="morphs"
ARG PASSWORD="asdfasdf"
RUN useradd -m $USERNAME -p $PASSWORD -s "/bin/bash"
RUN mkdir -p /home/$USERNAME
RUN chown -R $USERNAME:$USERNAME /home/$USERNAME
RUN usermod -aG sudo $USERNAME
RUN echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

USER $USERNAME
WORKDIR /home/$USERNAME
RUN mkdir -p /home/$USERNAME/SHARING
RUN sudo chown -R $USERNAME:$USERNAME /home/$USERNAME/SHARING

# RUN mkdir -p /home/$USERNAME/SHARING/RaspiCameraMotionTracker/
# COPY ./v2 /home/$USERNAME/SHARING/RaspiCameraMotionTracker/
# COPY ./go.mod /home/$USERNAME/SHARING/RaspiCameraMotionTracker/
# COPY ./go.sum /home/$USERNAME/SHARING/RaspiCameraMotionTracker/
# COPY ./config.json /home/$USERNAME/SHARING/RaspiCameraMotionTracker/
# RUN sudo chown -R $USERNAME:$USERNAME /home/$USERNAME/SHARING/RaspiCameraMotionTracker

# Build OpenCV for GoVersion
RUN mkdir -p /home/$USERNAME/SHARING/opencv/
WORKDIR /home/$USERNAME/SHARING/opencv/
ENV OPENCV_VERSION="4.5.0"
RUN wget https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip
RUN unzip ${OPENCV_VERSION}.zip
WORKDIR /home/$USERNAME/SHARING/opencv/opencv-$OPENCV_VERSION
RUN mkdir -p /home/$USERNAME/SHARING/opencv/opencv-$OPENCV_VERSION/cmake_binary
WORKDIR /home/$USERNAME/SHARING/opencv/opencv-$OPENCV_VERSION/cmake_binary
RUN cmake \
-D OPENCV_GENERATE_PKGCONFIG=ON \
-D WITH_CUDA=OFF \
-D BUILD_opencv_python=OFF \
-D BUILD_opencv_python2=OFF \
-D BUILD_opencv_python3=OFF \
-D CMAKE_BUILD_TYPE=RELEASE \
-D BUILD_SHARED_LIBS=ON \
-D CMAKE_INSTALL_PREFIX=/usr \
-D INSTALL_C_EXAMPLES=OFF \
-D INSTALL_PYTHON_EXAMPLES=OFF \
-D BUILD_PYTHON_SUPPORT=OFF \
-D BUILD_NEW_PYTHON_SUPPORT=OFF \
-D WITH_TBB=ON \
-D WITH_PTHREADS_PF=ON \
-D WITH_OPENNI=OFF \
-D WITH_OPENNI2=ON \
-D WITH_EIGEN=ON \
-D BUILD_DOCS=OFF \
-D BUILD_TESTS=OFF \
-D BUILD_PERF_TESTS=OFF \
-D BUILD_EXAMPLES=OFF \
-D WITH_OPENCL=$OPENCL_ENABLED \
-D USE_GStreamer=ON \
-D WITH_GDAL=ON \
-D WITH_CSTRIPES=ON \
-D ENABLE_FAST_MATH=1 \
-D WITH_OPENGL=ON \
-D WITH_QT=OFF \
-D WITH_IPP=OFF \
-D WITH_FFMPEG=ON \
-D WITH_PROTOBUF=ON \
-D BUILD_PROTOBUF=ON \
-D CMAKE_SHARED_LINKER_FLAGS=-Wl,-Bsymbolic \
-D WITH_V4L=ON ..
# -D WITH_NGRAPH=ON \
# RUN make
#RUN n=$(nproc) && ((c=$n-2)) && make -j $c
RUN make
RUN sudo make install
RUN sudo ldconfig
RUN sudo chown $USERNAME:video /dev/video0 || echo ""

WORKDIR /home/$USERNAME

# # RUN rm /${OPENCV_VERSION}.zip
# # RUN rm -r /opencv-${OPENCV_VERSION}

COPY ./go_install.sh /home/$USERNAME/go_install.sh
RUN sudo chmod +x /home/$USERNAME/go_install.sh
RUN sudo chown $USERNAME:$USERNAME /home/$USERNAME/go_install.sh
RUN /home/$USERNAME/go_install.sh
# base64Encode 'export GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$(stat --format="%s" /home/morphs/go.tar.gz) / 1000" | bc -l)) && echo Extracting [$TAR_CHECKPOINT] of $GO_TAR_KILOBYTES kilobytes /usr/local/go'
RUN sudo tar --checkpoint=100 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo ZXhwb3J0IEdPX1RBUl9LSUxPQllURVM9JChwcmludGYgIiUuM2ZcbiIgJChlY2hvICIkKHN0YXQgLS1mb3JtYXQ9IiVzIiAvaG9tZS9tb3JwaHMvZ28udGFyLmd6KSAvIDEwMDAiIHwgYmMgLWwpKSAmJiBlY2hvIEV4dHJhY3RpbmcgWyRUQVJfQ0hFQ0tQT0lOVF0gb2YgJEdPX1RBUl9LSUxPQllURVMga2lsb2J5dGVzIC91c3IvbG9jYWwvZ28= | base64 -d ; echo); eval $cmd"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
RUN sudo ln -s /usr/local/go/bin/go /usr/local/bin/go

RUN mkdir -p /home/$USERNAME/MOTION_TRACKER
RUN sudo chown -R $USERNAME:$USERNAME /home/$USERNAME/MOTION_TRACKER
COPY . /home/$USERNAME/MOTION_TRACKER
RUN sudo chown -R $USERNAME:$USERNAME /home/$USERNAME/MOTION_TRACKER
WORKDIR /home/$USERNAME/MOTION_TRACKER
RUN go mod download

COPY ./core.go /home/morphs/go/pkg/mod/gocv.io/x/gocv@v0.25.0/core.go
RUN sudo chmod +r /home/morphs/go/pkg/mod/gocv.io/x/gocv@v0.25.0/core.go
RUN sudo rm ./core.go

# RUN go get -u gocv.io/x/gocv
# go get -u gocv.io/x/gocv@v0.27.0
# RUN chmod +w /home/morphs/go/pkg/mod/gocv.io/x/gocv@v0.25.0/core.go
# RUN python3 -c "exec(__import__('base64').b64decode('JycnCmZpbGVfcGF0aCA9ICIvaG9tZS9tb3JwaHMvZ28vcGtnL21vZC9nb2N2LmlvL3gvZ29jdkB2MC4yNS4wL2NvcmUuZ28iCmZpbGVfZGF0YSA9IE5vbmUKd2l0aCBvcGVuKCBmaWxlX3BhdGggLCAiciIgKSBhcyBmaWxlOgoJZmlsZV9kYXRhID0gZmlsZS5yZWFkKCkKCWZpbGVfZGF0YSA9IGZpbGVfZGF0YS5yZXBsYWNlKCAiMSA8PCAzMCIgLCAiMSA8PCAyMCIgKQp3aXRoIG9wZW4oIGZpbGVfcGF0aCAsICJ3IiApIGFzIGZpbGU6CglmaWxlLndyaXRlKCBmaWxlX2RhdGEgKQonJyc=').decode('utf-8'))"
# RUN chmod -w /home/morphs/go/pkg/mod/gocv.io/x/gocv@v0.25.0/core.go

RUN go build -o motion-tracker-server
RUN chmod +x motion-tracker-server

# COPY /dev/video0 /dev/video0
RUN sudo mknod /dev/video0 c 81 0 || echo ""
RUN sudo chmod 666 /dev/video0 || echo ""
RUN sudo chgrp video /dev/video0 || echo ""
RUN sudo chown $USERNAME:video /dev/video0 || echo ""

# ENV DISPLAY=:10.0
# ENTRYPOINT [ "/bin/bash" ]
# ENTRYPOINT [ "/home/morphs/MOTION_TRACKER/motion-tracker-server" ]
ENTRYPOINT [ "/home/morphs/MOTION_TRACKER/entrypoint.sh" ]