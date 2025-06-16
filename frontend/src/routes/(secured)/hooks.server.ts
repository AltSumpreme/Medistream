
import type { Handle } from '@sveltejs/kit';
export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('token');

	if (token) {
		const res = await fetch('http://localhost:8080/auth/verify', {
			headers: { Authorization: `Bearer ${token}` }
		});
		if (res.ok) {
			const data = await res.json();
			event.locals.user = {
				user_id: data.user_id,
				role: data.role
			};
		}
	}
	return resolve(event);
};
