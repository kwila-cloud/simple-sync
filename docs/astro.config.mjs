// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'Simple-Sync Docs',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/kwila-cloud/simple-sync' }],
			sidebar: [
				{ label: 'Overview 🏠', slug: 'overview' },
				{ label: 'Tech Stack ⚙️', slug: 'tech-stack' },
				{
					label: 'API 🚀',
					autogenerate: { directory: 'api' },
				},
				{ label: 'ACL 🛡️', slug: 'acl' },
				{ label: 'Download as PDF 📄', link: 'https://kwila.github.io/simple-sync/docs.pdf', attrs: { target: '_blank' } },
			],
		}),
	],
});
