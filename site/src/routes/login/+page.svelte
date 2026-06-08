<div class="container mx-auto px-4 py-16">
	<div class="max-w-sm mx-auto">
		<div class="card p-8 bg-surface-900/50 border-surface-700">
			<h1 class="text-2xl font-bold text-center mb-2">Вход</h1>
			<p class="text-sm text-surface-500 text-center mb-8">Войдите в панель управления</p>

			<form method="POST" action="/api/v1/auth/login" class="space-y-4">
				<div class="space-y-1">
					<label class="label text-sm" for="username">Имя пользователя</label>
					<input class="input" type="text" id="username" name="username" required />
				</div>
				<div class="space-y-1">
					<label class="label text-sm" for="password">Пароль</label>
					<input class="input" type="password" id="password" name="password" required />
				</div>
				<button type="submit" class="btn variant-filled w-full">Войти</button>
			</form>

			<p class="text-sm text-center mt-6 text-surface-400">
				Нет аккаунта? <a href="/register" class="text-primary-400 hover:text-primary-300">Зарегистрироваться</a>
			</p>
		</div>
	</div>
</div>
