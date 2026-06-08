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
	<div class="absolute inset-0 bg-gradient-to-b from-primary-700/30 via-primary-900/10 to-transparent pointer-events-none"></div>
	<div class="absolute inset-0" style="background: radial-gradient(ellipse at 50% 0%, oklch(45% 0.2 280deg / 0.15) 0%, transparent 60%)"></div>
	<div class="container mx-auto px-4 relative">
		<h1 class="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-white via-primary-300 to-secondary-300 bg-clip-text text-transparent leading-tight">
			CubeNet CMS
		</h1>
		<p class="text-lg md:text-xl text-surface-400 max-w-2xl mx-auto mb-10">
			Открытая CMS для Minecraft модовых серверов.<br />
			Управляй сборками, серверами и модами в одном месте.
		</p>
		<div class="flex gap-4 justify-center">
			<a href="/servers" class="btn variant-filled px-8">Серверы</a>
			<a href="/login" class="btn variant-soft px-8">Начать</a>
		</div>
	</div>
</section>

<!-- Фичи -->
<section class="py-16">
	<div class="container mx-auto px-4">
		<h2 class="text-2xl font-bold text-center mb-12">Возможности</h2>
		<div class="grid md:grid-cols-3 gap-6 max-w-5xl mx-auto">
			<div class="card p-6 text-center bg-surface-900/50 border-primary-700/20">
				<div class="w-12 h-12 rounded-full bg-primary-600/20 flex items-center justify-center mx-auto mb-4">
					<span class="text-primary-300 text-xl font-bold">S</span>
				</div>
				<h3 class="font-semibold mb-2">Мониторинг серверов</h3>
				<p class="text-sm text-surface-400">TPS, онлайны, статус — всё в реальном времени через WebSocket.</p>
			</div>
			<div class="card p-6 text-center bg-surface-900/50 border-secondary-600/20">
				<div class="w-12 h-12 rounded-full bg-secondary-600/20 flex items-center justify-center mx-auto mb-4">
					<span class="text-secondary-300 text-xl font-bold">B</span>
				</div>
				<h3 class="font-semibold mb-2">Сборки и моды</h3>
				<p class="text-sm text-surface-400">Раздавай модпаки и конфиги через встроенный лаунчер.</p>
			</div>
			<div class="card p-6 text-center bg-surface-900/50 border-primary-700/20">
				<div class="w-12 h-12 rounded-full bg-primary-600/20 flex items-center justify-center mx-auto mb-4">
					<span class="text-primary-300 text-xl font-bold">O</span>
				</div>
				<h3 class="font-semibold mb-2">Open Source</h3>
				<p class="text-sm text-surface-400">Go + SvelteKit + Avalonia — полностью открытый код.</p>
			</div>
		</div>
	</div>
</section>

<!-- Серверы -->
<section class="py-16">
	<div class="container mx-auto px-4">
		<div class="flex items-center justify-between mb-8">
			<h2 class="text-2xl font-bold">Активные серверы</h2>
			<a href="/servers" class="text-sm text-primary-400 hover:text-primary-300 transition-colors">Все →</a>
		</div>

		{#if servers.length === 0}
			<div class="card p-12 text-center bg-surface-900/50 border-dashed border-surface-700">
				<p class="text-surface-500">Серверы пока не добавлены</p>
			</div>
		{:else}
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each servers as server}
					<a href="/servers/{server.slug}" class="card p-5 bg-surface-900/50 border-surface-700 hover:border-primary-500/50 transition-colors">
						<div class="flex items-center justify-between mb-3">
							<h3 class="font-semibold">{server.name}</h3>
							<span class="badge variant-{server.status === 'online' ? 'success' : 'warning'} text-xs">{server.status}</span>
						</div>
						<div class="text-sm text-surface-400">
							Игроков: {server.players}/{server.max_players}
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</div>
</section>
