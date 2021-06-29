package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
	_ "net/http/pprof"
	bcrypt "golang.org/x/crypto/bcrypt"
	auth "github.com/abbot/go-http-auth"
	//mjpeg "github.com/hybridgroup/mjpeg"
	mjpeg "github.com/0187773933/RaspiCameraMotionTracker/v2/mjpeg"
	motion "github.com/0187773933/RaspiCameraMotionTracker/v2"
)

// Switch to Fiber
// https://docs.gofiber.io/api/middleware/basicauth
// https://docs.gofiber.io/api/middleware#helmet
// https://docs.gofiber.io/api/app#handler

var stream *mjpeg.Stream

// func Secret( user , realm string ) string {
// 	if user == "john" {
// 		// password is "hello"
// 		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
// 	}
// 	return ""
// }

// https://github.com/abbot/go-http-auth/issues/56#issuecomment-393651189
// https://pkg.go.dev/golang.org/x/crypto/bcrypt#GenerateFromPassword
func Secret( user , realm string ) string {
	if user == "admin" {
		hashedPassword , err := bcrypt.GenerateFromPassword( []byte( "waduwaduwadu" ) , bcrypt.DefaultCost )
		if err == nil {
			return string( hashedPassword )
		}
	}
	return ""
}

// https://golang.org/pkg/net/http/#ServeMux.HandleFunc
// https://github.com/abbot/go-http-auth/blob/7f557639efd97bd84723b69471931553e1e0ade9/basic.go#L134
// https://github.com/hybridgroup/mjpeg/blob/master/stream.go#L18
// https://golang.org/pkg/net/http/#ResponseWriter
func handle( w http.ResponseWriter , r *auth.AuthenticatedRequest ) {
	w.Header().Add( "Content-Type" , "multipart/x-mixed-replace;boundary=MJPEGBOUNDARY" )
	c := make( chan []byte )
	stream.Lock.Lock()
	stream.M[ c ] = true
	stream.Lock.Unlock()
	for {
		time.Sleep( stream.FrameInterval )
		b := <-c
		_, err := w.Write( b )
		if err != nil {
			break
		}
	}
	stream.Lock.Lock()
	delete( stream.M , c )
	stream.Lock.Unlock()
	log.Println( "Stream:" , r.RemoteAddr , "disconnected" )
}

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

func main() {

	config := ParseConfig()
	stream = mjpeg.NewStream( config.FrameIntervalMilliseconds )

	motion_tracker := motion.NewTracker( stream , &config )
	go motion_tracker.Start()

	// start http server
	// https://medium.com/better-programming/hands-on-with-jwt-in-golang-8c986d1bb4c0
	// http://localhost:9363/frame.jpeg

	//http.Handle( "/frame.jpeg" , stream )

	//authenticator := auth.NewBasicAuthenticator( "localhost" , Secret )
	//http.HandleFunc( "/frame.jpeg" , authenticator.Wrap( handle ) )
	//http.HandleFunc( "/frame.jpeg" , handle  )
	// log.Fatal( http.ListenAndServe( "0.0.0.0:9767" , nil ) )

	http.Handle( config.ServerMJPEGEndpointURL , stream )
	http.ListenAndServe( fmt.Sprintf( "0.0.0.0:%s" , config.ServerPort ) , nil )
}