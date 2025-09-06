<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';

	// Проверяем авторизацию при загрузке
	onMount(() => {
		// Ждем инициализации
		const unsubscribe = authStore.subscribe((auth) => {
			if (auth.initialized && !auth.isAuthenticated) {
				goto('/login');
			}
		});

		return unsubscribe;
	});

	// Функции для управления приложением
	const openDocumentation = () => {
		window.open('http://localhost:8081/docs', '_blank');
	};

	const openProfiling = () => {
		window.open('http://localhost:8082/debug/pprof', '_blank');
	};

	const openMetrics = () => {
		window.open('http://localhost:8082/metrics', '_blank');
	};

	const openSwagger = () => {
		window.open('http://localhost:8083', '_blank');
	};
</script>

<svelte:head>
	<title>Админ панель</title>
	<meta name="description" content="Административная панель управления приложением" />
</svelte:head>

{#if $authStore.loading || !$authStore.initialized}
	<!-- Экран загрузки -->
	<div class="min-h-screen flex items-center justify-center">
		<div class="text-center">
			<svg class="animate-spin h-12 w-12 text-primary-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			<h2 class="text-lg font-medium text-gray-900">Загрузка...</h2>
			<p class="text-sm text-gray-500">Проверка авторизации</p>
		</div>
	</div>
{:else if $authStore.isAuthenticated}
	<div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
		<!-- Заголовок -->
		<div class="border-b border-gray-200 pb-5 mb-6">
			<h1 class="text-3xl font-bold leading-tight text-gray-900">
				Административная панель
			</h1>
			<p class="mt-2 text-sm text-gray-600">
				Управление приложением и мониторинг системы
			</p>
		</div>

		<!-- Быстрая статистика -->
		<div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8">
			<div class="card p-5">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<div class="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center">
							<svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
							</svg>
						</div>
					</div>
					<div class="ml-5 w-0 flex-1">
						<dl>
							<dt class="text-sm font-medium text-gray-500 truncate">Статус системы</dt>
							<dd class="text-lg font-medium text-gray-900">Активна</dd>
						</dl>
					</div>
				</div>
			</div>

			<div class="card p-5">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<div class="w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center">
							<svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
								<path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3z"/>
							</svg>
						</div>
					</div>
					<div class="ml-5 w-0 flex-1">
						<dl>
							<dt class="text-sm font-medium text-gray-500 truncate">Активные пользователи</dt>
							<dd class="text-lg font-medium text-gray-900">1</dd>
						</dl>
					</div>
				</div>
			</div>

			<div class="card p-5">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<div class="w-8 h-8 bg-purple-500 rounded-md flex items-center justify-center">
							<svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
							</svg>
						</div>
					</div>
					<div class="ml-5 w-0 flex-1">
						<dl>
							<dt class="text-sm font-medium text-gray-500 truncate">API запросы</dt>
							<dd class="text-lg font-medium text-gray-900">--</dd>
						</dl>
					</div>
				</div>
			</div>

			<div class="card p-5">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<div class="w-8 h-8 bg-yellow-500 rounded-md flex items-center justify-center">
							<svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
							</svg>
						</div>
					</div>
					<div class="ml-5 w-0 flex-1">
						<dl>
							<dt class="text-sm font-medium text-gray-500 truncate">Uptime</dt>
							<dd class="text-lg font-medium text-gray-900">--</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<!-- Основные действия -->
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2 mb-8">
			<!-- Управление системой -->
			<div class="card p-6">
				<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
					Управление системой
				</h3>
				<div class="space-y-3">
					<button
						on:click={openDocumentation}
						class="w-full flex items-center justify-between p-3 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
					>
						<div class="flex items-center">
							<svg class="w-5 h-5 text-blue-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clip-rule="evenodd" />
							</svg>
							<span class="text-sm font-medium text-gray-900">Документация API</span>
						</div>
						<svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
						</svg>
					</button>

					<button
						on:click={openSwagger}
						class="w-full flex items-center justify-between p-3 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
					>
						<div class="flex items-center">
							<svg class="w-5 h-5 text-green-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M12.316 3.051a1 1 0 01.633 1.265l-4 12a1 1 0 11-1.898-.632l4-12a1 1 0 011.265-.633zM5.707 6.293a1 1 0 010 1.414L3.414 10l2.293 2.293a1 1 0 11-1.414 1.414l-3-3a1 1 0 010-1.414l3-3a1 1 0 011.414 0zm8.586 0a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 11-1.414-1.414L16.586 10l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
							</svg>
							<span class="text-sm font-medium text-gray-900">Swagger UI</span>
						</div>
						<svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
						</svg>
					</button>
				</div>
			</div>

			<!-- Мониторинг -->
			<div class="card p-6">
				<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
					Мониторинг и отладка
				</h3>
				<div class="space-y-3">
					<button
						on:click={openProfiling}
						class="w-full flex items-center justify-between p-3 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
					>
						<div class="flex items-center">
							<svg class="w-5 h-5 text-red-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M3 3a1 1 0 000 2v8a2 2 0 002 2h2.586l-1.293 1.293a1 1 0 101.414 1.414L10 15.414l2.293 2.293a1 1 0 001.414-1.414L12.414 15H15a2 2 0 002-2V5a1 1 0 100-2H3zm11.707 4.707a1 1 0 00-1.414-1.414L10 9.586 8.707 8.293a1 1 0 00-1.414 0l-2 2a1 1 0 101.414 1.414L8 10.414l1.293 1.293a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
							</svg>
							<span class="text-sm font-medium text-gray-900">Профилирование (pprof)</span>
						</div>
						<svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
						</svg>
					</button>

					<button
						on:click={openMetrics}
						class="w-full flex items-center justify-between p-3 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
					>
						<div class="flex items-center">
							<svg class="w-5 h-5 text-purple-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
								<path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z"/>
							</svg>
							<span class="text-sm font-medium text-gray-900">Метрики системы</span>
						</div>
						<svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
						</svg>
					</button>
				</div>
			</div>
		</div>

		<!-- Логи и управление -->
		<div class="card p-6">
			<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
				Управление логированием
			</h3>
			<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
				<div>
					<label for="log-level" class="block text-sm font-medium text-gray-700 mb-2">
						Уровень логирования
					</label>
					<select
						id="log-level"
						class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
					>
						<option value="debug">Debug</option>
						<option value="info" selected>Info</option>
						<option value="warn">Warning</option>
						<option value="error">Error</option>
					</select>
				</div>
				<div class="flex items-end">
					<button class="btn btn-primary w-full">
						Применить
					</button>
				</div>
				<div class="flex items-end">
					<button class="btn btn-secondary w-full">
						Очистить логи
					</button>
				</div>
			</div>
		</div>
	</div>
{:else}
	<div class="min-h-screen flex items-center justify-center">
		<div class="text-center">
			<h1 class="text-2xl font-bold text-gray-900 mb-4">Доступ запрещен</h1>
			<p class="text-gray-600 mb-4">Необходима авторизация для доступа к админ панели</p>
			<a href="/login" class="btn btn-primary">
				Войти в систему
			</a>
		</div>
	</div>
{/if}
