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

- atom.xml
- feed.xml
- rss.xml
- index.xml
- atom.json
- feed.json
- rss.json
- index.json
- feed/
- rss/

**Third party services**:

- GitHub: Based on [official rules](https://docs.github.com/en/rest/activity/feeds?apiVersion=2022-11-28)
- Reddit: Based on [official wiki](https://www.reddit.com/wiki/rss/)

## Env variables

| variable     | required | description                                            |
| ------------ | -------- | ------------------------------------------------------ |
| `USER_AGENT` | false    | HTTP `User-Agent` in request, default `rss-finder/1.0` |
