import { createStyles, Anchor, rem, Text, Flex } from "@mantine/core";

const useStyles = createStyles((theme) => ({
  footer: {
    borderTop: `${rem(1)} solid ${
      theme.colorScheme === "dark" ? theme.colors.dark[5] : theme.colors.gray[2]
    }`,
  },
}));

interface link {
  label: string;
  link: string;
}

const links: link[] = [
  { label: "GitHub", link: "https://github.com/0x2E" },
  {
    label: "Rook1e",
    link: "https://rook1e.com",
  },
];

export default function Footer() {
  const { classes } = useStyles();
  const items = links.map((link) => (
    <Anchor<"a">
      color="dimmed"
      key={link.label}
      href={link.link}
      onClick={(event) => event.preventDefault()}
      size="sm"
    >
      {link.label}
    </Anchor>
  ));

  return (
    <Flex
      className={classes.footer}
      justify="space-between"
      align="center"
      direction={{ base: "column", sm: "row" }}
      gap={{ base: "xs" }}
      wrap="wrap"
      py={{ base: "sm", sm: "xl" }}
      px={{ base: "lg", sm: "5rem" }}
    >
      <Text color="dimmed" size="sm">
        Host on Vercel & Azure
      </Text>
      <Flex
        justify="space-between"
        align="center"
        direction={{ base: "column", sm: "row" }}
        wrap="wrap"
        gap={{ base: "xs", sm: "lg" }}
      >
        {items}
      </Flex>
    </Flex>
  );
}
