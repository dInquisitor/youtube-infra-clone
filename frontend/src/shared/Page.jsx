import React from "react";
import styled from "styled-components";

const Root = styled.div`
  background: #292d3e;
  color: #ffffff;
  height: 100vh;
  font-size: 16px;
  & input,
  button {
    font-size: 16px;
    color: #ffffff;
  }
`;

const Page = (props) => <Root {...props} />;

export default Page;
