package main

import (
	"encoding/binary"
	"flag"
	"log"
	"os"
	"time"

	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
)

// This code is based on: https://github.com/bluenviron/gortsplib/blob/v2.0.0/examples/client-read-format-h264-save-to-disk/main.go

// Converts the values of the NTP Timestamp.
// The parameter secs are the seconds that have passed since the 1st of January of 1900 until that timestamp
// The parameter frac is used to get a more precise timestamp (up to nanoseconds). It consist of the fraction of a second that has passed since the last second, and it goes from 0 to 4294967295, the max value of a 32bits integer
func toUNIXTime(secs uint64, frac uint64) time.Time {
	seconds := int64(secs - 2208988800)               // We substract the seconds that passed between 1900 and 1970 so that we start counting at Unix time
	nanoseconds := int64(((frac * 1000000000) >> 32)) // We divide the second's fraction by the max value of a 32bits integer (2 elevated to 32) to get the nanoseconds in second unit and then pass it to nanoseconds (it's done the other way around since starting with the division makes the value too small to be accountable)
	return time.Unix(seconds, nanoseconds).UTC()
}

func retrieveTimeStamp(pkt *rtp.Packet) time.Time {
	// We get the timestamp in seconds and the fraction of the current second (more info on how seconds and second fraction work on the toUNIX function) and transform it to an int from a byte array
	// Those values are found on the first and second sections of the header extension respectively, representing in total the 8 first bytes of it
	headerExtension := pkt.Header.GetExtension(0)
	secBits := binary.BigEndian.Uint32(headerExtension[0:4])
	fracBits := binary.BigEndian.Uint32(headerExtension[4:8])

	//Since the timestamp is defined in NTP format, we most transform it back to UNIX to display it properly

	return toUNIXTime(uint64(secBits), uint64(fracBits))
}

func main() {

	// Parse command line arguments
	var rtspEndpoint = flag.String("rtsp-endpoint", "", "The rtsp endpoint from which the video comes from")
	flag.Parse()

	// Configure logger to output to stdout with microsecond precision
	log.SetOutput(os.Stdout)
	logger := log.Default()
	logger.SetFlags(log.Lmicroseconds)

	rtspClient := gortsplib.Client{}
	// Parse the URL
	rtspEndpointUrl, err := url.Parse(*rtspEndpoint)
	if err != nil {
		logger.Println("Couldn't parse the url", err)
		return
	}

	// Connect to the server
	err = rtspClient.Start(rtspEndpointUrl.Scheme, rtspEndpointUrl.Host)
	if err != nil {
		logger.Println("Couldn't connect to the RTSP server", err)
		return
	}
	defer rtspClient.Close()

	// Find published medias
	medias, baseURL, _, err := rtspClient.Describe(rtspEndpointUrl)
	if err != nil {
		logger.Println("Couldn't receive a response from the RTSP URL", err)
		return
	}

	// Find the H264 media and format
	var format *formats.H264
	h264Media := medias.FindFormat(&format)
	if h264Media == nil {
		logger.Println("Media not found")
		return
	}

	// Setup the chosen h264 media only
	_, err = rtspClient.Setup(h264Media, baseURL, 0, 0)
	if err != nil {
		logger.Println("Couldn't setup the media", err)
		return
	}

	// This method is called when a RTP packet arrives
	rtspClient.OnPacketRTP(h264Media, format, func(pkt *rtp.Packet) {
		timeStamp := retrieveTimeStamp(pkt)
		logger.Println("Timestamp: ", timeStamp)
	})

	// Start playing
	_, err = rtspClient.Play(nil)
	if err != nil {
		logger.Println("Failed to play the H264 media", err)
		return
	}

	// Wait until a fatal error
	logger.Println(rtspClient.Wait())
	return
}
