// Простой и надежный API клиент с JWT токеном
const API_BASE_URL = 'http://localhost:8080';

// Получить токен из localStorage
function getToken() {
  return localStorage.getItem('token');
}

// Сохранить токен в localStorage
function setToken(token) {
  if (token) {
    localStorage.setItem('token', token);
  } else {
    localStorage.removeItem('token');
  }
}

// Основная функция для API вызовов
export async function apiCall(url, options = {}) {
  const token = getToken();
  
  // Настройка заголовков
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };
  
  // Добавляем токен если есть
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers,
    });
    
    const contentType = response.headers.get('content-type');
    let data = {};
    
    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    }
    
    return {
      success: response.ok,
      data,
    };
  } catch (error) {
    console.error('API call error:', error);
    return {
      success: false,
      data: { error: 'Network error: ' + error.message },
    };
  }
}

// Авторизация
export async function login(username, password) {
  const result = await apiCall('/user/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
  
  if (result.success && result.data && result.data.token) {
    setToken(result.data.token);
    return { success: true, user_id: result.data.user_id };
  }
  
  return { success: false, error: result.data?.error || 'Ошибка входа' };
}

// Регистрация
export async function register(username, email, password) {
  const result = await apiCall('/user/register', {
    method: 'POST',
    body: JSON.stringify({ username, email, password }),
  });
  
  return result;
}

// Выход
export function logout() {
  setToken(null);
}

// Проверка авторизации
export function isAuthenticated() {
  return !!getToken();
}

export { getToken, setToken };
