const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export async function registerUser(payload) {
    const response = await fetch(`${API_BASE_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
    });

    if (!response.ok) {
        throw new Error('Registration failed');
    }

    return response;
}
