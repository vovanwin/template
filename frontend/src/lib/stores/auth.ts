import { writable } from 'svelte/store';
import type { AuthToken, UserMe } from '../api/models';
import { axiosInstance } from '../api/client';

interface AuthState {
	isAuthenticated: boolean;
	token: string | null;
	user: UserMe | null;
	loading: boolean;
	initialized: boolean;
}

const initialState: AuthState = {
	isAuthenticated: false,
	token: null,
	user: null,
	loading: false,
	initialized: false
};

export const authStore = writable<AuthState>(initialState);

export const setAuth = (token: string, user: UserMe) => {
	localStorage.setItem('auth_token', token);
	authStore.set({
		isAuthenticated: true,
		token,
		user,
		loading: false,
		initialized: true
	});
};

export const clearAuth = () => {
	localStorage.removeItem('auth_token');
	authStore.set({
		...initialState,
		initialized: true
	});
};

export const setLoading = (loading: boolean) => {
	authStore.update(state => ({ ...state, loading }));
};

export const initAuth = async () => {
	setLoading(true);
	
	const token = localStorage.getItem('auth_token');
	if (!token) {
		authStore.update(state => ({ 
			...state, 
			loading: false, 
			initialized: true 
		}));
		return null;
	}

	try {
		// Устанавливаем токен в store
		authStore.update(state => ({
			...state,
			token,
			isAuthenticated: true
		}));

		// Пытаемся получить данные пользователя
		const response = await axiosInstance.get('/auth/me');
		const user = response.data;

		authStore.update(state => ({
			...state,
			user,
			loading: false,
			initialized: true
		}));

		return token;
	} catch (error) {
		console.error('Failed to verify token:', error);
		// Если токен недействителен, очищаем его
		clearAuth();
		return null;
	}
};
