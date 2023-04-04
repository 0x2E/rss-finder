import { createStyles, Text, rem, Button, Group } from "@mantine/core";
import { IconBrandGithub, IconExternalLink } from "@tabler/icons-react";

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

export default function Header() {
  const { classes } = useStyles();
  return (
    <>
      <h1 className={classes.title}>
        Another{" "}
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
        Enter the url to get its feed link automatically. Based on webpage
        parsing and well known paths.
      </Text>
      <Group align="center">
        <Button
          component="a"
          target="_blank"
          rel="noopener noreferrer"
          href="https://github.com/0x2E/rss-finder#how-it-works"
          variant="light"
          color="gray"
          leftIcon={<IconExternalLink size={rem(18)} />}
        >
          How it works
        </Button>
        <Button
          component="a"
          target="_blank"
          href="https://github.com/0x2E/rss-finder"
          variant="light"
          color="gray"
          leftIcon={<IconBrandGithub size={rem(18)} />}
        >
          Source Code
        </Button>
      </Group>
    </>
  );
}
