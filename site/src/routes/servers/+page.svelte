<script lang="ts">
	let servers = $state<Array<{ id: string; name: string; slug: string; status: string; players: number; max_players: number }>>([]);

	async function load() {
		const res = await fetch('http://localhost:8080/api/v1/servers');
		if (res.ok) servers = await res.json();
	}

	$effect(() => { load(); });
</script>

<div class="container mx-auto p-8">
	<h1 class="text-3xl font-bold mb-6">Servers</h1>

	{#if servers.length === 0}
		<p class="text-surface-500">No servers found.</p>
	{:else}
		<div class="grid gap-4">
			{#each servers as server}
				<div class="card p-4">
					<h2 class="text-xl font-semibold">{server.name}</h2>
					<p class="text-sm text-surface-500">{server.slug}</p>
					<div class="flex gap-2 mt-2">
						<span class="badge variant-{server.status === 'online' ? 'success' : 'warning'}">{server.status}</span>
						<span class="text-sm">{server.players}/{server.max_players} players</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
