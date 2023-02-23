import React from "react";
import "./App.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "react-query";
import SignIn from "./pages/sign-in/SignIn";
import ProtectedRoute from "./auth/ProtectedRoute";
import Home from "./pages/home/Home";
import Upload from "./pages/upload/Upload";
import SignOut from "./pages/sign-out/SignOut";
import Watch from "./pages/watch/Watch";

const queryClient = new QueryClient();

const App = () => {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <SignIn />,
    },
    {
      path: "/sign-out",
      element: <SignOut />,
    },
    {
      path: "/home",
      element: (
        <ProtectedRoute>
          <Home />
        </ProtectedRoute>
      ),
    },
    {
      path: "/upload",
      element: (
        <ProtectedRoute>
          <Upload />
        </ProtectedRoute>
      ),
    },
    {
      path: "/watch",
      element: <Watch />,
    },
  ]);
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  );
};

export default App;
