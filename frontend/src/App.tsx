import { FluentProvider, teamsDarkV21Theme } from "@fluentui/react-components";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Outlet } from "react-router";

function App() {
  return (
    <FluentProvider theme={teamsDarkV21Theme}>
      <div className="flex flex-col h-screen w-screen overflow-hidden">
        <Outlet />
      </div>
      <ReactQueryDevtools initialIsOpen={false} position="right" />
    </FluentProvider>
  );
}

export default App;
