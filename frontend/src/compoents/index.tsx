import {
  createStyles,
  Container,
  Text,
  rem,
  TextInput,
  ActionIcon,
  Stack,
} from "@mantine/core";
import { IconArrowRight, IconLink, IconX } from "@tabler/icons-react";
import List, { feed } from "./list";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { useForm } from "@mantine/form";
import Footer from "./footer";

const useStyles = createStyles((theme) => ({
  title: {
    fontSize: rem(62),
    fontWeight: 900,
    lineHeight: 1.1,
    margin: 0,
    padding: 0,

    [theme.fn.smallerThan("sm")]: {
      fontSize: rem(42),
      lineHeight: 1.2,
    },
  },

  description: {
    fontSize: rem(24),

    [theme.fn.smallerThan("sm")]: {
      fontSize: rem(18),
    },
  },

  control: {
    height: rem(54),
    paddingLeft: rem(38),
    paddingRight: rem(38),

    [theme.fn.smallerThan("sm")]: {
      height: rem(54),
      paddingLeft: rem(18),
      paddingRight: rem(18),
      flex: 1,
    },
  },
}));

export default function Index() {
  const { classes } = useStyles();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<{ data: feed[] }>();

  const form = useForm({
    initialValues: {
      url: "",
    },

    validate: {
      url: (url: string) => {
        if (
          url != "" &&
          !url.startsWith("http://") &&
          !url.startsWith("https://")
        ) {
          url = "https://" + url;
          form.setFieldValue("url", url);
        }
        try {
          new URL(url);
        } catch (error) {
          return "Invalid URL";
        }
        return null;
      },
    },
  });

  interface dataRespI {
    data: feed[];
  }
  interface errRespI {
    error: string;
  }

  async function submit(value: { url: string }) {
    setLoading(true);
    setData({ data: [] });
    try {
      const resp = await fetch(
        import.meta.env.VITE_HOST +
          "/api/find?url=" +
          encodeURIComponent(value.url),
        {
          method: "GET",
        }
      );
      if (!resp.ok) {
        switch (resp.status) {
          case 400:
            const errData: errRespI = await resp.json();
            throw new TypeError(errData.error);
            break;
          case 500:
            throw new TypeError("some errors occurred on the server side");
            break;
          default:
            throw new TypeError(resp.statusText);
        }
      }
      let respData: dataRespI = await resp.json();
      console.log(respData);
      let tmpData: dataRespI = { data: [] };
      for (const fr of respData.data) {
        let exist = false;
        for (let ft of tmpData.data) {
          if (ft.title == fr.title && ft.link == fr.link) {
            exist = true;
            break;
          }
        }
        if (!exist) {
          tmpData.data.push(fr);
        }
      }
      setData(tmpData);
    } catch (error) {
      console.log(error);
      if (error instanceof Error) {
        notifications.show({
          title: "Error",
          message: error.message,
          icon: <IconX size="1.1rem" />,
          color: "red",
          autoClose: 5000,
        });
      }
    }
    setLoading(false);
  }
  return (
    <Stack justify="space-between" mih="100vh">
      <Container size={700}>
        <Stack pt="12rem">
          <h1 className={classes.title}>
            A simple{" "}
            <Text
              component="span"
              variant="gradient"
              gradient={{ from: "blue", to: "cyan" }}
              inherit
            >
              RSS Finder
            </Text>
          </h1>

          <Text className={classes.description} color="dimmed">
            Enter a website url to automatically find its subscription links
            from pages and common paths
          </Text>
          <form onSubmit={form.onSubmit(submit)}>
            <TextInput
              icon={<IconLink />}
              size="lg"
              placeholder="https://rook1e.com"
              // label="URL"
              // withAsterisk
              rightSection={
                <ActionIcon
                  size={32}
                  radius="xl"
                  color={"blue"}
                  variant="filled"
                  type="submit"
                  loading={loading}
                >
                  <IconArrowRight size="1.1rem" stroke={1.5} />
                </ActionIcon>
              }
              {...form.getInputProps("url")}
            />
          </form>

          {!loading && data && <List data={data.data} />}
        </Stack>
      </Container>
      <Footer />
    </Stack>
  );
}
