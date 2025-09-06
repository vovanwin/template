<script lang="ts">
	import { createQuery, createMutation } from '@tanstack/svelte-query';
	import { goto } from '$app/navigation';
	import { getAuthMe, postAuthLogout } from '$lib/api/default/default';
	import { authStore, clearAuth } from '$lib/stores/auth';
	import { onMount } from 'svelte';

	// Проверяем авторизацию при загрузке
	onMount(() => {
		const unsubscribe = authStore.subscribe((auth) => {
			if (auth.initialized && !auth.isAuthenticated) {
				goto('/login');
			}
		});

		return unsubscribe;
	});

	// Получаем данные пользователя
	const userQuery = createQuery({
		queryKey: ['user', 'me'],
		queryFn: () => getAuthMe(),
		enabled: $authStore.isAuthenticated,
		onError: () => {
			// Если не удалось получить данные пользователя, перенаправляем на логин
			clearAuth();
			goto('/login');
		}
	});

	// Мутация для выхода
	const logoutMutation = createMutation({
		mutationFn: () => postAuthLogout(),
		onSettled: () => {
			// Независимо от результата, очищаем локальную авторизацию
			clearAuth();
			goto('/login');
		}
	});

	const handleLogout = () => {
		$logoutMutation.mutate();
	};

	$: user = $userQuery.data?.data;
</script>

<svelte:head>
	<title>Профиль пользователя</title>
</svelte:head>

<div class="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
	{#if $userQuery.isLoading}
		<div class="flex justify-center items-center h-64">
			<svg class="animate-spin h-8 w-8 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
		</div>
	{:else if $userQuery.isError}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
			Ошибка загрузки данных пользователя
		</div>
	{:else if user}
		<div class="space-y-6">
			<!-- Заголовок -->
			<div class="flex justify-between items-center">
				<h1 class="text-3xl font-bold text-gray-900">Профиль пользователя</h1>
				<button
					on:click={handleLogout}
					class="btn btn-secondary"
					disabled={$logoutMutation.isPending}
				>
					{#if $logoutMutation.isPending}
						Выход...
					{:else}
						Выйти
					{/if}
				</button>
			</div>

			<!-- Карточка с информацией о пользователе -->
			<div class="card p-6">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<div>
						<h3 class="text-lg font-medium text-gray-900 mb-4">Основная информация</h3>
						<dl class="space-y-3">
							<div>
								<dt class="text-sm font-medium text-gray-500">ID пользователя</dt>
								<dd class="text-sm text-gray-900 font-mono">{user.id}</dd>
							</div>
							<div>
								<dt class="text-sm font-medium text-gray-500">Email</dt>
								<dd class="text-sm text-gray-900">{user.email}</dd>
							</div>
							{#if user.role}
								<div>
									<dt class="text-sm font-medium text-gray-500">Роль</dt>
									<dd class="text-sm text-gray-900">
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800">
											{user.role.value || user.role}
										</span>
									</dd>
								</div>
							{/if}
							{#if user.tenant}
								<div>
									<dt class="text-sm font-medium text-gray-500">Тенант</dt>
									<dd class="text-sm text-gray-900">{user.tenant}</dd>
								</div>
							{/if}
							<div>
								<dt class="text-sm font-medium text-gray-500">Дата создания</dt>
								<dd class="text-sm text-gray-900">
									{new Date(user.created_at).toLocaleDateString('ru-RU', {
										year: 'numeric',
										month: 'long',
										day: 'numeric',
										hour: '2-digit',
										minute: '2-digit'
									})}
								</dd>
							</div>
						</dl>
					</div>

					<div>
						<h3 class="text-lg font-medium text-gray-900 mb-4">Дополнительно</h3>
						<dl class="space-y-3">
							{#if user.components && user.components.length > 0}
								<div>
									<dt class="text-sm font-medium text-gray-500">Доступные компоненты</dt>
									<dd class="text-sm text-gray-900">
										<div class="flex flex-wrap gap-2 mt-1">
											{#each user.components as component}
												<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
													{component}
												</span>
											{/each}
										</div>
									</dd>
								</div>
							{/if}
							{#if user.settings}
								<div>
									<dt class="text-sm font-medium text-gray-500">Настройки</dt>
									<dd class="text-sm text-gray-900 font-mono text-xs bg-gray-50 p-2 rounded">
										{user.settings}
									</dd>
								</div>
							{/if}
						</dl>
					</div>
				</div>
			</div>

			<!-- Действия -->
			<div class="card p-6">
				<h3 class="text-lg font-medium text-gray-900 mb-4">Действия</h3>
				<div class="flex space-x-4">
					<button class="btn btn-primary">
						Редактировать профиль
					</button>
					<button class="btn btn-secondary">
						Изменить пароль
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>
