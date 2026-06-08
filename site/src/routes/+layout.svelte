<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<nav class="navbar px-4 py-2">
	<div class="flex items-center justify-between">
		<a href="/" class="text-xl font-bold">CubeNet</a>
		<div class="flex gap-4">
			<a href="/servers" class="nav-link">Servers</a>
			<a href="/admin" class="nav-link">Admin</a>
			<a href="/login" class="btn variant-filled btn-sm">Sign in</a>
		</div>
	</div>
</nav>

<main>
	{@render children()}
</main>
