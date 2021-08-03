package motion

import (
	"fmt"
	"time"
	"bytes"
	"strings"
	"image"
	// "reflect"
	// "io/ioutil"
	"encoding/json"
	"encoding/base64"
	//"image/color"
	opencv "gocv.io/x/gocv"
	"net/http"

	try "github.com/manucorporat/try"
	mjpeg "github.com/0187773933/RaspiCameraMotionTracker/v2/mjpeg"
	// twilio https://github.com/sfreiberg/gotwilio
)

// https://github.com/hybridgroup/gocv/tree/release/cmd
// https://github.com/hybridgroup/gocv/blob/6240320eed51651fa2be9cfd304605b7497f4b6f/cmd/motion-detect/main.go#L3

type FrameClipping struct {
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
}
type TrackerConfig struct {
	ServerPort string `json:"server_port"`
	ServerMJPEGEndpointURL string `json:"server_mjpeg_endpoint_url"`
	DeviceID int `json:"device_id"`
	FrameIntervalMilliseconds int `json:"frame_interval_milliseconds"`
	DeltaThreshold int `json:"delta_threshold"`
	ShowDisplay bool `json:"show_display"`
	MinimumArea float64 `json:"minimum_area"`
	MinimumMotionCounterBeforeEvent int `json:"minimum_motion_counter_before_event"`
	MinimumEventsBeforeAlert int `json:"minimum_events_before_alert"`
	AlertCooloffDurationSeconds float64 `json:"alert_cooloff_duration_seconds"`
	AlertServerPostURL string `json:"alert_server_post_url"`
	FrameClipping FrameClipping `json:"frame_clipping"`
	TimeZone string `json:"time_zone"`
	JWTSecret string `json:"jwt_secret"`
	CookieSecret string `json:"cookie_32_secret"`
	CookieSalt string `json:"cookie_32_salt"`
}
type Tracker struct {
	Stream *mjpeg.Stream `json:"mjpeg_stream"`
	ConfigFilePath string `json:"config_file_path"`
	Config *TrackerConfig `json:"config"`
	LastAlertTime time.Time
}

func NewTracker( stream *mjpeg.Stream , config *TrackerConfig ) ( *Tracker ) {
	return &Tracker{
		Stream: stream ,
		Config: config ,
	}
}

func ( tracker *Tracker ) Alert( frame_buffer []uint8 ) {
	location , _ := time.LoadLocation( tracker.Config.TimeZone )
	tracker.LastAlertTime = time.Now().In( location )
	fmt.Println( "sending sms alert" )
	try.This( func() {
		frame_buffer_b64_string := base64.StdEncoding.EncodeToString( frame_buffer )
		post_data , _ := json.Marshal(map[string]string{
			"frame_buffer_b64_string": frame_buffer_b64_string ,
			"time_stamp": GetFormattedTimeString( tracker.LastAlertTime ) ,
		})
		client := &http.Client{}
		request , request_error := http.NewRequest( "POST" , tracker.Config.AlertServerPostURL , bytes.NewBuffer( post_data ) )
		request.Header.Set( "Content-Type" , "application/json" )
		if request_error != nil { fmt.Println( request_error ); }
		response , _ := client.Do( request )
		defer response.Body.Close()
		// body , body_error := ioutil.ReadAll( response.Body )
		// if body_error != nil { fmt.Println( body_error ); }
		// fmt.Println( body )
	}).Catch( func ( e try.E ) {
		// fmt.Println( e )
		fmt.Println( "Recieving Server Offline" )
	})
}

func GetFormattedTimeString( time_object time.Time ) ( result string ) {
	// https://stackoverflow.com/a/51915792
	// month_name := strings.ToUpper( time_object.Format( "Feb" ) ) // monkaHmm
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

// https://github.com/hybridgroup/gocv/issues/639
// https://github.com/hybridgroup/gocv/issues/535

// func (icv *ImageCV) Crop(left, top, right, bottom int) *ImageCV {
// 	croppedMat := icv.mat.Region(image.Rect(left, top, right, bottom))
//         resultMat := croppedMat.Clone()
// 	return &ImageCV{mat: resultMat}
// }

func ( tracker *Tracker ) Start() {
	fmt.Println( "Starting Motion Tracker" )

	tracker.LastAlertTime = time.Now()
	// this breaks if you change cool off constant
	duration , _ := time.ParseDuration( fmt.Sprintf( "-%ds" , tracker.Config.AlertCooloffDurationSeconds ) )
	tracker.LastAlertTime = tracker.LastAlertTime.Add( duration )

	// Vars
	motion_counter := 0
	events_counter := 0

	var webcam *opencv.VideoCapture
	var webcam_error error
	try.This( func() {
		fmt.Printf( "Opening Video Capture From : /dev/video%v\n" , tracker.Config.DeviceID )
		fmt.Println( tracker.Config.DeviceID )
		webcam , webcam_error = opencv.OpenVideoCapture( tracker.Config.DeviceID )
	}).Catch( func ( e try.E ) {
		fmt.Println( e )
		fmt.Printf( "probably don't have read access to /dev/video%v" , tracker.Config.DeviceID );
	})
	// webcam , webcam_error := opencv.OpenVideoCapture( tracker.Config.DeviceID )
	if webcam_error != nil { panic( webcam_error ) }
	defer webcam.Close()

	var window *opencv.Window
	if tracker.Config.ShowDisplay == true {
		window = opencv.NewWindow( "Motion Window" )
		defer window.Close()
	}

	frame := opencv.NewMat()
	defer frame.Close()
	delta := opencv.NewMat()
	defer delta.Close()

	threshold := opencv.NewMat()
	defer threshold.Close()

	mog2 := opencv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	//status := "Ready"
	//status_color_ready := color.RGBA{ 0 , 255 , 0 , 0 }
	//status_color_motion := color.RGBA{ 255 , 0 , 0 , 0 }

	// frame_clipping_min := image.Point{ tracker.Config.FrameClipping.X1 , tracker.Config.FrameClipping.Y1 }
	// frame_clipping_max := image.Point{ tracker.Config.FrameClipping.X2 , tracker.Config.FrameClipping.Y2 }
	// frame_clipping := image.Rectangle{ frame_clipping_min , frame_clipping_max }

	fmt.Printf( "Start reading device: %v\n" , tracker.Config.DeviceID )
	for {

		ok := webcam.Read( &frame );
		if !ok { fmt.Printf( "Device closed: %v\n" , tracker.Config.DeviceID ); panic( "device closed" ) }
		if frame.Empty() { continue }

		// 1.) Crop Frame
		// Default Raspi Camera Size
		// frame.Size() == [480 640]
		// x_test , _ := frame.ToImage()
		// x_test_bounds := x_test.Bounds()
		// fmt.Println( reflect.TypeOf( frame ) )
		// cropped := frame.Region( frame_clipping )
		// fmt.Println( reflect.TypeOf( cropped ) )

		// cropped := frame.Region( frame_clipping )
		// frame = cropped
		// defer cropped.Close()

		// Somehow if you mess with the min aka (x1,y1) point, it panics
		// it does successfully clip it though, but something runs away with memory.
		// no clipping ????
		// frame = frame.Region( frame_clipping )
		// delta = delta.Region( frame_clipping )
		// threshold = threshold.Region( frame_clipping )
		// fmt.Println( frame.Rows() )
		// ?? frame.Reshape( 400 , 200 )

		//status = "Ready"

		// Create Delta ( foreground only )
		mog2.Apply( frame , &delta )

		// Create Threshold from Delta
		opencv.Threshold( delta , &threshold , 25 , 255 , opencv.ThresholdBinary )

		// Dilate Threshold , still no intervals , can we just do this twice ?
		kernel := opencv.GetStructuringElement( opencv.MorphRect, image.Pt( 3 , 3 ) )
		defer kernel.Close()
		opencv.Dilate( threshold , &threshold , kernel )
		opencv.Dilate( threshold , &threshold , kernel )

		// Find Contours
		contours := opencv.FindContours( threshold , opencv.RetrievalExternal , opencv.ChainApproxSimple )
		now := time.Now()
		inside_cooloff := true
		difference := now.Sub( tracker.LastAlertTime )
		if difference.Seconds() > tracker.Config.AlertCooloffDurationSeconds {
			inside_cooloff = false
		}
		fmt.Println( difference , inside_cooloff )
		for _ , contour := range contours {
			area := opencv.ContourArea( contour )
			if area < tracker.Config.MinimumArea {
				continue
			} else {
				if inside_cooloff == false {
					motion_counter += 1
				}
			}
			//status = "Motion detected"
			// if ShowDisplay == true {
			// 	opencv.DrawContours( &frame , contours , i , status_color_motion , 2 )
			// 	rect := opencv.BoundingRect( contour )
			// 	opencv.Rectangle( &frame , rect , color.RGBA{ 0 , 0 , 255 , 0 } , 2 )
			// }

			//opencv.DrawContours( &frame , contours , i , status_color_motion , 2 )
			//rect := opencv.BoundingRect( contour )
			//opencv.Rectangle( &frame , rect , color.RGBA{ 0 , 0 , 255 , 0 } , 2 )
		}

		fmt.Printf( "Motion Counter === %d === Events Counter === %d\n" , motion_counter , events_counter )

		// Show Display
		//opencv.PutText( &frame , status , image.Pt( 10 , 20 ) , opencv.FontHersheyPlain , 1.2 , status_color_ready , 2 )
		if tracker.Config.ShowDisplay == true {
			window.IMShow( frame )
			if window.WaitKey( 1 ) == 27 {
				break
			}
		}

		// Update MJPEG Stream
		frame_buffer , _ := opencv.IMEncode ( ".jpg" , frame )
		tracker.Stream.UpdateJPEG( frame_buffer )
		//if frame_buffer_error == nil {
		//	stream.UpdateJPEG( frame_buffer )
		//} else {
		//	fmt.Println( frame_buffer_error )
		//}

		// Calculate State Decisions Based on Current Value of motion_counter
		if motion_counter >= tracker.Config.MinimumMotionCounterBeforeEvent { events_counter += 1; motion_counter = 0 }
		if events_counter >= tracker.Config.MinimumEventsBeforeAlert { go tracker.Alert( frame_buffer ); events_counter = 0; motion_counter = 0 }

	}

}