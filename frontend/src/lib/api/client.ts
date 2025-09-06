import axios from 'axios';
import type { AxiosRequestConfig } from 'axios';
import { authStore, clearAuth } from '../stores/auth';
import { get } from 'svelte/store';

// Создаем экземпляр axios с базовой конфигурацией
const axiosInstance = axios.create({
	baseURL: 'http://localhost:8080',
	timeout: 10000,
	headers: {
		'Content-Type': 'application/json',
	},
});

// Интерцептор для добавления токена к запросам
axiosInstance.interceptors.request.use(
	(config) => {
		const auth = get(authStore);
		if (auth.token) {
			config.headers.Authorization = `Bearer ${auth.token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

// Интерцептор для обработки ошибок авторизации
axiosInstance.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response?.status === 401) {
			// При получении 401 ошибки очищаем авторизацию
			clearAuth();
			// Перенаправляем на страницу логина
			if (typeof window !== 'undefined') {
				window.location.href = '/login';
			}
		}
		return Promise.reject(error);
	}
);

// Экспортируем функцию для Orval
export const apiClient = <T = any>(config: AxiosRequestConfig): Promise<T> => {
	return axiosInstance(config);
};

// Экспортируем экземпляр для прямого использования
export { axiosInstance };
