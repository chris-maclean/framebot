# Framebot

## What is Framebot?
Framebot is a automated Twitter bot the posts frames from a video file, in sequence, once every 5 minutes. It's a great way to bring beautiful images from your favorite film to your timeline. We recommend [running Framebot in a Docker container](#docker-usage).

## Quickstart
Start by cloning the [Github repo](https://github.com/chris-maclean/framebot).

To post the first frame of a video file, save the following file to your filesystem
```json
{
    "file": "/path/to/your/video/file.mp4",
    "title": "Your Video File's Title!",
    "nextFrame": 1
}
```

Build framebot
```bash
$ go build -v -o . ./...
```

Set environment variables representing OAuth keys:
```bash
export TWITTER_CONSUMER_KEY=...
export TWITTER_CONSUMER_SECRET=...
export TWITTER_ACCESS_TOKEN=...
export TWITTER_ACCESS_SECRET=...
```

Run the following command:

```bash
$ framebot --stateFile /path/to/state/file.json
```

Check the Twitter profile of the account associated with the OAuth keys. A new tweet should have been posted containing the first frame of the video with text showing the title and frame number.

## Dependencies
Framebot relies on `ffmpeg` to extract frame images from the video file. This library is built into the container image. Non-container users should have the `ffmpeg` library on their machine; and `ffmpeg`, `ffprobe`, and `bash` available on their PATH. The user is responsible for providing the video file.

Framebot also uses the Twitter API to post tweets containing the frame image and a short bit of text. The user must have a [Twitter Developer Platform](https://developer.twitter.com/en) account and provide OAuth user context keys for the account making the tweets. Those keys must be provided in environment variables under the following names:

```bash
TWITTER_CONSUMER_KEY
TWITTER_CONSUMER_SECRET
TWITTER_ACCESS_TOKEN
TWITTER_ACCESS_SECRET
```

## Running Framebot
Framebot posts frames from the desired video file in sequence, but the Framebot program itself is not a long-running process. Each invocation of Framebot performs the following steps, in order:
* read a state file from the filesystem
* use the state information to identify the frame it needs to post
* extract the frame image
* post the tweet including the image
* update state information to indicate the next frame to be posted
* write the state file back to the filesystem

All of the information Framebot needs to run is contained in the state file. 

### The State File
The Framebot state file is a JSON object that includes up to 4 properties. The user is responsible for creating the first state file, and after that Framebot will generate updated values and save the file as part of its run.

A nominal state file looks like this:
```json
{
    "file": "/opt/framebot/2001.mp4",
    "title": "2001: A Space Odyssey (1968)",
    "nextFrame": 17492,
    "totalFrames": 214154
}
```
* `file`: fully-qualified path to the video file
* `title`: the title of the video. This is used only as text in the posted tweets
* `nextFrame`: the frame of the video file that will be posted on the next execution of Framebot. The Framebot code increments this value on each execution
* `totalFrames`: the total number of frames in the video file. This is used only as text in the posted tweets.

All values are optional! The default behavior for each property is as follows:
* `file`: the file location will be set to `/opt/framebot/movie`
* `title`: the title of the video will not be printed in the tweet text
* `nextFrame`: frame #1 will be extracted and posted
* `totalFrames`: the total number of frames will be calculated with `ffprobe` and saved for future executions. 

The call to `ffprobe` to calculate the total number of frames can be very slow! If the user has prior knowledge of the number of frames, or a faster algorithm for calculating it, they should provide `totalFrames` in the initial state file and prevent Framebot from calculating it. However, even if `totalFrames` is not known at the first execution, the calculation will only be performed once. After that, the value will be saved to the state file for future executions.


## Runtime options
* `--stateFile`: a path to the state file. Defaults to `/opt/framebot/framebot-state.json` if not provided

## Docker Usage
We highly recommend running the Framebot image in a Docker container. The image includes the necessary `ffmpeg` library, a prebuilt executable, and a `crontab` entry that automatically executes Framebot once every 5 minutes. 

The user will need to specify some options when starting the container image:
* `-v /path/to/initialStateFile.json:/opt/framebot/framebot-state.json`
    * Mount an initial state file to the location `/opt/framebot/framebot-state.json`. This location on the container filesystem is required because `framebot` will be invoked without any runtime options. Ensure that the `file` property of the state file represents the filepath _on the container filesystem!_
* `-v /path/to/moviefile:/opt/framebot/movie`
    * Mount the video file from which frames will be extracted. The container file location must match the value specified in the state file, or be the default `/opt/framebot/movie` if the state file does not specify a filepath
* `--env-file /path/to/envfile`
    * Attach environment variables representing the Twitter account's OAuth keys. See [Dependencies](#Dependencies)

The container image automatically adds a Framebot execution to the root user's `crontab`. This ensures that Framebot will be posting indefinitely as long as the container is alive and the video file has frames remaining. This is the best way to use Framebot!

## Thank you for using Framebot!
