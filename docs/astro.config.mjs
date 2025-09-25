// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	base: '/simple-sync',
	site: 'https://kwila-cloud.github.io/simple-sync',
	integrations: [
		starlight({
			title: 'Simple Sync Docs',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/kwila-cloud/simple-sync' }],
			sidebar: [
				{ label: 'Overview 🏠', slug: 'overview' },
				{ label: 'Tech Stack ⚙️', slug: 'tech-stack' },
				{
					label: 'API 🚀',
					autogenerate: { directory: 'api' },
				},
				{ label: 'ACL 🛡️', slug: 'acl' },
				{ label: 'Internal Events 📊', slug: 'internal-events' },
				{ label: 'Release History 📋', link: 'https://github.com/kwila-cloud/simple-sync/blob/main/CHANGELOG.md', attrs: { target: '_blank' } },
				{ label: 'Download as PDF 📄', link: '/docs.pdf', attrs: { target: '_blank' } },
			],
		}),
	],
});
