import App from "@/App";
import HomePage from "@/home/HomePage";
import "@/index.css";
import { Router } from "@/router";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import { client } from "./api/client.gen";

const queryClient = new QueryClient();

client.setConfig({
  baseUrl: "http://localhost:8080",
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route element={<App />}>
            <Route path={Router.HOME} element={<HomePage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>
);
