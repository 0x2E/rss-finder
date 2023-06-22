import { TextInput, ActionIcon } from "@mantine/core";
import { IconArrowRight, IconLink, IconX } from "@tabler/icons-react";
import List, { feed } from "./list";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { useForm } from "@mantine/form";

interface dataRespI {
  data: feed[];
}
interface errRespI {
  error: string;
}

const urlDefault = "https://rook1e.com";

function normalizeURL(url: string): string {
  if (url == "") {
    url = urlDefault;
  } else if (!url.startsWith("http://") && !url.startsWith("https://")) {
    url = "https://" + url;
  }
  return url;
}

export default function Form() {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<{ data: feed[] }>();

  const form = useForm({
    initialValues: {
      url: "",
    },

    validate: {
      url: (url: string) => {
        url = normalizeURL(url);
        form.setFieldValue("url", url);
        try {
          new URL(url);
        } catch (error) {
          return "Invalid URL";
        }
        return null;
      },
    },
    transformValues: (values) => ({
      url: normalizeURL(values.url),
    }),
  });

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
          case 504:
            throw new TypeError("timeout");
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
    <>
      <form onSubmit={form.onSubmit(submit)}>
        <TextInput
          icon={<IconLink />}
          size="lg"
          placeholder={urlDefault}
          // label="URL"
          // withAsterisk
          autoFocus
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
    </>
  );
}
