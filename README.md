# RSS Finder

A modern web application to discover RSS feeds from any website. Powered by feedfinder-ts for reliable RSS feed detection.

## Features

- **Robust Detection**: Uses [feedfinder-ts](https://github.com/0x2E/feedfinder-ts) under the hood for comprehensive RSS feed discovery
- **Fast & Lightweight**: Optimized for performance with minimal dependencies

## Getting Started

### Prerequisites

- Node.js 22+
- pnpm (recommended) or npm

### Installation

```bash
# Clone the repository
git clone <your-repo-url>
cd rss-finder

# Install dependencies
pnpm install
```

### Development

```bash
# Start development server
pnpm dev
```

Visit `http://localhost:5173` to see the application.

### Building

```bash
# Build for production
pnpm build
```

## Deployment

Since this app is built with SvelteKit, it's compatible with multiple runtimes. Check the [SvelteKit docs](https://svelte.dev/docs/kit/building-your-app) for platform-specific instructions.

## API

### Find RSS Feeds

**POST** `/api/find-feeds`

```json
{
	"url": "https://example.com"
}
```

**Response:**

```json
{
	"feeds": [
		{
			"title": "Example Blog RSS",
			"url": "https://example.com/rss.xml",
			"type": "application/rss+xml"
		}
	]
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT

## Credits

- [feedfinder-ts](https://github.com/0x2E/feedfinder-ts) - For reliable RSS feed detection
