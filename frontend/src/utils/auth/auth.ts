export interface User {
    name: string;
    email: string;
}


export const parseJwt = (token: string): User | null => {
    try {
        const base64Url = token.split('.')[1];
        console.log(base64Url);
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(
            atob(base64)
                .split('')
                .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                .join('')
        );

        const jsonData = JSON.parse(jsonPayload);
        const res: User = {
            name: jsonData["name"] as string,
            email: jsonData["email"] as string,
        };
        return res
    } catch (error) {
        console.error('Error parsing JWT:', error);
        return null;
    }
};

