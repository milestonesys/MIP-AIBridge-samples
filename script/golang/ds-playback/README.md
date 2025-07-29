# DS-Playback

This sample script demonstrates how to use the Direct Streaming protocol to fetch playback recordings from XProtect. The Direct Streaming protocol is based on gRPC. The proto file is currently published in the `AI Bridge developer reference` at [http://<ai-bridge-external-hostname>:4000/](http://<ai-bridge-external-hostname>:4000/).

The proto file can be downloaded from `http://<ai-bridge-external-hostname>:4000/`. Make sure that Milestone AI Bridge is running in debug mode.

## Getting Started

The directory contains the following files:

```txt
.
├── README.md
└── src
    ├── Makefile
    ├── go.mod
    ├── go.sum
    └── main.go

1 directory, 5 files
```

1. **Set up your Go environment.**

    Download and install Go by following the official guide: [Download and install Go](https://go.dev/doc/install).

2. **Install required packages for Makefile.**

    On Ubuntu, run:

    ```bash
    sudo apt-get -y install wget make curl jq
    ```

3. **Install dependencies and download required files.**

    In this directory, run:

    ```bash
    make -C src/ setup SYSTEM_IP=<ai-bridge-external-hostname>
    ```

    When the setup is finished:

    - Go dependencies for compiling the proto file and generating the client code are located in `./bin`.
    - The proto file will be downloaded to [./src/bin/protos/ds/v1/ds.proto](./src/bin/protos/ds/v1/ds.proto).
    - Generated Go code for the playback client will be located at [./src/generated/protos/ds/v1/ds.pb.go](./src/generated/protos/ds/v1/ds.pb.go).

4. **Build the command-line app.**

    Run:

    ```bash
    make -C src/ build
    ```

    The binary is located in the `./src` directory.

5. **View help for usage instructions.**

    Run:

    ```bash
    ./src/bin/ds-playback --help
    ```

    Output:

    ```txt
    Usage of ./bin/ds-playback:
      -device-id string
            The camera ID used to extract the video recording
      -end-time string
            The end time of the recording you want to fetch (default "yyyy-mm-dd 24:00")
      -external-hostname string
            Hostname or IP address of the machine running Milestone AI Bridge in debug mode (default "localhost")
      -start-time string
            The start time of the recording you want to fetch (default "yyyy-mm-dd 24:00")
      -stream-id string
            The stream ID of the specified device
    ```

6. **Extract a video recording from XProtect.**

    Run:

    ```bash
    ./src/bin/ds-playback -start-time "<yyyy-mm-dd HH:MM>" -end-time "<yyyy-mm-dd HH:MM>" -external-hostname "localhost" -device-id "<camera-id>" -stream-id "<stream-id>"
    ```

    Example:

    ```bash
    ./src/bin/ds-playback -start-time "2025-05-28 16:11" -end-time "2025-05-28 16:14" -external-hostname "localhost" -device-id "2884f669-3b7a-4fbc-8f01-5cbfa8d3279e" -stream-id "28dc44c3-079e-4c94-8ec9-60363451eb40"
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
./src/bin/ds-playback -start-time "2025-05-28 16:11" -end-time "2025-05-28 16:14" -external-hostname "localhost" -device-id $device_id -stream-id $stream_id
```
