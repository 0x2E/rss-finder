import { MantineProvider } from "@mantine/core";
import Index from "./compoents/index";
import "./App.css";

export default function App() {
  return (
    <MantineProvider withGlobalStyles withNormalizeCSS>
      <Index />
    </MantineProvider>
  );
}
