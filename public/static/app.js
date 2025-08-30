class RSSFinder {
  constructor() {
    this.debounceTimeout = null;
    this.init();
  }

  init() {
    this.form = document.getElementById("searchForm");
    this.urlInput = document.getElementById("urlInput");
    this.searchBtn = document.getElementById("searchBtn");
    this.buttonText = this.searchBtn.querySelector(".button-text");
    this.results = document.getElementById("results");
    this.successMessage = document.getElementById("successMessage");

    this.bindEvents();
  }

  bindEvents() {
    this.form.addEventListener("submit", (e) => this.handleSubmit(e));
    this.urlInput.addEventListener("input", (e) => this.handleInputChange(e));
    this.urlInput.addEventListener("keydown", (e) => this.handleKeyDown(e));
    document.addEventListener("keydown", (e) => this.handleGlobalKeyDown(e));
  }

  handleInputChange(e) {
    const url = e.target.value.trim();

    // Clear previous timeout
    if (this.debounceTimeout) {
      clearTimeout(this.debounceTimeout);
    }

    // Enable/disable button based on input
    this.searchBtn.disabled = !url;

    // Auto-submit if URL looks valid (debounced)
    if (url && this.isValidUrl(url)) {
      this.debounceTimeout = setTimeout(() => {
        this.performSearch(url);
      }, 800);
    }
  }

  handleKeyDown(e) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      this.handleSubmit(e);
    }
  }

  handleGlobalKeyDown(e) {
    // Allow Escape key to clear results
    if (e.key === "Escape") {
      this.clearResults();
      this.urlInput.focus();
    }

    // Allow Ctrl/Cmd + K to focus search
    if ((e.ctrlKey || e.metaKey) && e.key === "k") {
      e.preventDefault();
      this.urlInput.focus();
      this.urlInput.select();
    }
  }

  normalizeUrl(url) {
    if (!url?.trim()) {
      throw new Error("URL is required");
    }

    const normalizedUrl = url.trim().startsWith("http")
      ? url.trim()
      : `https://${url.trim()}`;

    try {
      new URL(normalizedUrl);
      return normalizedUrl;
    } catch {
      throw new Error("Invalid URL format");
    }
  }

  isValidUrl(string) {
    try {
      this.normalizeUrl(string);
      return true;
    } catch (_) {
      return false;
    }
  }

  clearResults() {
    this.results.innerHTML = "";
  }

  async handleSubmit(e) {
    e.preventDefault();

    const url = this.urlInput.value.trim();
    if (!url) return;

    this.performSearch(url);
  }

  async performSearch(url) {
    // Clear any pending debounced search
    if (this.debounceTimeout) {
      clearTimeout(this.debounceTimeout);
      this.debounceTimeout = null;
    }

    // Normalize URL first
    let normalizedUrl;
    try {
      normalizedUrl = this.normalizeUrl(url);
    } catch (error) {
      this.showError(error.message);
      return;
    }

    this.setLoadingState(true);
    this.showLoadingMessage();

    try {
      const response = await fetch("/api/find-feeds", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: normalizedUrl }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Failed to find feeds");
      }

      this.displayFeeds(data.feeds);
    } catch (error) {
      this.showError(error.message);
    } finally {
      this.setLoadingState(false);
    }
  }

  setLoadingState(loading) {
    this.searchBtn.disabled = loading;
    this.buttonText.textContent = loading ? "Discovering..." : "Find RSS";
  }

  showLoadingMessage() {
    this.results.innerHTML = `
      <div class="loading-skeleton" role="status" aria-live="polite">
        <div class="skeleton-item">
          <div class="skeleton-bar"></div>
          <div class="skeleton-actions">
            <div class="skeleton-button"></div>
            <div class="skeleton-button"></div>
          </div>
        </div>
        <div class="skeleton-item">
          <div class="skeleton-bar"></div>
          <div class="skeleton-actions">
            <div class="skeleton-button"></div>
            <div class="skeleton-button"></div>
          </div>
        </div>
      </div>
    `;
  }

  showError(message) {
    this.results.innerHTML = `<div class="error" role="alert">Error: ${this.escapeHtml(
      message
    )}</div>`;
  }

  displayFeeds(feeds) {
    if (!feeds || feeds.length === 0) {
      this.results.innerHTML =
        '<div class="error" role="alert">No RSS feeds found for this URL.</div>';
      return;
    }

    const feedsHtml = feeds
      .map((feed, index) => {
        const title = feed.title || feed.name || "Untitled Feed";
        const url =
          feed.url ||
          feed.href ||
          feed.link ||
          (typeof feed === "string" ? feed : "");

        return `
        <article class="feed-item" role="article">
          <div class="feed-info">
            <h3 class="feed-title">${this.escapeHtml(title)}</h3>
            <div class="feed-url" aria-label="Feed URL">${this.escapeHtml(
              url
            )}</div>
          </div>
          <div class="feed-actions">
            <button class="btn-small" onclick="rssFinderApp.openInNewTab('${this.escapeHtml(
              url
            )}')" aria-label="Open ${this.escapeHtml(
          title
        )} in new tab">Open</button>
            <button class="btn-small btn-secondary" onclick="rssFinderApp.copyToClipboard('${this.escapeHtml(
              url
            )}')" aria-label="Copy ${this.escapeHtml(title)} URL">Copy</button>
          </div>
        </article>
      `;
      })
      .join("");

    this.results.innerHTML = `
      <div class="feed-list" role="list" aria-label="Found RSS feeds">
        <h2 class="results-title">Found ${feeds.length} RSS feed${
      feeds.length !== 1 ? "s" : ""
    }</h2>
        ${feedsHtml}
      </div>
    `;
  }

  openInNewTab(url) {
    window.open(url, "_blank");
  }

  async copyToClipboard(url) {
    try {
      await navigator.clipboard.writeText(url);
      this.showSuccessMessage();
    } catch (err) {
      console.error("Failed to copy: ", err);
    }
  }

  showSuccessMessage() {
    this.successMessage.classList.add("show");
    setTimeout(() => {
      this.successMessage.classList.remove("show");
    }, 2000);
  }

  escapeHtml(text) {
    if (!text || typeof text !== "string") return String(text || "");
    const map = {
      "&": "&amp;",
      "<": "&lt;",
      ">": "&gt;",
      '"': "&quot;",
      "'": "&#039;",
    };
    return text.replace(/[&<>"']/g, (m) => map[m]);
  }
}

// Initialize when DOM is loaded
document.addEventListener("DOMContentLoaded", () => {
  window.rssFinderApp = new RSSFinder();
});
