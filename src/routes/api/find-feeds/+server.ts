import { json } from '@sveltejs/kit';
import { find } from 'feedfinder-ts';
import { normalizeUrl } from '$lib/utils.js';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request }) => {
	try {
		const { url } = await request.json();

		if (!url) {
			return json({ error: 'URL is required' }, { status: 400 });
		}

		let normalizedUrl: string;
		try {
			normalizedUrl = normalizeUrl(url);
		} catch (error) {
			return json({ error: (error as Error).message }, { status: 400 });
		}

		const feeds = await find(normalizedUrl, {
			userAgent: 'rss-finder/2.0'
		});

		return json({ feeds });
	} catch (error) {
		console.error('Error finding RSS feeds:', error);
		return json({ error: 'Failed to find RSS feeds' }, { status: 500 });
	}
};
