export type AppointmentStatus = 'PENDING' | 'CONFIRMED' | 'CANCELLED' | 'COMPLETED';

export interface Appointment {
  id: string; 
  patientId: string;
  doctorId: string;
  appointmentDate: string; 
  status: AppointmentStatus;
  duration: number; 
  Location: string;
  AppointmentType: 'Consultation' | 'Follow-up' | 'Check-up' | 'Emergency';
  notes?: string;
  Mode: 'Online' | 'In-Person';
  createdAt: string; 
  updatedAt: string; 

}

export type Appointments = Appointment[];