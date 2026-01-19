const API_URL = import.meta.env.VITE_API_URL;
//this is what i had for local development
// const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export async function fetchAPI(endpoint: string, options: RequestInit = {}) {
    const response = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,

        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.statusText}`);

    }
    return response.json()
}

