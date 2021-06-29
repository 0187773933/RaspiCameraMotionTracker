package utils

import (
	"os"
	"strconv"
	"io/ioutil"
	"encoding/json"
	motion "github.com/0187773933/RaspiCameraMotionTracker/v2/motion"
)

func ParseConfig() ( config motion.TrackerConfig ) {
	if len( os.Args ) < 2 {
		server_port := os.Getenv( "RPMT_SERVER_PORT" )
		server_mjpeg_endpoint_url := os.Getenv( "RPMT_SERVER_MJPEG_ENDPOINT_URL" )
		device_id , _ := strconv.Atoi( os.Getenv( "RPMT_DEVICE_ID" ) )
		frame_interval_milliseconds , _ := strconv.Atoi( os.Getenv( "RPMT_FRAME_INTERVAL_MILLISECONDS" ) )
		delta_threshold , _ := strconv.Atoi( os.Getenv( "RPMT_DELTA_THRESHOLD" ) )
		show_display , _ := strconv.ParseBool( os.Getenv( "RPMT_SHOW_DISPLAY" ) )
		minimum_area , _ := strconv.ParseFloat( os.Getenv( "RPMT_MINIMUM_AREA" ) , 64 )
		minimum_motion_counter_before_event , _ := strconv.Atoi( os.Getenv( "RPMT_MINIMUM_MOTION_COUNTER_BEFORE_EVENT" ) )
		minimum_events_before_alert , _ := strconv.Atoi( os.Getenv( "RPMT_MINIMUM_EVENTS_BEFORE_ALERT" ) )
		alert_cooloff_duration_seconds , _ := strconv.ParseFloat( os.Getenv( "RPMT_ALERT_COOLOFF_DURATION_SECONDS" ) , 64 )
		alert_server_post_url := os.Getenv( "RPMT_ALERT_SERVER_POST_URL" )

		config.ServerPort  = server_port
		config.ServerMJPEGEndpointURL = server_mjpeg_endpoint_url
		config.DeviceID = device_id
		config.FrameIntervalMilliseconds = frame_interval_milliseconds
		config.DeltaThreshold = delta_threshold
		config.ShowDisplay = show_display
		config.MinimumArea = minimum_area
		config.MinimumMotionCounterBeforeEvent = minimum_motion_counter_before_event
		config.MinimumEventsBeforeAlert = minimum_events_before_alert
		config.AlertCooloffDurationSeconds = alert_cooloff_duration_seconds
		config.AlertServerPostURL = alert_server_post_url
	} else {
		config_json_data , _ := ioutil.ReadFile( os.Args[ 1 ] )
		json.Unmarshal( config_json_data , &config )
	}
	return
}