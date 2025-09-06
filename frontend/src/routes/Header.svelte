<script lang="ts">
	import { page } from '$app/state';
	import { authStore } from '$lib/stores/auth';
</script>

<header class="bg-white shadow-sm border-b border-gray-200">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex justify-between items-center h-16">
			<!-- Logo -->
			<div class="flex items-center">
				<div class="flex-shrink-0">
					<h1 class="text-xl font-bold text-gray-900">Admin Panel</h1>
				</div>
			</div>

			<!-- Navigation -->
			<nav class="flex space-x-8">
				<a 
					href="/" 
					class="text-gray-900 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors {page.url.pathname === '/' ? 'bg-primary-100 text-primary-700' : ''}"
				>
					Главная
				</a>
				{#if $authStore.isAuthenticated}
					<a 
						href="/profile" 
						class="text-gray-900 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors {page.url.pathname === '/profile' ? 'bg-primary-100 text-primary-700' : ''}"
					>
						Профиль
					</a>
				{:else}
					<a 
						href="/login" 
						class="text-gray-900 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors {page.url.pathname === '/login' ? 'bg-primary-100 text-primary-700' : ''}"
					>
						Вход
					</a>
				{/if}
			</nav>

			<!-- User Info -->
			{#if $authStore.isAuthenticated}
				<div class="flex items-center space-x-4">
					<span class="text-sm text-gray-700">
						{$authStore.user?.email}
					</span>
					<div class="h-8 w-8 bg-primary-600 rounded-full flex items-center justify-center">
						<span class="text-white text-sm font-medium">
							{$authStore.user?.email?.charAt(0).toUpperCase()}
						</span>
					</div>
				</div>
			{/if}
		</div>
	</div>
</header>
