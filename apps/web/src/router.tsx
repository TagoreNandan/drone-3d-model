import { createBrowserRouter } from "react-router";
import { lazy, Suspense } from "react";

import AppShell from "./app-shell";
import Home from "./routes/home";

const DroneView = lazy(() => import("./routes/drone-view"));

function NotFound() {
  return (
    <main className="container mx-auto max-w-3xl px-4 py-8">
      <h1 className="text-2xl font-semibold">404</h1>
      <p className="text-muted-foreground">The requested page could not be found.</p>
    </main>
  );
}

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppShell />,
    children: [
      { index: true, element: <Home /> },
      {
        path: "drone-view",
        element: (
          <Suspense fallback={<div className="flex h-full items-center justify-center p-8">Loading...</div>}>
            <DroneView />
          </Suspense>
        ),
      },
      { path: "*", element: <NotFound /> },
    ],
  },
]);
