<script lang="ts">
  import { Clock, User, MapPin, Phone, Video, MoreHorizontal, Edit, Trash2, CheckCircle } from 'lucide-svelte';
  import type { Appointment } from '$lib/types/appointment';
 

  export let searchQuery: string = '';
  export let filterStatus: string = 'all';
  export let data: {
    user: { role: 'PATIENT' | 'DOCTOR' | 'ADMIN'; user_id: string },
    appointments: Appointment[],
  };

  let selectedAppointmentId: string | null = null;
  let appointments: Appointment[] = data.appointments;

 // const { user } = data;



  

  function getStatusColor(status: string) {
    switch (status) {
      case 'confirmed': return 'bg-blue-100 text-blue-800';
      case 'completed': return 'bg-green-100 text-green-800';
      case 'upcoming': return 'bg-yellow-100 text-yellow-800';
      case 'cancelled': return 'bg-red-100 text-red-800';
      case 'no-show': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }

  function getButtonText(appointment: Appointment) {
    switch (appointment.status) {
      case 'COMPLETED': return 'View Notes';
      case 'PENDING': return appointment.AppointmentType === 'Consultation' ? 'Join Call' : 'Start';
      case 'CONFIRMED': return 'Start';
      default: return 'View';
    }
  }

  function getButtonColor(status: string) {
    if (status === 'upcoming' || status === 'confirmed') {
      return 'bg-blue-600 text-white hover:bg-blue-700';
    }
    return 'border border-gray-300 text-gray-700 hover:bg-gray-50';
  }

  function toggleMenu(appointmentId: string) {
    selectedAppointmentId = selectedAppointmentId === appointmentId ? null : appointmentId;
  }

  $: filteredAppointments = appointments.filter((appointment: Appointment) => {
    const matchesSearch =
      searchQuery === '' || appointment.patientId.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesFilter = filterStatus === 'all' || appointment.status === filterStatus;

    return matchesSearch && matchesFilter;
  });
</script>

<!-- UI -->
<div class="bg-white rounded-lg border border-gray-200">
  <div class="p-6 border-b border-gray-200 flex items-center justify-between">
    <h2 class="text-lg font-semibold text-gray-900">All Appointments</h2>
    <span class="text-sm text-gray-500">{filteredAppointments.length} appointments</span>
  </div>

  <div class="divide-y divide-gray-200">
    {#each filteredAppointments as appointment}
      <div class="p-6 hover:bg-gray-50 transition-colors">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">

            <!-- Time -->
            <div class="flex flex-col items-center min-w-[60px]">
              <Clock class="size-4 text-gray-400 mb-1" />
              <span class="text-sm font-medium">
                {new Date(appointment.appointmentDate).toLocaleTimeString('en-US', {
                  hour: 'numeric',
                  minute: '2-digit',
                })}
              </span>
              <span class="text-xs text-gray-500">
                {new Date(appointment.appointmentDate).toLocaleDateString('en-US', {
                  month: 'short',
                  day: 'numeric',
                  year: 'numeric'
                })}
              </span>
            </div>

            <!-- Patient Info -->
            <img src={"/placeholder.svg"} alt="avatar" class="size-12 rounded-full" />
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <h3 class="font-medium text-gray-900">{appointment.patientId}</h3>
                {#if appointment.Mode === 'Online'}
                  <Video class="size-4 text-blue-500" />
                {:else if appointment.Mode === 'In-Person'}
                  <span class="px-2 py-1 text-xs font-medium bg-green-100 text-green-800 rounded-full">In-Person</span>
                {/if}
              </div>

              <div class="flex items-center gap-4 text-sm text-gray-500">
                <span>{appointment.Mode}</span>
                <div class="flex items-center gap-1">
                  <MapPin class="size-3" />
                  <span>{appointment.Location}</span>
                </div>
                <span>{appointment.duration} min</span>
              </div>

              {#if appointment.notes}
                <p class="text-sm text-gray-600 mt-1">{appointment.notes}</p>
              {/if}
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-3">
            <span class="px-3 py-1 text-xs font-medium rounded-full {getStatusColor(appointment.status)}">
              {appointment.status}
            </span>

            <button class="px-4 py-2 text-sm font-medium rounded-lg transition-colors {getButtonColor(appointment.status)}">
              {getButtonText(appointment)}
            </button>

            <!-- Menu -->
            <div class="relative">
              <button
                on:click={() => toggleMenu(String(appointment.id))}
                class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <MoreHorizontal class="size-4" />
              </button>

              {#if selectedAppointmentId === String(appointment.id)}
                <div class="absolute right-0 top-full mt-1 w-48 bg-white border border-gray-200 rounded-lg shadow-lg z-10">
                  <div class="p-1">
                    <button class="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md">
                      <Edit class="size-4" />
                      Edit Appointment
                    </button>
                    <button class="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md">
                      <CheckCircle class="size-4" />
                      Mark as Completed
                    </button>
                    <button class="flex items-center gap-2 w-full px-3 py-2 text-sm text-red-600 hover:bg-red-50 rounded-md">
                      <Trash2 class="size-4" />
                      Cancel Appointment
                    </button>
                  </div>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>
    {/each}

    {#if filteredAppointments.length === 0}
      <div class="p-12 text-center">
        <Clock class="size-12 text-gray-300 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">No appointments found</h3>
        <p class="text-gray-500">Try adjusting your search or filter criteria.</p>
      </div>
    {/if}
  </div>
</div>
