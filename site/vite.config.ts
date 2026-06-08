import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': 'http://localhost:8080',
			'/api/v1/ws': {
				target: 'ws://localhost:8080',
				ws: true
			},
			'/health': 'http://localhost:8080'
		}
	}
});
