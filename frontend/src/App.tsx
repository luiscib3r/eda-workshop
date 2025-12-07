import { FluentProvider, teamsDarkV21Theme } from "@fluentui/react-components";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Link, Outlet } from "react-router";

function App() {
  return (
    <FluentProvider theme={teamsDarkV21Theme}>
      <div className="flex flex-col h-screen w-screen overflow-hidden">
        <div className="flex items-center w-full h-12 bg-[#0F0F0F] px-6">
          <Link to="/">EDA Workshop</Link>
        </div>
        <Outlet />
      </div>
      <ReactQueryDevtools initialIsOpen={false} position="right" />
    </FluentProvider>
  );
}

export default App;
