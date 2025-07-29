import argparse
import asyncio
import datetime
import logging
import os
import uuid

import grpc
import google.protobuf.timestamp_pb2 as timestamp_pb2

import generated.protos.ds.v1.ds_pb2_grpc as ds_pb2_grpc
import generated.protos.ds.v1.ds_pb2 as ds_pb2


DS_PORT = 9898
DATETIME_FORMAT = '%Y-%m-%d %H:%M'
DEFAULT_START_DATETIME = (datetime.datetime.now() - datetime.timedelta(minutes=30)).strftime(DATETIME_FORMAT)
DEFAULT_END_DATETIME = datetime.datetime.now().strftime(DATETIME_FORMAT)

def parse_datetime(datetime_str: str) -> datetime.datetime:
    try:
        return datetime.datetime.strptime(datetime_str, DATETIME_FORMAT).replace(tzinfo=datetime.timezone.utc)
    except ValueError as e:
        logging.critical("Unable to parse date: %s", e)
        raise


async def run(start_date_utc:datetime.datetime, end_date_utc:datetime.datetime, streaming_url:str, camera_stream_id:str, data_folder:str) -> None:
    async with grpc.aio.insecure_channel(streaming_url) as channel:
        stub = ds_pb2_grpc.DirectStreamingStub(channel)
        file_path = os.path.join(data_folder, f"playback-video-{uuid.uuid4()}.bin")
        with open(file_path, "ab") as f:
            logging.info("Starting playback")

            # Cast the end time to timestamp_pb2.Timestamp
            end_time = timestamp_pb2.Timestamp()
            end_time.FromDatetime(end_date_utc)

            # Define the playback request
            request = ds_pb2.PlaybackRequest(
                stream_id = camera_stream_id,
                to=end_time
            )

            # WORKAROUND: Define the from attr for the playback request
            setattr(request, 'from', start_date_utc)

            # Start the request
            stream_call:grpc.aio._call.UnaryStreamCall = stub.PlaybackStream(request)

            logging.info("start of sequence")
            while True:
                # Fetch stream
                it = await stream_call.read()
                if not it:
                    break
                # Parse the video data to generate the video file
                if it.video_data:
                    for frame in it.video_data.frames:
                        logging.info(frame.timestamp.ToDatetime())
                        f.write(frame.data)
            logging.info("end of sequence")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--start-time', default=DEFAULT_START_DATETIME,
                        help='The start time of the recording you want to fetch', type=str)
    parser.add_argument('--end-time', default=DEFAULT_END_DATETIME,
                        help='The end time of the recording you want to fetch', type=str)
    parser.add_argument('--external-hostname', default="localhost",
                        help='The hostname or ip address of the machine running Milestone AI Bridge in debug mode', type=str)
    parser.add_argument('--device-id', default="",
                        help='The camera id to use for extracting the video recording', type=str)
    parser.add_argument('--stream-id', default="",
                        help='The stream id of the specified device', type=str)
    args = parser.parse_args()

    logging.basicConfig(level=logging.INFO)
    logging.info("Arguments: %s", vars(args))

    start_date_utc = parse_datetime(args.start_time)
    end_date_utc = parse_datetime(args.end_time)
    streaming_url = f"{args.external_hostname}:{DS_PORT}"
    camera_stream_id = f"{args.device_id}/{args.stream_id}"

    logging.info("Start Time UTC: %s", start_date_utc)
    logging.info("End Time UTC: %s", end_date_utc)
    logging.info("Streaming URL: %s", streaming_url)
    logging.info("Camera/Stream ID: %s", camera_stream_id)

    data_folder_name = "data"
    if not os.path.isabs(data_folder_name):
        data_folder_name = os.path.join(os.path.dirname(os.path.abspath(__file__)), data_folder_name)
    
    try:
        if not os.path.exists(data_folder_name):
            os.makedirs(data_folder_name)
    except Exception as e:
        logging.critical("Failed to create data folder %s: %s", data_folder_name, e)

    asyncio.run(run(start_date_utc, end_date_utc, streaming_url, camera_stream_id, data_folder_name))
