import { createStyles, Anchor, rem, Text, Flex } from "@mantine/core";

const useStyles = createStyles((theme) => ({
  footer: {
    borderTop: `${rem(1)} solid ${
      theme.colorScheme === "dark" ? theme.colors.dark[5] : theme.colors.gray[2]
    }`,
  },
}));

export default function Footer() {
  const { classes } = useStyles();

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
        © 2023{" "}
        <Anchor href="https://github.com/0x2E" target="_blank" color="dimmed" underline>
          Rook1e
        </Anchor>
        . All Rights Reserved.
      </Text>
      <Text color="dimmed" size="sm">
        Hosted on{" "}
        <Anchor href="https://vercel.com/" target="_blank" color="dimmed" underline>
          Vercel
        </Anchor>{" "}
        and{" "}
        <Anchor href="https://www.cloudflare.com/" target="_blank" color="dimmed" underline>
          CloudFlare
        </Anchor>
      </Text>
    </Flex>
  );
}
