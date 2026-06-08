<script lang="ts">
	interface Server {
		id: string; name: string; slug: string; description: string;
		address: string; version: string; status: string;
		tps: number; players: number; max_players: number; mods: string[];
	}

	let servers = $state<Server[]>([]);
	let loading = $state(true);

	async function load() {
		const res = await fetch('/api/v1/servers');
		if (res.ok) servers = await res.json();
		loading = false;
	}

	$effect(() => { load(); });
</script>

<div class="container mx-auto px-4 py-12">
	<h1 class="text-3xl font-bold mb-8">Servers</h1>

	{#if loading}
		<div class="card p-12 text-center">
			<p class="text-surface-500">Loading...</p>
		</div>
	{:else if servers.length === 0}
		<div class="card p-12 text-center">
			<p class="text-surface-500">No servers found.</p>
		</div>
	{:else}
		<div class="grid gap-4 max-w-3xl">
			{#each servers as server}
				<div class="card p-5">
					<div class="flex items-start justify-between mb-3">
						<div>
							<h2 class="text-xl font-semibold">{server.name}</h2>
							<p class="text-sm text-surface-500">{server.slug}</p>
						</div>
						<div class="flex items-center gap-3">
							<span class="badge variant-{server.status === 'online' ? 'success' : 'warning'}">{server.status}</span>
						</div>
					</div>
					<div class="flex flex-wrap gap-4 text-sm text-surface-400">
						<span>Players: {server.players}/{server.max_players}</span>
						<span>TPS: {server.tps.toFixed(1)}</span>
						<span>Version: {server.version || '—'}</span>
						<span>Address: {server.address || '—'}</span>
					</div>
					{#if server.mods.length > 0}
						<div class="flex flex-wrap gap-2 mt-3">
							{#each server.mods as mod}
								<span class="badge variant-soft-surface text-xs">{mod}</span>
							{/each}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
