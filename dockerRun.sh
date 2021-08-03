#!/bin/bash
APP_NAME="rpmt-motion-tracker"
sudo docker rm $APP_NAME || echo ""
IMAGE_NAME="xp6qhg9fmuolztbd2ixwdbtd1/raspi-motion-tracker:arm32test"
id=$(sudo docker run -dit \
--name $APP_NAME \
-p 9767:9767 \
-v $(pwd)/config.json:/home/morphs/MOTION_TRACKER/config.json:ro \
--env="DISPLAY" \
-v /tmp/.X11-unix:/tmp/.X11-unix:rw \
--device /dev/video0 \
$IMAGE_NAME config.json)
sudo docker logs -f $id

# --entrypoint "sudo chown "morphs:video" /dev/video0 && /home/morphs/MOTION_TRACKER/motion-tracker-server /home/morphs/MOTION_TRACKER/config.json" \

# -v /dev/video0:/dev/video0 \
# --device "/dev/video0:/dev/video0" \
# -v ${PWD}/PythonVersion/built_wheels:/home/morphs/built_wheels \
# -v /tmp/.X11-unix:/tmp/.X11-unix \
# --device "/dev/video0:/dev/video0" \

# sudo docker exec -it 4ba44b079ac9 /bin/bash


# sudo docker run -it --rm --entrypoint "/bin/bash" xp6qhg9fmuolztbd2ixwdbtd1/raspi-motion-tracker:arm32test
