package main

import (
	"fmt"
	"log"
	"time"
	// "io/ioutil"
	// "reflect"
	"net/http"
	"encoding/base64"
	// json "encoding/json"
	_ "net/http/pprof"
	bcrypt "golang.org/x/crypto/bcrypt"
	auth "github.com/abbot/go-http-auth"
	//mjpeg "github.com/hybridgroup/mjpeg"
	mjpeg "github.com/0187773933/RaspiCameraMotionTracker/v2/mjpeg"
	motion "github.com/0187773933/RaspiCameraMotionTracker/v2/motion"
	utils "github.com/0187773933/RaspiCameraMotionTracker/v2/utils"
	jwt "github.com/dgrijalva/jwt-go"
	securecookie "github.com/gorilla/securecookie"
)

// Switch to Fiber
// https://docs.gofiber.io/api/middleware/basicauth
// https://docs.gofiber.io/api/middleware#helmet
// https://docs.gofiber.io/api/app#handler

var stream *mjpeg.Stream
var JWT_SECRET = []byte("asdf")
var COOKIE_SECRET = []byte("asdf")
var COOKIE_SALT = []byte("asdf")
var COOKIE *securecookie.SecureCookie

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

func isAuthorized( endpoint func( http.ResponseWriter , *http.Request ) ) http.Handler {
	return http.HandlerFunc( func( w http.ResponseWriter , r *http.Request ) {
		cookie_value := make(map[string]string)
		if cookie , err := r.Cookie("rpmt-cookie"); err == nil {
			if err = COOKIE.Decode("rpmt-cookie", cookie.Value, &cookie_value); err == nil {
				fmt.Println( cookie_value )
				token , err := jwt.Parse( cookie_value["token"] , func( token *jwt.Token ) ( interface{} , error ) {
					fmt.Println( "here again" )
					if _ , ok := token.Method.( *jwt.SigningMethodHMAC ); !ok {
						fmt.Printf("There was an error")
					}
					return JWT_SECRET , nil
				})
				if err != nil { fmt.Fprintf(w, err.Error()) }
				if token.Valid {
					fmt.Println( token )
					endpoint( w , r )
					return
				}
			}
		}
		http.Redirect( w , r , "/" , http.StatusUnauthorized )
		return
	})
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "asdf"
	claims["exp"] = time.Now().Add(time.Minute * 1036800).Unix()
	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

type LoginResult struct {
	Token string `json:"token"`
}
// func login( w http.ResponseWriter , r *auth.AuthenticatedRequest ) {
func login( w http.ResponseWriter , r *http.Request ) {
	validToken , err := GenerateJWT()
	if err != nil { fmt.Println("Failed to generate token") }
	// w.Header().Set( "Content-Type" , "application/json" )
	// data := LoginResult{}
	// data.Token = validToken
	// w.WriteHeader( http.StatusCreated )
	// logincookie , _ := securecookie.New( "olahmb" , COOKIE_SECRET , securecookie.Params{
	// 	Path:     "/frame.jpeg" ,           // cookie received only when URL starts with this path
	// 	Domain:   "localhost" ,    // cookie received only when URL domain matches this one
	// 	MaxAge:   3600 ,             // cookie becomes invalid 3600 seconds after it is set
	// 	HTTPOnly: true ,             // disallow access by remote javascript code
	// 	Secure:   false ,             // cookie received only with HTTPS, never with HTTP
	// 	SameSite: securecookie.Lax , // cookie received with same or sub-domain names
	// })
	// fmt.Println( logincookie )
	// logincookie.SetValue( w , []byte( validToken ) )
	value := map[string]string{
		"token": validToken ,
	}
	encoded , _ := COOKIE.Encode( "rpmt-cookie" , value );
	cookie := &http.Cookie{
		Name:  "rpmt-cookie",
		Value: encoded ,
		Path:  "/frame.jpeg" ,
		// Path:  "/test" ,
		Secure: false ,
		HttpOnly: false ,
	}
	fmt.Println( cookie )
	http.SetCookie( w , cookie )

	http.Redirect( w , r , "/frame.jpeg" , http.StatusTemporaryRedirect )
	// http.Redirect( w , r , "/test" , http.StatusTemporaryRedirect )
	// http.Redirect( w ,r , "/test" , http.StatusOK )
	// http.Redirect( w ,r , "/login" , http.StatusUnauthorized )

	// json.NewEncoder( w ).Encode( data )
	// jsonResp , _ := json.Marshal( data )
	// w.Write( jsonResp )
	return
}

type TestAuthResult struct {
	Result string `json:"result"`
}
func test_authenticated( w http.ResponseWriter , r *http.Request ) {
	// w.Header().Set( "Content-Type" , "application/json" )
	// data := TestAuthResult{}
	// data.Result = "success"
	// w.WriteHeader( http.StatusCreated )
	// json.NewEncoder( w ).Encode( data )
	// jsonResp , _ := json.Marshal( data )
	// w.Write( jsonResp )
	fmt.Println( "here???" )
}

func main() {

	config := utils.ParseConfig()
	JWT_SECRET = []byte( config.JWTSecret )
	// COOKIE_SECRET = []byte( config.CookieSecret )
	COOKIE_SECRET , _ = base64.StdEncoding.DecodeString( config.CookieSecret )
	COOKIE_SALT , _ = base64.StdEncoding.DecodeString( config.CookieSalt )
	COOKIE = securecookie.New( COOKIE_SECRET , COOKIE_SALT )
	stream = mjpeg.NewStream( config.FrameIntervalMilliseconds , COOKIE , JWT_SECRET )

	motion_tracker := motion.NewTracker( stream , &config )
	// fmt.Println( motion_tracker )
	go motion_tracker.Start()

	// start http server
	// https://medium.com/better-programming/hands-on-with-jwt-in-golang-8c986d1bb4c0
	// http://localhost:9363/frame.jpeg

	//http.Handle( "/frame.jpeg" , stream )

	//authenticator := auth.NewBasicAuthenticator( "localhost" , Secret )
	//http.HandleFunc( "/frame.jpeg" , authenticator.Wrap( handle ) )
	//http.HandleFunc( "/frame.jpeg" , handle  )
	// log.Fatal( http.ListenAndServe( "0.0.0.0:9767" , nil ) )

	// somehow it never calls ServeHTTP() ??? http.Handle implicitly calls it ???
	// https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/

	http.Handle( config.ServerMJPEGEndpointURL , stream )
	http.HandleFunc( "/login" , login )
	http.Handle( "/test" , isAuthorized( test_authenticated ) )
	http.ListenAndServe( fmt.Sprintf( "0.0.0.0:%s" , config.ServerPort ) , nil )
}