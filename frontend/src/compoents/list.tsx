import {
  ActionIcon,
  CopyButton,
  Group,
  Table,
  Text,
  Tooltip,
  createStyles,
} from "@mantine/core";
import { IconCheck, IconCopy } from "@tabler/icons-react";

const useStyles = createStyles((theme) => ({
  table: {
    wordBreak: "break-word",
  },
}));

export interface feed {
  title: string;
  link: string;
}

export default function List(props: { data: feed[] }) {
  const { classes } = useStyles();

  if (props.data.length == 0) {
    return <Text fz={"lg"}>Oops, I didn't find any rss links.</Text>;
  }

  return (
    <Table
      striped
      highlightOnHover
      horizontalSpacing="sm"
      verticalSpacing="sm"
      fontSize="md"
      className={classes.table}
    >
      <thead>
        <tr>
          <th>#</th>
          <th>Title</th>
          <th>Link</th>
        </tr>
      </thead>
      <tbody>
        {props.data.map((feed, index) => (
          <tr key={feed.link}>
            <td>{index + 1}</td>
            <td>{feed.title}</td>
            <td>
              <Group spacing="xs">
                {feed.link}
                <CopyButton value={feed.link} timeout={3000}>
                  {({ copied, copy }) => (
                    <Tooltip
                      label={copied ? "Copied" : "Copy"}
                      withArrow
                      position="right"
                    >
                      <ActionIcon
                        color={copied ? "teal" : "gray"}
                        onClick={copy}
                      >
                        {copied ? (
                          <IconCheck size="1rem" />
                        ) : (
                          <IconCopy size="1rem" />
                        )}
                      </ActionIcon>
                    </Tooltip>
                  )}
                </CopyButton>
              </Group>
            </td>
          </tr>
        ))}
      </tbody>
    </Table>
  );
}
