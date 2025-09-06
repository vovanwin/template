<script lang="ts">
	import { createMutation } from '@tanstack/svelte-query';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { postAuthLogin } from '$lib/api/default/default';
	import { setAuth, setLoading, authStore } from '$lib/stores/auth';
	import type { LoginRequest } from '$lib/api/models';

	let email = '';
	let password = '';
	let errorMessage = '';

	const loginMutation = createMutation({
		mutationFn: (data: LoginRequest) => postAuthLogin(data),
		onMutate: () => {
			setLoading(true);
			errorMessage = '';
		},
		onSuccess: (data) => {
			const token = data.data.token;
			const user = {
				id: data.data.user_id,
				email: data.data.user_email,
				tenant: '',
				created_at: new Date().toISOString(),
				settings: '',
				components: []
			};
			setAuth(token, user);
			goto('/');
		},
		onError: (error: any) => {
			setLoading(false);
			errorMessage = error.response?.data?.message || 'Ошибка входа';
		}
	});

	const handleSubmit = () => {
		if (!email || !password) {
			errorMessage = 'Пожалуйста, заполните все поля';
			return;
		}
		$loginMutation.mutate({ email, password });
	};

	// Проверяем авторизацию при загрузке
	onMount(() => {
		const unsubscribe = authStore.subscribe((auth) => {
			if (auth.initialized && auth.isAuthenticated) {
				goto('/');
			}
		});

		return unsubscribe;
	});

	$: isLoading = $authStore.loading || $loginMutation.isPending;
</script>

<svelte:head>
	<title>Вход в систему</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
				Войти в аккаунт
			</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Введите свои учетные данные
			</p>
		</div>
		
		<form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
			<div class="rounded-md shadow-sm -space-y-px">
				<div>
					<label for="email" class="sr-only">Email</label>
					<input
						id="email"
						name="email"
						type="email"
						required
						class="form-input rounded-t-md"
						placeholder="Email адрес"
						bind:value={email}
						disabled={isLoading}
					/>
				</div>
				<div>
					<label for="password" class="sr-only">Пароль</label>
					<input
						id="password"
						name="password"
						type="password"
						required
						class="form-input rounded-b-md"
						placeholder="Пароль"
						bind:value={password}
						disabled={isLoading}
					/>
				</div>
			</div>

			{#if errorMessage}
				<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
					{errorMessage}
				</div>
			{/if}

			<div>
				<button
					type="submit"
					class="btn btn-primary w-full flex justify-center py-2 px-4 text-sm"
					disabled={isLoading}
				>
					{#if isLoading}
						<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Вход...
					{:else}
						Войти
					{/if}
				</button>
			</div>

			<div class="text-center">
				<p class="text-sm text-gray-600">
					Тестовые данные:
					<br>
					Email: admin@example.com
					<br>
					Пароль: password
				</p>
			</div>
		</form>
	</div>
</div>
