<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import { Calendar, Clock, MapPin, User, FileText } from 'lucide-svelte';
	import { createAppointment } from '../../../lib/services/appointments';
	import type { Appointment } from '../../../lib/types/appointment';

	export let onClose: () => void;
	export let userRole: 'PATIENT' | 'DOCTOR' | 'ADMIN'; // pass in from parent

	let isSubmitting = false;
	let error = '';

	let formData: Partial<Appointment> = {
		patientId: '',
		appointmentDate: '',
		duration: 15,
		notes: '',
		Mode: 'In-Person',
		AppointmentType: 'Consultation'
	};

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		error = '';
		isSubmitting = true;

		try {
			// If patient, do not allow location to be set from client
			if (userRole === 'PATIENT') {
				delete formData.Location;
			}
			await createAppointment(formData);
			onClose();
		} catch (err) {
			console.error('Error creating appointment:', err);
			error = 'Failed to create appointment. Please try again.';
		} finally {
			isSubmitting = false;
		}
	};
</script>

<Modal title="Schedule Appointment" {onClose}>
	<form on:submit={handleSubmit} class="space-y-4">
		<!-- Patient -->
		<div>
			<label class="mb-1 flex items-center gap-1 text-sm font-medium">
				<User class="w-4" /> Patient
			</label>
			<input required bind:value={formData.patientId} class="w-full rounded border p-2" />
		</div>

		<!-- Date & Time -->
		<div class="grid grid-cols-2 gap-4">
			<div>
				<label class="mb-1 flex items-center gap-1 text-sm font-medium">
					<Calendar class="w-4" /> Date
				</label>
				<input type="date" required bind:value={formData.appointmentDate} class="w-full rounded border p-2" />
			</div>
			<div>
				<label class="mb-1 flex items-center gap-1 text-sm font-medium">
					<Clock class="w-4" /> Time
				</label>
				<input type="time" required bind:value={formData} class="w-full rounded border p-2" />
			</div>
		</div>

		<!-- Type & Duration -->
		<div class="grid grid-cols-2 gap-4">
			<div>
				<label for="type-select" class="mb-1 flex text-sm font-medium">Type</label>
				<select id="type-select" bind:value={formData.AppointmentType} class="w-full rounded border p-2">
					<option value="Consultation">Consultation</option>
					<option value="Follow-up">Follow-up</option>
					<option value="Check-up">Check-up</option>
					<option value="Emergency">Emergency</option>
				</select>
			</div>
			<div>
				<label for="duration-select" class="mb-1 flex text-sm font-medium">Duration</label>
				<select id="duration-select" bind:value={formData.duration} class="w-full rounded border p-2">
					<option>15</option>
					<option>30</option>
					<option>45</option>
					<option>60</option>
				</select>
			</div>
		</div>

		<!-- Room (only for Doctor/Admin) -->
		{#if userRole === 'DOCTOR' || userRole === 'ADMIN'}
			<div>
				<label class="mb-1 flex items-center gap-1 text-sm font-medium">
					<MapPin class="w-4" /> Room
				</label>
				<select bind:value={formData.Location} class="w-full rounded border p-2">
					<option value="Room 101">Room 101</option>
					<option value="Room 102">Room 102</option>
					<option value="Online">Online</option>
				</select>
			</div>
		{/if}

		<!-- Notes -->
		<div>
			<label class="mb-1 flex items-center gap-1 text-sm font-medium">
				<FileText class="w-4" /> Notes
			</label>
			<textarea bind:value={formData.notes} class="w-full resize-none rounded border p-2" rows="3"></textarea>
		</div>

		<div class="flex justify-end gap-2 pt-4 border-t">
			<button type="button" on:click={onClose} class="rounded border px-4 py-2">Cancel</button>
			<button type="submit" class="rounded bg-blue-600 px-4 py-2 text-white" disabled={isSubmitting}>
				{isSubmitting ? 'Scheduling...' : 'Schedule'}
			</button>
		</div>
	</form>
	{#if error}
		<div class="text-red-500 mt-2">
			{error}
		</div>
	{/if}
</Modal>
