// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'Simple-Sync Documentation',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/kwila/simple-sync' }],
			sidebar: [
				{ label: 'Overview', slug: 'index' },
				{
					label: 'API',
					autogenerate: { directory: 'api' },
				},
				{ label: 'ACL', slug: 'acl' },
				{ label: 'Tech Stack', slug: 'tech-stack' },
			],
		}),
	],
});
