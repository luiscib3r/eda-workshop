import App from "@/App";
import HomePage from "@/home/HomePage";
import "@/index.css";
import { Router } from "@/router";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import { client } from "./api/client.gen";
import FilePagesPage from "./file_pages/FilePagesPage";

const queryClient = new QueryClient();

client.setConfig({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route element={<App />}>
            <Route path={Router.HOME} element={<HomePage />} />
            <Route path={Router.FILE_PAGES} element={<FilePagesPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>
);
