import api from "$lib/config/api";
import type { Appointment } from "$lib/types/appointment";

export async function createAppointment(payload: Partial<Appointment>) {
  const res = await api.post("/appointments", payload);
  return res.data;
}

export async function getAllAppointments(page = 1, limit = 10) {
  const res = await api.get(`/appointments?page=${page}&limit=${limit}`);
  return res.data;
}

export async function getAppointmentById(id: string) {
  const res = await api.get(`/appointments/${id}`);
  return res.data;
}

export async function updateAppointment(id: string, payload: Appointment) {
  const res = await api.put(`/appointments/${id}`, payload);
  return res.data;
}

export async function deleteAppointment(id: string) {
  const res = await api.delete(`/appointments/${id}`);
  return res.data;
}

export async function rescheduleAppointment(id: string, payload: { date: string; duration: number }) {
  const res = await api.put(`/appointments/${id}/reschedule`, payload);
  return res.data;
}

export async function cancelAppointment(id: string) {
  const res = await api.post(`/appointments/${id}/cancel`);
  return res.data;
}

export async function changeAppointmentStatus(id: string, payload: Appointment) {
  const res = await api.post(`/appointments/${id}/status`, payload);
  return res.data;
}

export async function getAppointmentsByDoctorId(doctorId: string, limit = 10, offset = 0) {
  const res = await api.get(`/appointments/doctor/${doctorId}?limit=${limit}&offset=${offset}`);
  return res.data;
}

export async function getAppointmentsByPatientId(patientId: string, limit = 10, offset = 0) {
  const res = await api.get(`/appointments/patient/${patientId}?limit=${limit}&offset=${offset}`);
  return res.data;
}
