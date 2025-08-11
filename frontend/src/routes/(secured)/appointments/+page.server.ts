import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { Appointment } from '$lib/types/appointment';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const user = locals.user;
	if (!user) throw redirect(302, '/login');
	

	let appointments: Appointment[] = [];

	const urlMap: Record<string, string> = {
		PATIENT: `http://localhost:8080/appointments/patient/${user.user_id}`,
		DOCTOR: `http://localhost:8080/appointments/doctor/${user.user_id}`,
		ADMIN: `http://localhost:8080/appointments`
	};

	const headers: HeadersInit = {
		'Authorization': `Bearer ${user.token}`
	};

	const url = urlMap[user.role];
	const res = await fetch(url, user.role === 'DOCTOR' ? undefined : { headers });

	if (res.ok) {
		const result = await res.json();
		appointments = result.appointments || result; 
	}

	return {
		user,
		appointments
	};
};
