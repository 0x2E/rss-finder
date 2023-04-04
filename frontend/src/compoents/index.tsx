import { Container, Stack } from "@mantine/core";
import Footer from "./footer";
import Header from "./header";
import Form from "./form";

export default function Index() {
  return (
    <Stack justify="space-between" mih="100vh">
      <Container size={800}>
        <Stack pt="11rem">
          <Header />
          <Form />
        </Stack>
      </Container>
      <Footer />
    </Stack>
  );
}
