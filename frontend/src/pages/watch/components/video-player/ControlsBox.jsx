import React from "react";
import styled from "styled-components";
// import PlayButtonSrc from "./assets/play-button.svg";
// import PauseButtonSrc from "./assets/pause-button.svg";

const Container = styled.div`
  width: 848px;
  height: 480px;
  /* position: relative; */
  & > video {
    width: 100%;
    height: 100%;
  }
`;

// const ControlBar = styled.div`
//   background: rgba(0, 0, 0, 0.5);
//   width: 100%;
//   height: 36px;
//   z-index: 1;
//   position: absolute;
//   bottom: 0;
//   left: 0;
//   padding: 6px 12px;
//   display: flex;
//   align-items: center;
// `;

// const PlayPauseButton = styled.img`
//   width: 24px;
//   height: 24px;
//   cursor: pointer;
// `;

// const PlayTime = styled.div`
//   margin-left: 12px;
// `;

// const ProgressContainer = styled.div`
//   width: 60%;
//   background: #aaa;
//   margin-left: 36px;
//   height: 12px;
//   cursor: pointer;
// `;

// const PlayedProgress = styled.div`
//   background: #fff;
//   width: ${(props) => `${props.percentage || 0}%`};
//   height: 100%;
// `;

// const padTime = (string, pad, length) =>
//   (new Array(length + 1).join(pad) + string).slice(-length);

// const videoTimeFormat = (secondsFloat) => {
//   const minutes = Math.floor(secondsFloat / 60);
//   const seconds = Math.floor(secondsFloat - minutes * 60);
//   return `${padTime(minutes, "0", 2)}:${padTime(seconds, "0", 2)}`;
// };

// get videoPlayer as props from caller
const ControlsBox = ({ children }) => (
  // const [isPaused, setIsPaused] = useState(true);
  // const [playBackTime, setPlaybackTime] = useState(null);

  // useEffect(() => {
  //   if (!videoPlayer) {
  //     return;
  //   }
  //   videoPlayer.on("playbackTimeUpdated", (e) => {
  //     setPlaybackTime(e.time);
  //     console.log(videoPlayer.getDashMetrics().getCurrentBufferLevel());
  //   });
  // }, [setPlaybackTime, videoPlayer]);

  // const togglePlay = useCallback(() => {
  //   if (isPaused) {
  //     videoPlayer.play();
  //     setIsPaused(false);
  //     return;
  //   }
  //   videoPlayer.pause();
  //   setIsPaused(true);
  // }, [videoPlayer, isPaused, setIsPaused]);

  <Container>
    {children}
    {/* {videoPlayer && (
        <ControlBar>
          <PlayPauseButton
            src={isPaused ? PlayButtonSrc : PauseButtonSrc}
            alt=""
            onClick={togglePlay}
          />
          <PlayTime>
            {videoTimeFormat(playBackTime < 1 ? 0 : playBackTime)} /{" "}
            {!Number.isNaN(videoPlayer.duration()) &&
              videoTimeFormat(videoPlayer.duration())}
          </PlayTime>
          <ProgressContainer>
            <PlayedProgress
              percentage={Math.floor(
                (playBackTime / videoPlayer.duration()) * 100
              )}
            />
          </ProgressContainer>
        </ControlBar>
      )} */}
  </Container>
);
export default ControlsBox;
