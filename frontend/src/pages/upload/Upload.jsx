import React, { useCallback, useState } from "react";
import styled from "styled-components";
import apiClient from "../../requestor/axiosClient";
import HVCenter from "../../shared/HVCenter";
import Page from "../../shared/Page";

const UploadBox = styled.div`
  width: 320px;
`;

const SubmitButton = styled.button`
  background: #c2791a;
  outline: none;
  border: 0;
  cursor: pointer;
  &:disabled {
    background-color: rgba(194, 121, 26, 0.5);
  }
`;

const Upload = () => {
  const [isAttemptingUpload, setIsAttemptingUpload] = useState(false);
  const [uploadErrorMessage, setUploadErrorMessage] = useState(null);
  // const [uploadProgress, setUploadProgress] = useState(null);

  const onUploadProgress = (e) => {
    // can use this for progressbar
    console.log("upload progress", e);
  };

  const uploadFile = useCallback(async (videoFile, videoID) => {
    const formData = new FormData();

    formData.append("videoFile", videoFile);
    formData.append("videoID", videoID);

    apiClient.post("/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      onUploadProgress,
    });
  }, []);

  const attemptUpload = useCallback(async (videoFile) => {
    setIsAttemptingUpload(true);

    if (!videoFile) {
      setUploadErrorMessage("Please select a file to upload");
      setIsAttemptingUpload(false);
      return;
    }

    let beginUploadResponse;

    try {
      beginUploadResponse = await apiClient.post("/begin-upload");
    } catch (e) {
      setUploadErrorMessage("An unknown error occurred, please try again.");
      setIsAttemptingUpload(false);
      return;
    }

    if (beginUploadResponse.status !== 200) {
      setUploadErrorMessage("An unknown error occurred, please try again.");
      setIsAttemptingUpload(false);
      return;
    }

    const beginUploadStatus = beginUploadResponse.data.status;
    if (beginUploadStatus !== "SUCCESS") {
      setUploadErrorMessage("An unknown error occurred, please try again.");
      setIsAttemptingUpload(false);
    }

    // actual upload
    await uploadFile(videoFile, beginUploadResponse.data.videoData.id);

    setIsAttemptingUpload(false);
  }, []);

  const [videoFile, setVideoFile] = useState(null);

  return (
    <Page>
      <HVCenter>
        {uploadErrorMessage}
        <UploadBox>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              attemptUpload(videoFile);
            }}
          >
            <div>Choose file</div>
            <input
              type="file"
              onChange={(e) => setVideoFile(e.target.files[0])}
              disabled={isAttemptingUpload}
            />
            <SubmitButton type="submit" disabled={isAttemptingUpload}>
              Upload video
            </SubmitButton>
          </form>
        </UploadBox>
      </HVCenter>
    </Page>
  );
};

export default Upload;
