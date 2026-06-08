<div class="container mx-auto max-w-sm p-8">
	<h1 class="text-2xl font-bold mb-6 text-center">Login</h1>

	<form method="POST" action="/api/v1/auth/login" class="space-y-4">
		<div>
			<label class="label" for="username">Username</label>
			<input class="input" type="text" id="username" name="username" required />
		</div>
		<div>
			<label class="label" for="password">Password</label>
			<input class="input" type="password" id="password" name="password" required />
		</div>
		<button type="submit" class="btn variant-filled w-full">Sign in</button>
	</form>
</div>
