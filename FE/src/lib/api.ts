export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/query";

export const apiRequest = async (query: string, variables: any = {}, token?: string | null) => {
    const headers: HeadersInit = {
        "Content-Type": "application/json",
    };

    if (token) {
        headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(API_URL, {
        method: "POST",
        headers,
        body: JSON.stringify({ query, variables }),
    });

    const result = (await response.json()) as { data?: any; errors?: any[] };
    if (result.errors) {
        throw new Error(result.errors[0]?.message || "Something went wrong");
    }

    return result.data;
};
