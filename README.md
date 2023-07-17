# [RSS Finder](https://rss-finder.rook1e.com/)

A tool for finding and sniffing rss links.

## How It Works

**Parsing HTML**:

- `<link>` with type `application/rss+xml`
- `<link>` with type `application/atom+xml`
- `<link>` with type `application/json`
- `<link>` with type `application/feed+json`
- `<a>` contains `rss` word

**Well-known paths**:

- `atom.xml`, `feed.xml`, `rss.xml`, `index.xml`
- `atom.json`, `feed.json`, `rss.json`, `index.json`
- `feed/`, `rss/`

**Third party services**:

- GitHub: [official rules](https://docs.github.com/en/rest/activity/feeds?apiVersion=2022-11-28)
- Reddit: [official wiki](https://www.reddit.com/wiki/rss/)
- YouTube: [ref](https://authory.com/blog/create-a-youtube-rss-feed-with-vastly-increased-limits)

## Contributing

Please keep code clean, and thanks for your contribution!

1. frontend: front-end code is in `frontend`

```bash
npm run dev
```

2. serverless: use [vercel cli](https://vercel.com/docs/cli) to run locally

```bash
vercel dev
```

3. Test the changes, e.g. `go test . /... ` to test Go code. Then open a pr to the main branch. It is recommended that one pr does only one thing.

## Env variables

| variable     | required | description                                            |
| ------------ | -------- | ------------------------------------------------------ |
| `USER_AGENT` | false    | HTTP `User-Agent` in request, default `rss-finder/1.0` |
| `WEB_DOMAIN` | true     | Domain name of the web ui                              |
