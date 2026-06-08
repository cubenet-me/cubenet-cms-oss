<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<div class="flex flex-col min-h-screen" data-theme="catppuccin">
	<header class="h-16 px-4 flex items-center border-b border-surface-800/50 backdrop-blur-sm bg-surface-950/80 sticky top-0 z-50">
		<div class="container mx-auto flex items-center justify-between">
			<a href="/" class="text-xl font-bold bg-gradient-to-r from-primary-300 to-secondary-300 bg-clip-text text-transparent">
				CubeNet
			</a>
			<nav class="flex items-center gap-6">
				<a href="/servers" class="text-sm font-medium text-surface-300 hover:text-primary-300 transition-colors">Servers</a>
				<a href="/admin" class="text-sm font-medium text-surface-300 hover:text-primary-300 transition-colors">Admin</a>
				<a href="/login" class="btn variant-filled btn-sm">Войти</a>
			</nav>
		</div>
	</header>

	<main class="flex-1">
		{@render children()}
	</main>

	<footer class="py-8 border-t border-surface-800/50">
		<div class="container mx-auto px-4 text-center text-sm text-surface-500">
			<p class="mb-1">Сделано на <a href="https://github.com/cubenet-me/cubenet-cms-oss" class="text-primary-400 hover:text-primary-300">CubeNet CMS OSS</a></p>
			<p class="text-xs text-surface-600">CubeNet CMS License v1.0 — Использование в коммерческих целях без сохранения копирайта запрещено</p>
		</div>
	</footer>
</div>
