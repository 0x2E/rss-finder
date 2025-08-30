export function normalizeUrl(url: string): string {
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
