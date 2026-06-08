<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<div class="flex flex-col min-h-screen">
	<header class="h-16 px-4 flex items-center border-b border-surface-200/20 backdrop-blur-sm bg-surface-950/80 sticky top-0 z-50">
		<div class="container mx-auto flex items-center justify-between">
			<a href="/" class="text-xl font-bold bg-gradient-to-r from-primary-400 to-secondary-400 bg-clip-text text-transparent">
				CubeNet
			</a>
			<nav class="flex items-center gap-6">
				<a href="/servers" class="text-sm font-medium text-surface-300 hover:text-white transition-colors">Servers</a>
				<a href="/admin" class="text-sm font-medium text-surface-300 hover:text-white transition-colors">Admin</a>
				<a href="/login" class="btn variant-filled btn-sm">Sign in</a>
			</nav>
		</div>
	</header>

	<main class="flex-1">
		{@render children()}
	</main>

	<footer class="py-8 border-t border-surface-200/10">
		<div class="container mx-auto px-4 text-center text-sm text-surface-500">
			<p>CubeNet CMS — Open Source</p>
		</div>
	</footer>
</div>
