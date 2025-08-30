<script lang="ts">
	import { onMount } from 'svelte';
	import { normalizeUrl } from '$lib/utils.js';

	let urlInput = $state('');
	let isLoading = $state(false);
	let feeds = $state<unknown[]>([]);
	let error = $state('');
	let successMessage = $state(false);
	let hasSearched = $state(false);

	function clearResults() {
		feeds = [];
		error = '';
		hasSearched = false;
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();

		const url = urlInput.trim();
		if (!url) return;

		await performSearch(url);
	}

	async function performSearch(url: string) {
		let normalizedUrl: string;
		try {
			normalizedUrl = normalizeUrl(url);
		} catch (err) {
			error = (err as Error).message;
			feeds = [];
			hasSearched = true;
			return;
		}

		isLoading = true;
		error = '';
		feeds = [];
		hasSearched = true;

		try {
			const response = await fetch('/api/find-feeds', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ url: normalizedUrl })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Failed to find feeds');
			}

			feeds = data.feeds || [];
		} catch (err) {
			error = (err as Error).message;
		} finally {
			isLoading = false;
		}
	}

	function handleInputChange() {
		// Clear previous results when input changes
		feeds = [];
		error = '';
		hasSearched = false;
	}

	function openInNewTab(url: string) {
		window.open(url, '_blank');
	}

	async function copyToClipboard(url: string) {
		try {
			await navigator.clipboard.writeText(url);
			showSuccessMessage();
		} catch (err) {
			console.error('Failed to copy: ', err);
		}
	}

	function showSuccessMessage() {
		successMessage = true;
		setTimeout(() => {
			successMessage = false;
		}, 2000);
	}

	onMount(() => {
		function handleGlobalKeyDown(e: KeyboardEvent) {
			// Allow Escape key to clear results
			if (e.key === 'Escape') {
				clearResults();
				const input = document.getElementById('urlInput') as HTMLInputElement;
				if (input) input.focus();
			}

			// Allow Ctrl/Cmd + K to focus search
			if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
				e.preventDefault();
				const input = document.getElementById('urlInput') as HTMLInputElement;
				if (input) {
					input.focus();
					input.select();
				}
			}
		}

		document.addEventListener('keydown', handleGlobalKeyDown);

		return () => {
			document.removeEventListener('keydown', handleGlobalKeyDown);
		};
	});
</script>

<main class="mx-auto flex min-h-screen max-w-4xl flex-col px-8">
	<section class="flex min-h-[70vh] flex-1 flex-col items-center justify-center">
		<header class="mb-12 text-center">
			<h1 class="gradient-text mb-4 text-6xl font-bold tracking-tight">RSS Finder</h1>
			<p class="text-xl font-normal text-muted-foreground">Discover RSS feeds from any website</p>
		</header>

		<form
			class="mb-16 w-full max-w-2xl"
			id="searchForm"
			role="search"
			aria-label="RSS feed search"
			onsubmit={handleSubmit}
		>
			<div class="flex items-stretch gap-4 rounded-2xl border-2 border-border bg-white p-2">
				<label for="urlInput" class="sr-only">Website URL</label>
				<input
					id="urlInput"
					name="url"
					placeholder="Enter a website URL..."
					aria-describedby="url-help"
					required
					autocomplete="url"
					class="flex-1 rounded-xl border-0 bg-transparent px-6 py-5 text-lg transition-all duration-300 ease-out outline-none placeholder:font-normal placeholder:text-muted-foreground focus:bg-accent"
					bind:value={urlInput}
					oninput={handleInputChange}
				/>
				<button
					type="submit"
					id="searchBtn"
					aria-describedby="search-help"
					disabled={!urlInput.trim() || isLoading}
					class="button-shimmer relative cursor-pointer overflow-hidden rounded-xl border-0 bg-gradient-to-br from-primary to-primary-hover px-10 py-5 text-lg font-semibold text-primary-foreground transition-all duration-200 ease-out hover:translate-y-[-2px] focus-visible:outline focus-visible:outline-3 focus-visible:outline-offset-2 focus-visible:outline-ring/40 active:translate-y-[-1px] disabled:transform-none disabled:cursor-not-allowed disabled:bg-muted-foreground disabled:opacity-60"
				>
					<span class="button-text">{isLoading ? 'Discovering...' : 'Find RSS'}</span>
				</button>
			</div>
			<small id="url-help" class="mt-3 block text-center text-sm text-muted-foreground">
				Enter any website URL to discover its RSS feeds
			</small>
		</form>

		<section class="w-full max-w-3xl" id="results" aria-live="polite" aria-label="Search results">
			{#if isLoading}
				<div class="flex flex-col gap-4" role="status" aria-live="polite">
					<div class="flex items-center rounded-2xl border border-border bg-white p-7">
						<div
							class="mr-4 h-5 flex-1 animate-pulse rounded bg-gradient-to-r from-border via-muted to-border bg-[length:200%_100%]"
						></div>
						<div class="flex gap-2">
							<div
								class="h-8 w-16 animate-pulse rounded-lg bg-gradient-to-r from-border via-muted to-border bg-[length:200%_100%]"
							></div>
							<div
								class="h-8 w-16 animate-pulse rounded-lg bg-gradient-to-r from-border via-muted to-border bg-[length:200%_100%]"
							></div>
						</div>
					</div>
					<div class="flex items-center rounded-2xl border border-border bg-white p-7">
						<div
							class="mr-4 h-5 flex-1 animate-pulse rounded bg-gradient-to-r from-border via-muted to-border bg-[length:200%_100%]"
						></div>
						<div class="flex gap-2">
							<div
								class="h-8 w-16 animate-pulse rounded-lg bg-gradient-to-r from-border via-muted to-border bg-[length:200%_100%]"
							></div>
							<div
								class="h-8 w-16 animate-pulse rounded-lg bg-gradient-to-r from-border via-muted to-border bg-[length:200%_100%]"
							></div>
						</div>
					</div>
				</div>
			{:else if error}
				<div
					class="mb-4 rounded-xl border border-destructive/20 bg-destructive/5 px-5 py-4 text-center text-sm text-destructive"
					role="alert"
				>
					Error: {error}
				</div>
			{:else if feeds.length === 0 && hasSearched && urlInput.trim()}
				<div
					class="mb-4 rounded-xl border border-destructive/20 bg-destructive/5 px-5 py-4 text-center text-sm text-destructive"
					role="alert"
				>
					No RSS feeds found for this URL.
				</div>
			{:else if feeds.length > 0}
				<div class="flex flex-col gap-4" role="list" aria-label="Found RSS feeds">
					<h2 class="text-foreground mb-6 text-center text-xl font-medium tracking-tight">
						Found {feeds.length} RSS feed{feeds.length !== 1 ? 's' : ''}
					</h2>
					{#each feeds as feed, index (index)}
						{@const feedObj = feed as {
							title?: string;
							name?: string;
							url?: string;
							href?: string;
							link?: string;
						}}
						{@const title = feedObj.title || feedObj.name || 'Untitled Feed'}
						{@const url =
							feedObj.url || feedObj.href || feedObj.link || (typeof feed === 'string' ? feed : '')}
						<article
							class="flex items-center rounded-3xl border border-border bg-gradient-to-br from-white to-muted p-8 transition-all duration-300 ease-out hover:translate-y-[-4px] hover:border-primary hover:bg-gradient-to-br hover:from-white hover:to-accent"
						>
							<div class="mr-4 flex-1">
								<h3 class="text-foreground mb-2 text-lg font-medium tracking-tight">
									{title}
								</h3>
								<div
									class="font-mono text-sm break-all text-muted-foreground"
									aria-label="Feed URL"
								>
									{url}
								</div>
							</div>
							<div class="flex gap-2">
								<button
									class="button-shimmer relative min-w-0 cursor-pointer overflow-hidden rounded-lg border-0 bg-gradient-to-br from-primary to-primary-hover px-5 py-2.5 text-sm font-semibold text-primary-foreground transition-all duration-200 ease-out hover:translate-y-[-2px] active:translate-y-[-1px] disabled:transform-none disabled:cursor-not-allowed disabled:bg-muted-foreground disabled:opacity-60"
									onclick={() => openInNewTab(url)}
									aria-label="Open {title} in new tab"
								>
									Open
								</button>
								<button
									class="button-shimmer relative min-w-0 cursor-pointer overflow-hidden rounded-lg border-0 bg-gradient-to-br from-secondary to-secondary-hover px-5 py-2.5 text-sm font-semibold text-secondary-foreground transition-all duration-200 ease-out hover:translate-y-[-2px] active:translate-y-[-1px] disabled:transform-none disabled:cursor-not-allowed disabled:bg-muted-foreground disabled:opacity-60"
									onclick={() => copyToClipboard(url)}
									aria-label="Copy {title} URL"
								>
									Copy
								</button>
							</div>
						</article>
					{/each}
				</div>
			{/if}
		</section>
	</section>

	<footer class="mt-16 py-8 pt-4 text-center">
		<nav aria-label="Footer navigation">
			<a
				href="https://github.com/0x2E/rss-finder"
				target="_blank"
				rel="noopener noreferrer"
				aria-label="View RSS Finder project on GitHub"
				class="rounded-lg px-4 py-2 text-sm font-normal text-muted-foreground no-underline transition-all duration-200 ease-out hover:bg-accent hover:text-primary"
			>
				GitHub Project
			</a>
		</nav>
	</footer>

	<div
		class="success-message bg-foreground text-background fixed top-8 right-8 z-50 rounded-lg px-5 py-3 text-sm"
		class:show={successMessage}
		id="successMessage"
		role="status"
		aria-live="assertive"
	>
		<span class="success-icon mr-2 font-bold" aria-hidden="true">✓</span>
		Copied to clipboard!
	</div>
</main>
