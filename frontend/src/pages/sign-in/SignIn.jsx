import React, { useCallback, useState } from "react";
import { Navigate } from "react-router-dom";
import styled from "styled-components";
import { useAuth } from "../../auth/useAuth";
import apiClient from "../../requestor/axiosClient";
import HVCenter from "../../shared/HVCenter";
import Page from "../../shared/Page";

const LoginBox = styled.div`
  width: 320px;
`;

const Label = styled.div`
  margin-bottom: 6px;
`;

const TextInput = styled.input`
  background: transparent;
  border: 1px solid #ffffff;
  outline: none;
  box-shadow: none;
  display: block;
  color: #ffffff;
  padding: 10px;
  margin-bottom: 12px;
  width: 100%;
`;

const SubmitButton = styled.button`
  background: #c2791a;
  display: block;
  outline: none;
  border: 0;
  padding: 12px 16px;
  width: 100%;
  cursor: pointer;
  &:disabled {
    background-color: rgba(194, 121, 26, 0.5);
  }
`;

const SignInError = styled.div`
  color: #ffffff;
  margin-bottom: 12px;
  background: #b95959;
  padding: 6px 12px;
  border-radius: 2px;
  font-size: 0.8em;
`;

const UNKNOWN_ERROR_MESSAGE = "An unknown error occurred, please try again.";

const SignIn = () => {
  const { user, clientLogin } = useAuth();

  if (user) {
    return <Navigate to="/home" />;
  }

  const [signInErrorMessage, setSignInErrorMessage] = useState(null);
  const [isAttemptingSignIn, setIsAttemptingSignIn] = useState(false);

  const attemptSignIn = useCallback(async (email, password) => {
    setIsAttemptingSignIn(true);
    setSignInErrorMessage(null);

    let signInResponse;

    try {
      signInResponse = await apiClient.post("/sign-in", {
        data: {
          email,
          password,
        },
      });
    } catch (e) {
      console.log(e);
      setSignInErrorMessage(UNKNOWN_ERROR_MESSAGE);
      setIsAttemptingSignIn(false);
      return;
    }

    if (signInResponse.status !== 200) {
      setSignInErrorMessage(UNKNOWN_ERROR_MESSAGE);
      setIsAttemptingSignIn(false);
      return;
    }

    const signInStatus = signInResponse.data.status;

    if (signInStatus === "SUCCESS") {
      clientLogin(signInResponse.data.userData);
      setIsAttemptingSignIn(false);
      return;
    }

    if (signInStatus === "LOGIN_INFO_INCORRECT") {
      setSignInErrorMessage("Email or password is incorrrect");
      setIsAttemptingSignIn(false);
      return;
    }

    setSignInErrorMessage(UNKNOWN_ERROR_MESSAGE);
    setIsAttemptingSignIn(false);
  }, []);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  return (
    <Page>
      <HVCenter>
        <LoginBox>
          {signInErrorMessage && (
            <SignInError>{signInErrorMessage}</SignInError>
          )}

          <form
            onSubmit={(e) => {
              e.preventDefault();
              attemptSignIn(email, password);
            }}
          >
            <Label>Email address</Label>
            <TextInput
              type="text"
              value={email}
              onChange={({ target }) => setEmail(target.value)}
            />
            <Label>Password</Label>
            <TextInput
              type="password"
              value={password}
              onChange={({ target }) => setPassword(target.value)}
            />
            <SubmitButton type="submit" disabled={isAttemptingSignIn}>
              Sign In / Up
            </SubmitButton>
          </form>
        </LoginBox>
      </HVCenter>
    </Page>
  );
};

export default SignIn;
