import React, { createContext, useContext } from "react";
import { useNavigate } from "react-router-dom";
import useLocalStorage from "./useLocalStorage";

const AuthContext = createContext(null);

export const useAuth = () => {
  const [user, setUser] = useLocalStorage("user", null);
  const navigate = useNavigate();

  const clientLogin = async (data) => {
    setUser(data);
    navigate("/home");
  };

  const clientLogout = () => {
    setUser(null);
    navigate("/", { replace: true });
  };
  return {
    user,
    clientLogin,
    clientLogout,
  };
};

export const AuthProvider = ({ children }) => {
  const auth = useAuth();

  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
};

const AuthConsumer = () => useContext(AuthContext);

export default AuthConsumer;
