#!/bin/bash
#sudo docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 -t raspi-motion-tracker-frame-consumer:latest --push .
sudo docker buildx build -m 8g --platform linux/arm/v6 -t xp6qhg9fmuolztbd2ixwdbtd1/raspi-motion-tracker:arm32test --push .