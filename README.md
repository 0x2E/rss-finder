# rss-finder

**well-known**:

Website url + :

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

**parse HTML**:

- `<link>` with type `application/rss+xml`
- `<link>` with type `application/atom+xml`
- `<link>` with type `application/json`
- `<link>` with type `application/feed+json`
- `<a>` contains `rss` word

## Env variables

| variable     | required | description                                            |
| ------------ | -------- | ------------------------------------------------------ |
| `USER_AGENT` | false    | HTTP `User-Agent` in request, default `rss-finder/1.0` |
