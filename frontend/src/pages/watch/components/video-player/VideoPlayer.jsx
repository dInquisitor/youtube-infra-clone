/* eslint-disable jsx-a11y/media-has-caption */
import React, { useEffect, useRef } from "react";
import dashjs from "dashjs";
import ControlsBox from "./ControlsBox";

const VideoPlayer = ({ videoID }) => {
  const videoElement = useRef(null);

  useEffect(() => {
    if (videoElement.current === null) {
      return;
    }

    const player = dashjs.MediaPlayer().create();
    player.initialize(
      videoElement.current,
      `/api/stream/${videoID}/manifest.mpd`,
      true
    );
  }, [videoElement]);

  return (
    <ControlsBox>
      <video controls ref={videoElement} />
    </ControlsBox>
  );
};

export default VideoPlayer;
