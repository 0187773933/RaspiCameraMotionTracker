package motion

import (
	"fmt"
	"image"
	opencv "gocv.io/x/gocv"
)

// https://vitux.com/opencv_ubuntu/
// nano ~/.bashrc
// export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
// source ~/.bashrc
// sudo apt-get install libopencv-dev

// sudo ln -s /usr/lib/x86_64-linux-gnu/pkgconfig/opencv4.pc /usr/lib/x86_64-linux-gnu/pkgconfig/opencv.pc

// https://github.com/hybridgroup/gocv/blob/release/cmd/mjpeg-streamer/main.go

// https://github.com/hybridgroup/gocv/blob/6240320eed51651fa2be9cfd304605b7497f4b6f/cmd/motion-detect/main.go#L3

func Test() {
	fmt.Println( "test" )
	webcam, _ := opencv.OpenVideoCapture( 0 )
	window_greyscale := opencv.NewWindow( "Greyscale" )
	window_threshold := opencv.NewWindow( "Threshold" )
	window_delta := opencv.NewWindow( "Delta" )
	img := opencv.NewMat()
	greyscale := opencv.NewMat()
	greyscale_blurred := opencv.NewMat()
	ksize := image.Point{ X: 21 , Y: 21 }

	delta := opencv.NewMat()
	average := opencv.NewMat()
	average_abs := opencv.NewMat()

	threshold := opencv.NewMat()
	threshold_size := image.Point{ X: 1 , Y: 1 }

	motion_counter := 0
	first_run := true
	for {
		webcam.Read( &img )

		// GreyScale
		//fmt.Println( img.ToImage() )
		// https://godoc.org/gocv.io/x/gocv#CvtColor
		// https://godoc.org/gocv.io/x/gocv#ColorConversionCode
		// https://github.com/ceberous/RaspiMotionAlarmRewrite/blob/master/py_scripts/motion_simple_rewrite_fixed.py#L358
		opencv.CvtColor( img , &greyscale , opencv.ColorBGRToGray )
		// https://godoc.org/gocv.io/x/gocv#GaussianBlur
		opencv.GaussianBlur( greyscale , &greyscale_blurred , ksize , 0.0 , 0.0 , 0 )


		// Delta
		// https://shimat.github.io/opencvsharp_docs/html/c733d983-da1f-4c39-d183-e80c7862450f.htm
		// https://godoc.org/gocv.io/x/gocv#AddWeighted
		if first_run == true { average = greyscale.Clone() }
		opencv.AddWeighted( greyscale , 0.5 , average , 0.5 , 0.5 , &delta )
		opencv.ConvertScaleAbs( average , &average_abs , 1 , 0 )
		opencv.AbsDiff( greyscale , average_abs , &delta )

		// Threshold
		opencv.Threshold( delta , &threshold , 5.0 , 255.0 , opencv.ThresholdBinary )
		// https://github.com/hybridgroup/gocv/issues/682
		// https://godoc.org/gocv.io/x/gocv#GetStructuringElement
		kernel_material := opencv.GetStructuringElement( opencv.MorphRect , threshold_size )
		opencv.Dilate( threshold , &threshold , kernel_material )
		// frameThreshold = cv2.dilate( frameThreshold , None , iterations=2 )

		// Find Movement
		countours := opencv.FindContours( threshold.Clone() , opencv.RetrievalExternal , opencv.ChainApproxSimple )
		for _ , point := range countours {
			contour_area := opencv.ContourArea( point )
			if contour_area < 500 {
				motion_counter = 0
			} else {
				motion_counter += 1
			}
		}
		fmt.Println( motion_counter )
		// # Search for Movment
		// ( cnts , _ ) = cv2.findContours( frameThreshold.copy() , cv2.RETR_EXTERNAL , cv2.CHAIN_APPROX_SIMPLE )
		// for c in cnts:
		// 	if cv2.contourArea( c ) < min_area:
		// 		motionCounter = 0 # ???
		// 		continue
		// 	motionCounter += 1

		if first_run == true { first_run = false }
		window_greyscale.IMShow( greyscale )
		window_greyscale.WaitKey( 1 )
		window_threshold.IMShow( threshold )
		window_threshold.WaitKey( 1 )
		window_delta.IMShow( delta )
		window_delta.WaitKey( 1 )
		//window.IMShow( threshold )
		//window.WaitKey( 1 )
	}
}