import adapterAuto from '@sveltejs/adapter-auto';
import adapterNode from '@sveltejs/adapter-node';

// Determine which adapter to use based on deployment type environment variable
const deploymentType = process.env.DEPLOYMENT_TYPE || 'auto';
const adapter = deploymentType === 'node' ? adapterNode : adapterAuto;

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter()
	}
};

export default config;
