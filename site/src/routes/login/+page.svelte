<div class="container mx-auto px-4 py-16">
	<div class="max-w-sm mx-auto">
		<div class="card p-8">
			<h1 class="text-2xl font-bold text-center mb-8">Sign in</h1>

			<form method="POST" action="/api/v1/auth/login" class="space-y-4">
				<div class="space-y-1">
					<label class="label text-sm" for="username">Username</label>
					<input class="input" type="text" id="username" name="username" required />
				</div>
				<div class="space-y-1">
					<label class="label text-sm" for="password">Password</label>
					<input class="input" type="password" id="password" name="password" required />
				</div>
				<button type="submit" class="btn variant-filled w-full">Sign in</button>
			</form>
		</div>
	</div>
</div>
