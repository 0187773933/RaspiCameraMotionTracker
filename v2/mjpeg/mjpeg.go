package mjpeg

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	securecookie "github.com/gorilla/securecookie"
)

// Stream represents a single video feed.
type Stream struct {
	M             map[chan []byte]bool
	Frame         []byte
	Lock          sync.Mutex
	FrameInterval time.Duration
	Cookie *securecookie.SecureCookie
	JWTSecret []byte
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
"--" + boundaryWord + "\r\n" +
"Content-Type: image/jpeg\r\n" +
"Content-Length: %d\r\n" +
"X-Timestamp: 0.000000\r\n" +
"\r\n"

// func ( s *Stream ) IsAuthorized( endpoint func( http.ResponseWriter , *http.Request ) ) http.Handler {
// 	return http.HandlerFunc( func( w http.ResponseWriter , r *http.Request ) {
// 		cookie_value := make(map[string]string)
// 		if cookie , err := r.Cookie("rpmt-cookie"); err == nil {
// 			if err = COOKIE.Decode("rpmt-cookie", cookie.Value, &cookie_value); err == nil {
// 				http.Redirect( w , r , "/" , http.StatusUnauthorized )
// 				return
// 			}
// 		} else {
// 			http.Redirect( w , r , "/" , http.StatusUnauthorized )
// 			return
// 		}
// 		fmt.Println( cookie_value )
// 		token , err := jwt.Parse( cookie_value["token"] , func( token *jwt.Token ) ( interface{} , error ) {
// 			fmt.Println( "here again" )
// 			if _ , ok := token.Method.( *jwt.SigningMethodHMAC ); !ok {
// 				return nil, fmt.Errorf("There was an error")
// 			}
// 			return JWT_SECRET , nil
// 		})
// 		if err != nil { fmt.Fprintf(w, err.Error()) }
// 		if token.Valid {
// 			fmt.Println( token )
// 			endpoint( w , r )
// 		} else {
// 			// fmt.Fprintf( w , "Not Authorized" )
// 			http.Redirect( w , r , "/" , http.StatusUnauthorized )
// 		}
// 		return
// 	})
// }

func ( s *Stream ) ServeHTTP( w http.ResponseWriter , r *http.Request ) {
	cookie_value := make(map[string]string)
	if cookie , err := r.Cookie("rpmt-cookie"); err == nil {
		if err = s.Cookie.Decode("rpmt-cookie", cookie.Value, &cookie_value); err == nil {
			fmt.Println( cookie_value )
			token , err := jwt.Parse( cookie_value["token"] , func( token *jwt.Token ) ( interface{} , error ) {
				fmt.Println( "here again" )
				if _ , ok := token.Method.( *jwt.SigningMethodHMAC ); !ok {
					fmt.Printf("There was an error")
				}
				return s.JWTSecret , nil
			})
			if err != nil { fmt.Fprintf(w, err.Error()) }
			if token.Valid {
				fmt.Println( token )
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
				return
			}
		}
	}
	fmt.Println( "here ???" )
	http.Redirect( w , r , "/login" , http.StatusTemporaryRedirect )
	return
}


// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
// as long as we impliment ServeHTTP we are fine??? this isn't class shit.
// https://pkg.go.dev/net/http#Handle
// func ( s *Stream ) ServeHTTP( w http.ResponseWriter , r *http.Request ) {
// 	log.Println( "Stream:" , r.RemoteAddr , "connected" )
// 	w.Header().Add( "Content-Type" , "multipart/x-mixed-replace;boundary=" + boundaryWord )
// 	c := make( chan []byte )
// 	s.Lock.Lock()
// 	s.M[ c ] = true
// 	s.Lock.Unlock()
// 	for {
// 		time.Sleep( s.FrameInterval )
// 		b := <-c
// 		_ , err := w.Write( b )
// 		if err != nil { break }
// 	}
// 	s.Lock.Lock()
// 	delete( s.M , c )
// 	s.Lock.Unlock()
// 	log.Println( "Stream:" , r.RemoteAddr , "disconnected" )
// }

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
func NewStream( frame_interval_milliseconds int , cookie *securecookie.SecureCookie , jwt_secret []byte ) *Stream {
	return &Stream{
		M: make( map[ chan []byte ]bool ) ,
		Frame: make( []byte , len( headerf ) ),
		FrameInterval: ( time.Duration( frame_interval_milliseconds ) * time.Millisecond ) ,
		Cookie: cookie ,
		JWTSecret: jwt_secret ,
	}
}