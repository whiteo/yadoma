import axios from 'axios';

const API_BASE_URL = '/yadoma/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authApi = {
  login: async (email, password) => {
    const response = await api.post('/authenticate', { email, password });
    return response.data;
  },

  register: async (email, password) => {
    const response = await api.post('/user/create', { email, password });
    return response.data;
  },

  validateToken: async () => {
    const response = await api.get('/authenticate');
    return response.data;
  },
};

export const userApi = {
  getCurrentUser: async () => {
    const response = await api.get('/user/me');
    return response.data;
  },

  getAllUsers: async () => {
    const response = await api.get('/user/all');
    return response.data;
  },

  deleteUser: async (userId) => {
    const response = await api.delete(`/user/delete/${userId}`);
    return response.data;
  },
};

export const containerApi = {
  getAllContainers: async (userId) => {
    const response = await api.get(`/container/${userId}/all`);
    return response.data;
  },

  getContainer: async (containerId) => {
    const response = await api.get(`/container/get/${containerId}`);
    return response.data;
  },

  createContainer: async (containerData) => {
    const response = await api.post('/container/create', containerData);
    return response.data;
  },

  startContainer: async (containerId) => {
    const response = await api.post(`/container/start/${containerId}`);
    return response.data;
  },

  stopContainer: async (containerId) => {
    const response = await api.post(`/container/stop/${containerId}`);
    return response.data;
  },

  restartContainer: async (containerId) => {
    const response = await api.post(`/container/restart/${containerId}`);
    return response.data;
  },

  deleteContainer: async (containerId) => {
    const response = await api.delete(`/container/delete/${containerId}`);
    return response.data;
  },

  getContainerLogsWsUrl: (containerId) => {
    const token = localStorage.getItem('token');
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    return `${protocol}//${window.location.host}/yadoma/ws/containers/${containerId}/logs?token=${encodeURIComponent(token)}`;
  },
};

export const systemApi = {
  getSystemInfo: async () => {
    const response = await api.get('/system/info');
    return response.data;
  },

  getDiskUsage: async () => {
    const response = await api.get('/system/disk-usage');
    return response.data;
  },
};

export default api;
