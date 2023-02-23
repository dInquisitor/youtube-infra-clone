import React from "react";
import { Link } from "react-router-dom";
import Page from "../../shared/Page";

const Home = () => (
  <Page>
    <Link to="/upload">Upload Video</Link>
  </Page>
);

export default Home;
