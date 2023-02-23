/* eslint-disable jsx-a11y/media-has-caption */

import React from "react";
import { useSearchParams } from "react-router-dom";
import styled from "styled-components";
import Page from "../../shared/Page";
import VideoPlayer from "./components/video-player/VideoPlayer";

const WatchContainer = styled.div`
  padding: 5% 10%;
`;

const Watch = () => {
  const [searchParams] = useSearchParams();

  return (
    <Page>
      <WatchContainer>
        {searchParams.get("v") && (
          <VideoPlayer videoID={searchParams.get("v")} />
        )}
        Comments here?
      </WatchContainer>
    </Page>
  );
};

export default Watch;
