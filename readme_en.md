# Dynamic Watermark Video Processing API Documentation

**Version:** 1.0
**Last Updated:** July 25, 2025

## Overview

This service adds a dynamic logo (watermark) that moves from corner to corner every 5 seconds to video files uploaded via HTTP. To avoid blocking the server, operations are executed asynchronously (in the background). When a video process is complete, a POST request is sent to a predefined Webhook URL to notify completion. Processed videos can also be downloaded from or deleted from the server.

It can be run in a Docker environment using the provided `docker-compose.yml` file.

Attention: Please ensure to set the `webhookURL` in the environment variables.

## Base URL

All API requests should be made to the following base URL:

`http://localhost:80` (For local development environment)

## Authentication

There is no authentication method (e.g., API Key, OAuth) in the current version. The endpoints are public. For security, it is recommended to run the API in an isolated network or place it behind an API Gateway.

---

## Endpoints

### 1. Upload Video and Start Processing

Uploads a video and initiates the watermarking process in the background. It returns an immediate response upon accepting the request.

- **Endpoint:** `POST /add-logo`
- **Description:** Creates a new video processing task.
- **Request Type:** `multipart/form-data`

**Request Body:**

| Field Name | Type   | Requirement | Description                                           |
| :--------- | :----- | :---------- | :---------------------------------------------------- |
| `video`    | File   | **Required** | The video file to be processed (mp4, mov, avi, etc.). |
| `video_id` | String | **Required** | A unique identifier for the process.                  |

**Example cURL Request:**

```bash
curl -X POST \
  http://localhost:80/add-logo \
  -F "video=@/path/to/your/local_video.mp4" \
  -F "video_id=a1b2c3d4-e5f6-7890-1234-567890abcdef"
```

Success Response (202-Accepted):

Indicates that the request has been accepted and processing has started in the background.

JSON

{
  "status": "processing",
  "video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
Error Responses:

400 Bad Request (if video file is missing):

JSON

{ "error": "Video file not provided" }
400 Bad Request (if video_id is missing):

JSON

{ "error": "Video ID is required" }
409 Conflict (if a process with the same video_id is already running):

JSON

{ "error": "Video is already being processed" }
2. Download Processed Video
Downloads a video that has completed processing from the server.

Endpoint: POST /get-video

Description: Fetches the processed video by its video_id.

Request Type: application/json

Request Body:

JSON

{
  "video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
Example cURL Request:

Bash

curl -X POST \
  http://localhost:80/get-video \
  -H "Content-Type: application/json" \
  -d '{"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"}' \
  --output processed_video.mp4
Success Response (200-OK):

Content-Type: video/mp4

Body: Raw video file data.

Error Responses:

400 Bad Request (if video_id is missing or invalid):

JSON

{ "error": "Video ID is required" }
404 Not Found (if the video is not found or not yet processed).

3. Delete Processed Video
Permanently deletes a processed video from the server.

Endpoint: DELETE /del-video

Description: Deletes the specified video and its record by video_id.

Request Type: application/json

Request Body:

JSON

{
  "video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
Example cURL Request:

Bash

curl -X DELETE \
  http://localhost:80/del-video \
  -H "Content-Type: application/json" \
  -d '{"video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"}'
Success Response (200-OK):

JSON

{ "message": "Video deleted successfully" }
Error Responses:

400 Bad Request (if video_id is missing or invalid).

404 Not Found (if the video to be deleted is not found).

500 Internal Server Error (if an error occurs during file deletion).

4. Service Health Check (Ping)
Used to check if the service is running.

Endpoint: GET /ping

Example cURL Request:

Bash

curl http://localhost:80/ping
Success Response (200-OK):

Returns a 200 status code with an empty body.

Webhook Notification
When a video processing task is successfully completed, the service sends an HTTP POST request to the webhookURL defined in handlers/video_upload.go in the following format.

Method: POST

Content-Type: application/json

Webhook Body Content:

JSON

{
  "video_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
This notification indicates that the process is complete and the video is now available for download or deletion.