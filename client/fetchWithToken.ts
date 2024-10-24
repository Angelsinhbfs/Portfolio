const fetchWithToken = async (url: string, options: RequestInit = {}) => {

    let tokenCookie = sessionStorage.getItem('token');

    // Add the token to the request headers
    const headers = {
        ...options.headers,
        'Authorization': `Bearer ${tokenCookie}`,
    };

    try {
        const response = await fetch(url, { ...options, headers });
        if (!response.ok) {
            throw new Error('Request failed');
        }
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
            return response.json();
        }
        return 0;
    } catch (error) {
        console.error('Error fetching data:', error);
        throw error;
    }
};

export default fetchWithToken;