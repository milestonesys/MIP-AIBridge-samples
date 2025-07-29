package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	protocols_ds "ds-playback/generated/protos/ds/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Parse command line arguments
	var startTime = flag.String("start-time", "yyyy-mm-dd 24:00", "The start time of the recording you want to fetch")
	var endTime = flag.String("end-time", "yyyy-mm-dd 24:00", "The end time of the recording you want to fetch")
	var externalHostname = flag.String("external-hostname", "localhost", "Hostname or IP address of the machine running Milestone AI Bridge in debug mode")
	var deviceID = flag.String("device-id", "", "The camera ID used to extract the video recording")
	var streamID = flag.String("stream-id", "", "The stream ID of the specified device")
	flag.Parse()

	// Configure logger to output to stdout with microsecond precision
	log.SetOutput(os.Stdout)
	logger := log.Default()
	logger.SetFlags(log.Lmicroseconds)

	var err error
	layout := "2006-01-02 15:04" // Expected input format

	var startDateUTC time.Time
	startDateUTC, err = time.ParseInLocation(layout, *startTime, time.Local)
	if err != nil {
		logger.Println("Failed to parse start date: ", err)
		return
	}

	var endDateUTC time.Time
	endDateUTC, err = time.ParseInLocation(layout, *endTime, time.Local)
	if err != nil {
		logger.Println("Failed to parse end date: ", err)
		return
	}

	var streamingUrl = *externalHostname + ":9898"
	var cameraStreamId = *deviceID + "/" + *streamID

	// Set up a channel to handle interrupt signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Ensure the data directory exists for storing recorded videos
	dataFolderPath := "data"
	if _, err := os.Stat(dataFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dataFolderPath, os.ModePerm)
		if err != nil {
			logger.Println("Failed to create data directory: ", err)
			return
		}
	}

	// Open a file for writing video frames
	f, err := os.OpenFile(fmt.Sprintf("%s/playback-video-%s.bin", dataFolderPath, uuid.New().String()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Println("Failed to open file for writing: ", err)
		return
	}
	defer f.Close()

	logger.Println("Initiating playback request...")
	// Connect to the VMS bridge
	conn, err := grpc.NewClient(streamingUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(64*1024*1024)))
	if err != nil {
		logger.Println("Failed to connect to streaming server: ", err)
		return
	}
	client := protocols_ds.NewDirectStreamingClient(conn)
	defer conn.Close()

	// Create a PlaybackStream client
	playbackStreamClient, err := client.PlaybackStream(context.Background(), &protocols_ds.PlaybackRequest{
		StreamId: cameraStreamId,
		From:     timestamppb.New(startDateUTC),
		To:       timestamppb.New(endDateUTC),
	})
	if err != nil {
		logger.Println("Failed to create playback stream client: ", err)
		return
	}
	defer playbackStreamClient.CloseSend()

	logger.Println("Playback sequence started.")
	defer logger.Println("Playback finished. Exiting.")

	// Continuously receive and process playback data
	for {
		// Exit if an interrupt signal is received
		select {
		case <-signals:
			logger.Println("Interrupt signal received. Exiting playback loop.")
			return
		default:
		}

		streamData, err := playbackStreamClient.Recv()
		if err != nil {
			if err == io.EOF {
				logger.Println("Playback sequence completed.")
			} else {
				logger.Println("Error receiving playback data: ", err)
				return
			}
			// Exit the loop to close the connection
			break
		}

		// Write video frames to the output file
		for _, frame := range streamData.GetVideoData().Frames {
			logger.Println("Frame timestamp:", frame.GetTimestamp().AsTime().Format(time.RFC3339Nano))
			f.Write(frame.Data)
		}
	}
}
