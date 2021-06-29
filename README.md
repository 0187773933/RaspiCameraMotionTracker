# Raspi Camera Motion Tracker

> - Exposes a mjpeg stream of raspi camera
> - POSTS Motion Frames to Python/Tensorflow Consumer Server

```bash
go run main.go config.json
```

```bash
export RPMT_SERVER_PORT=9767 && \
export RPMT_SERVER_MJPEG_ENDPOINT_URL="/frame.jpeg" && \
export RPMT_DEVICE_ID=0 && \
export RPMT_FRAME_INTERVAL_MILLISECONDS=50 && \
export RPMT_DELTA_THRESHOLD=5 && \
export RPMT_SHOW_DISPLAY="false" && \
export RPMT_MINIMUM_AREA=1000 && \
export RPMT_MINIMUM_MOTION_COUNTER_BEFORE_EVENT=25 && \
export RPMT_MINIMUM_EVENTS_BEFORE_ALERT=16 && \
export RPMT_ALERT_COOLOFF_DURATION_SECONDS=3 && \
export RPMT_ALERT_SERVER_POST_URL="https://example.com/" && \
go run main.go
```