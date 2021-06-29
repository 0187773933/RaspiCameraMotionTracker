package mjpeg

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Stream represents a single video feed.
type Stream struct {
	M             map[chan []byte]bool
	Frame         []byte
	Lock          sync.Mutex
	FrameInterval time.Duration
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
"--" + boundaryWord + "\r\n" +
"Content-Type: image/jpeg\r\n" +
"Content-Length: %d\r\n" +
"X-Timestamp: 0.000000\r\n" +
"\r\n"

// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func ( s *Stream ) ServeHTTP( w http.ResponseWriter , r *http.Request ) {
	log.Println( "Stream:" , r.RemoteAddr , "connected" )
	w.Header().Add( "Content-Type" , "multipart/x-mixed-replace;boundary=" + boundaryWord )
	c := make( chan []byte )
	s.Lock.Lock()
	s.M[ c ] = true
	s.Lock.Unlock()
	for {
		time.Sleep( s.FrameInterval )
		b := <-c
		_ , err := w.Write( b )
		if err != nil { break }
	}
	s.Lock.Lock()
	delete( s.M , c )
	s.Lock.Unlock()
	log.Println( "Stream:" , r.RemoteAddr , "disconnected" )
}

// UpdateJPEG pushes a new JPEG frame onto the clients.
func ( s *Stream ) UpdateJPEG( jpeg []byte ) {
	header := fmt.Sprintf( headerf , len( jpeg ) )
	if len( s.Frame ) < len( jpeg ) + len( header ) {
		s.Frame = make( []byte , ( len( jpeg ) + len( header ) ) * 2 )
	}

	copy( s.Frame , header )
	copy( s.Frame[ len( header ) : ] , jpeg )

	s.Lock.Lock()
	for c := range s.M {
		// Select to skip streams which are sleeping to drop frames.
		// This might need more thought.
		select {
			case c <- s.Frame:
			default:
		}
	}
	s.Lock.Unlock()
}

// NewStream initializes and returns a new Stream.
func NewStream( frame_interval_milliseconds int ) *Stream {
	return &Stream{
		M: make( map[ chan []byte ]bool ) ,
		Frame: make( []byte , len( headerf ) ),
		FrameInterval: ( time.Duration( frame_interval_milliseconds ) * time.Millisecond ) ,
	}
}