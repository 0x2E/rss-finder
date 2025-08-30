import { Hono } from "hono";
import { renderer } from "./renderer";
import { find } from "feedfinder-ts";

const app = new Hono();

app.use(renderer);

app.get("/", (c) => {
  return c.render(
    <main class="container">
      <section class="main-content">
        <header>
          <h1>RSS Finder</h1>
          <p class="subtitle">Discover RSS feeds from any website</p>
        </header>

        <form
          class="search-form"
          id="searchForm"
          role="search"
          aria-label="RSS feed search"
        >
          <div class="input-group">
            <label for="urlInput" class="visually-hidden">
              Website URL
            </label>
            <input
              type="url"
              id="urlInput"
              name="url"
              placeholder="Enter a website URL..."
              aria-describedby="url-help"
              required
              autocomplete="url"
            />
            <button type="submit" id="searchBtn" aria-describedby="search-help">
              <span class="button-text">Find RSS</span>
            </button>
          </div>
          <small id="url-help" class="help-text">
            Enter any website URL to discover its RSS feeds
          </small>
        </form>

        <section
          class="results"
          id="results"
          role="region"
          aria-live="polite"
          aria-label="Search results"
        ></section>
      </section>

      <footer class="footer">
        <nav aria-label="Footer navigation">
          <a
            href="https://github.com/0x2E/rss-finder"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="View RSS Finder project on GitHub"
          >
            GitHub Project
          </a>
        </nav>
      </footer>

      <div
        class="success-message"
        id="successMessage"
        role="status"
        aria-live="assertive"
      >
        <span class="success-icon" aria-hidden="true">
          ✓
        </span>
        Copied to clipboard!
      </div>
    </main>
  );
});

function normalizeUrl(url) {
  if (!url?.trim()) {
    throw new Error('URL is required');
  }

  const normalizedUrl = url.trim().startsWith('http') ? url.trim() : `https://${url.trim()}`;
  
  try {
    new URL(normalizedUrl);
    return normalizedUrl;
  } catch {
    throw new Error('Invalid URL format');
  }
}

app.post("/api/find-feeds", async (c) => {
  try {
    const { url } = await c.req.json();

    if (!url) {
      return c.json({ error: "URL is required" }, 400);
    }

    let normalizedUrl;
    try {
      normalizedUrl = normalizeUrl(url);
    } catch (error) {
      return c.json({ error: error.message }, 400);
    }

    const feeds = await find(normalizedUrl, {
      userAgent: "rss-finder/2.0",
    });

    return c.json({ feeds });
  } catch (error) {
    return c.json({ error: "Failed to find RSS feeds" }, 500);
  }
});

export default app;
