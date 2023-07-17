import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { Notifications } from "@mantine/notifications";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <Notifications position="top-right" />
    <App />
  </React.StrictMode>,
);
