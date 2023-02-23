import React from "react";
import styled from "styled-components";

const Container = styled.div`
  width: 848px;
  height: 480px;
  & > video {
    width: 100%;
    height: 100%;
  }
`;

const ControlsBox = ({ children }) => <Container>{children}</Container>;
export default ControlsBox;
