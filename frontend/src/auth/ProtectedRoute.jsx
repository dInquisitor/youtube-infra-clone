import { Navigate } from "react-router-dom";
import React from "react";
import { useAuth } from "./useAuth";

const ProtectedRoute = ({ children }) => {
  const { user } = useAuth();
  if (!user) {
    return <Navigate to="/" />;
  }
  return children;
};

export default ProtectedRoute;
