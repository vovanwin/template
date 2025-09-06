<script lang="ts">
	import { onMount } from 'svelte';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import Header from './Header.svelte';
	import '../app.css';
	import { initAuth } from '$lib/stores/auth';
	import '$lib/api/client';

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60 * 5, // 5 minutes
				refetchOnWindowFocus: false,
			},
		},
	});

	let { children } = $props();

	onMount(async () => {
		// Инициализируем аутентификацию при загрузке
		await initAuth();
	});
</script>

<QueryClientProvider client={queryClient}>
	<div class="min-h-screen flex flex-col bg-gray-50">
		<Header />

		<main class="flex-1">
			{@render children()}
		</main>

		<footer class="bg-white border-t border-gray-200 py-6">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<p class="text-center text-sm text-gray-500">
					Административная панель управления приложением
				</p>
			</div>
		</footer>
	</div>
</QueryClientProvider>
