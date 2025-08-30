<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta name="description" content="Discover RSS feeds from any website quickly and easily." />
	<link rel="icon" type="image/x-icon" href="/favicon.ico" />
	<title>RSS Finder - Discover RSS Feeds from Any Website</title>
</svelte:head>

{@render children?.()}
