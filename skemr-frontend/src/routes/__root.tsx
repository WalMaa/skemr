import { queryClient } from "@/lib/query-client";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Toaster } from "@/components/ui/sonner";
import { ThemeProvider } from "@/components/theme-provider";

export interface RouterContext {
  queryClient: typeof queryClient;
}

const RootLayout = () => {
  return (
    <ThemeProvider storageKey="vite-ui-theme" defaultTheme="system">
      <Outlet />
      <Toaster />
      <TanStackRouterDevtools />
      <ReactQueryDevtools buttonPosition="bottom-left"  initialIsOpen={false} />
    </ThemeProvider>
  );
};

export const Route = createRootRoute({
  component: RootLayout,
  context: () => ({ queryClient }),
});
