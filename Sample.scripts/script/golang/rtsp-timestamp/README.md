# RTSP-Timestamp

This sample script shows how to play an RTSP live feed from cameras connected to an AI Bridge installation and retrieve the timestamp from the RTSP header.

## Getting Started

The directory contains the following files:

```txt
.
├── README.md
└── src
    ├── go.mod
    ├── go.sum
    ├── Makefile
    └── main.go

1 directory, 5 files
```

1. **Set up your Go environment.**

    Download and install Go by following the official guide: [Download and install Go](https://go.dev/doc/install).

2. **Build the command line app.**

    Run the following command:

    ```bash
    make -C src/ build
    ```

    You can find the binary at the `./src` directory.

3. **View help for usage instructions.**

    Run the following command:

    ```bash
    ./src/bin/rtsp-timestamp --help
    ```

    Output:

    ```txt
    Usage of ./src/bin/rtsp-timestamp:
        -rtsp-endpoint string
            The rtsp endpoint from which the video comes from
    ```

4. **Extract a video recording from XProtect.**

    Run the following command:

    ```bash
    ./src/bin/rtsp-timestamp -rtsp-endpoint "<RTSP Endpoint>"
    ```

    Example:

    ```bash
    ./src/bin/rtsp-timestamp -rtsp-endpoint "rtsp://myhostname:8554/9d9c0802-5011-4a25-83ed-536c1f6680c0/28dc44c3-079e-4c94-8ec9-60363451eb40"
    ```

    Output:
    The console output shows the timestamp. You must stop the script manually, for example, by pressing Ctrl+C.

## Annex

### Obtaining an RTSP endpoint

Open [AI Bridge GraphiQL](http://localhost:4000/api/bridge/graphql) and run the following query:

```graphql
query get_rtsp_endpoints {
  cameras {
    name
    videoStreams {
      name
      streamAvailability {
        rtsp
      }
    }
  }
}
```

Example output:

```json
{
  "data": {
    "cameras": [
      {
        "name": "AXIS Q6054 Mk II PTZ Dome Network Camera (10.1.65.24) - Camera 1",
        "videoStreams": [
          {
            "name": "Video stream 1",
            "streamAvailability": {
              "rtsp": "rtsp://myhostname:8554/9d9c0802-5011-4a25-83ed-536c1f6680c0/28dc44c3-079e-4c94-8ec9-60363451eb40"
            }
          }
        ]
      }
    ]
  }
}
```

The value of the `"rtsp"` field is the RTSP endpoint which you can use for the `rtsp-endpoint` parameter.