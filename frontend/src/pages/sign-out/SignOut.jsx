import React, { useEffect } from "react";
import { useAuth } from "../../auth/useAuth";
import apiClient from "../../requestor/axiosClient";
import Page from "../../shared/Page";

const SignOut = () => {
  const { clientLogout } = useAuth();

  useEffect(() => {
    (async () => {
      await apiClient.get("/sign-out");
      clientLogout();
    })();
  }, []);
  return <Page />;
};

export default SignOut;
