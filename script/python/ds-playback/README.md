# DS-Playback

This sample script demonstrates how to use the Direct Streaming protocol to fetch playback recordings from XProtect. The Direct Streaming protocol is based on gRPC. The proto file is currently published in the `AI Bridge developer reference` at [http://<ai-bridge-external-hostname>:4000/](http://<ai-bridge-external-hostname>:4000/).

The proto file can be downloaded from `http://<ai-bridge-external-hostname>:4000/`. Make sure that Milestone AI Bridge is running in debug mode.

## Getting Started

The directory contains the following files:

```txt
.
├── Makefile
├── README.md
├── main.py
└── requirements.txt

0 directories, 4 files
```

1. **Install required packages for Makefile.**

    On Ubuntu, run:

    ```bash
    sudo apt-get -y install python3-pip python3-venv wget make curl jq
    ```

2. **Install dependencies and download required files.**

    In this directory, run:

    ```bash
    make setup SYSTEM_IP=<ai-bridge-external-hostname>
    ```

    When the setup is finished:

    - A python environment is created with all the required dependencies to compile the protocol buffers and run the grpc server.
    - The proto file will be downloaded to [./bin/protos/ds/v1/ds.proto](./bin/protos/ds/v1/ds.proto).
    - Generated Python code for the playback client will be located at [./generated/protos/ds/v1/ds_pb2.py](./generated/protos/ds/v1/ds_pb2.py).

3. **Activate the python environment**

    Run:

    ```bash
    source playback_venv/bin/activate
    ```

4. **View help for usage instructions.**

    Run:

    ```bash
    PYTHONPATH=./generated/protos/ds/v1/ python3 main.py --help
    ```

    Output:

    ```txt
    usage: main.py [-h] [--start-time START_TIME] [--end-time END_TIME] [--external-hostname EXTERNAL_HOSTNAME] [--device-id DEVICE_ID] [--stream-id STREAM_ID]

    options:
      -h, --help            show this help message and exit
      --start-time START_TIME
                            The start time of the recording you want to fetch
      --end-time END_TIME   The end time of the recording you want to fetch
      --external-hostname EXTERNAL_HOSTNAME
                            The hostname or ip address of the machine running Milestone AI Bridge in debug mode
      --device-id DEVICE_ID
                            The camera id to use for extracting the video recording
      --stream-id STREAM_ID
                            The stream id of the specified device
    ```

5. **Extract a video recording from XProtect.**

    Run:

    ```bash
    PYTHONPATH=./generated/protos/ds/v1/ python3 main.py --start-time "<yyyy-mm-dd HH:MM>" --end-time "<yyyy-mm-dd HH:MM>" --external-hostname "localhost" --device-id "<camera-id>" --stream-id "<stream-id>"
    ```

    Example:

    ```bash
    PYTHONPATH=./generated/protos/ds/v1/ python3 main.py --start-time "2025-05-28 16:11" --end-time "2025-05-28 16:14" --external-hostname "localhost" --device-id "2884f669-3b7a-4fbc-8f01-5cbfa8d3279e" --stream-id "28dc44c3-079e-4c94-8ec9-60363451eb40"
    ```

    Output:
    A video recording will be saved in the data folder. You can use VLC or other media players to view the recording.

    Note:
    To export a playback recording, the cameras must have recordings stored in the XProtect media database. Otherwise, the exported video will be empty.

## Annex

To obtain the `camera-id` and `stream-id`, open [AI Bridge GraphiQL](http://localhost:4000/api/bridge/graphql) and run the following query:

```graphql
query get_cameras {
  cameras {
    videoStreams {
        id
    }
  }
}
```

Example Output:

```json
{
  "data": {
    "cameras": [
        {
        "videoStreams": [
          {
            "id": "2884f669-3b7a-4fbc-8f01-5cbfa8d3279e/28dc44c3-079e-4c94-8ec9-60363451eb40"
          }
        ]
      }
    ]
  }
}
```

- To extract the `device-id`, use the part before the `/` (e.g., `2884f669-3b7a-4fbc-8f01-5cbfa8d3279e`).
- To extract the `stream-id`, use the part after the `/` (e.g., `28dc44c3-079e-4c94-8ec9-60363451eb40`).

---

Alternatively, you can use the following commands to get the `device-id` and `stream-id` using `curl` and `jq` (replace `localhost` with your Milestone AI Bridge external hostname):

Send the query to the GraphQL service:

```bash
result=$(curl -f -s -X POST \
-H "Accept: application/json" \
-H "Content-type: application/json" \
--data-binary '{ "query": "{ cameras { videoStreams { id } } }" }' http://localhost:4000/api/bridge/graphql \
| jq .data.cameras[0].videoStreams[0].id \
| tr -d '"')
```

Extract the `device-id` and `stream-id`:

```bash
device_id=$(echo $result | cut -d '/' -f1)
stream_id=$(echo $result | cut -d '/' -f2)
```

Now you can run the program using these values:

```bash
PYTHONPATH=./generated/protos/ds/v1/ python3 main.py --start-time "2025-05-28 16:11" --end-time "2025-05-28 16:14" --external-hostname "localhost" --device-id $device_id  --stream-id $stream_id
```
