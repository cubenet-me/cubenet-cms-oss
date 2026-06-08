<script lang="ts">
	interface Server {
		id: string; name: string; slug: string;
		status: string; players: number; max_players: number;
	}

	let servers = $state<Server[]>([]);

	async function load() {
		const res = await fetch('/api/v1/servers');
		if (res.ok) servers = await res.json();
	}

	$effect(() => { load(); });
</script>

<!-- Hero -->
<section class="py-24 text-center relative overflow-hidden">
	<div class="absolute inset-0 bg-gradient-to-b from-primary-500/10 via-secondary-500/5 to-transparent pointer-events-none"></div>
	<div class="container mx-auto px-4 relative">
		<h1 class="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-white via-primary-300 to-secondary-300 bg-clip-text text-transparent">
			CubeNet CMS
		</h1>
		<p class="text-lg md:text-xl text-surface-400 max-w-2xl mx-auto mb-10">
			Open source content management system for Minecraft modded servers.<br />
			Manage builds, servers, and mods — all in one place.
		</p>
		<div class="flex gap-4 justify-center">
			<a href="/servers" class="btn variant-filled px-8">Browse Servers</a>
			<a href="/login" class="btn variant-soft px-8">Get Started</a>
		</div>
	</div>
</section>

<!-- Features -->
<section class="py-16">
	<div class="container mx-auto px-4">
		<h2 class="text-2xl font-bold text-center mb-12">Why CubeNet?</h2>
		<div class="grid md:grid-cols-3 gap-6 max-w-5xl mx-auto">
			<div class="card p-6 text-center">
				<div class="w-12 h-12 rounded-full bg-primary-500/20 flex items-center justify-center mx-auto mb-4">
					<span class="text-primary-400 text-xl font-bold">S</span>
				</div>
				<h3 class="font-semibold mb-2">Server Management</h3>
				<p class="text-sm text-surface-400">Real-time monitoring, TPS tracking, player counts — all at a glance.</p>
			</div>
			<div class="card p-6 text-center">
				<div class="w-12 h-12 rounded-full bg-secondary-500/20 flex items-center justify-center mx-auto mb-4">
					<span class="text-secondary-400 text-xl font-bold">B</span>
				</div>
				<h3 class="font-semibold mb-2">Build Distribution</h3>
				<p class="text-sm text-surface-400">Distribute modpacks and configs to players via integrated launcher.</p>
			</div>
			<div class="card p-6 text-center">
				<div class="w-12 h-12 rounded-full bg-warning-500/20 flex items-center justify-center mx-auto mb-4">
					<span class="text-warning-400 text-xl font-bold">O</span>
				</div>
				<h3 class="font-semibold mb-2">Open Source</h3>
				<p class="text-sm text-surface-400">Go backend, SvelteKit frontend, Avalonia launcher — fully open source.</p>
			</div>
		</div>
	</div>
</section>

<!-- Servers Preview -->
<section class="py-16">
	<div class="container mx-auto px-4">
		<div class="flex items-center justify-between mb-8">
			<h2 class="text-2xl font-bold">Servers</h2>
			<a href="/servers" class="text-sm text-primary-400 hover:text-primary-300 transition-colors">View all →</a>
		</div>

		{#if servers.length === 0}
			<div class="card p-12 text-center">
				<p class="text-surface-500">No servers configured yet.</p>
			</div>
		{:else}
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each servers as server}
					<a href="/servers/{server.slug}" class="card p-5 hover:border-primary-500/50 transition-colors">
						<div class="flex items-center justify-between mb-3">
							<h3 class="font-semibold">{server.name}</h3>
							<span class="badge variant-{server.status === 'online' ? 'success' : 'warning'} text-xs">{server.status}</span>
						</div>
						<div class="text-sm text-surface-400">
							{server.players}/{server.max_players} players
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</div>
</section>
