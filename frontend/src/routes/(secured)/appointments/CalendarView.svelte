<script lang="ts">
  import { ChevronLeft, ChevronRight, Clock, User, MapPin } from 'lucide-svelte';
  import type { Appointment, AppointmentStatus } from '../../../lib/types/appointment'; 

  export let selectedDate: Date;
  export let searchQuery: string;
  export let filterStatus: string;
  export let onfoo: ((action: string, data?: any) => void) | undefined;

  let currentMonth = new Date(selectedDate);

  const appointments: Appointment[] = [
    {
      id: '1',
      patientId: 'patient_1',
      doctorId: 'doctor_1',
      appointmentDate: new Date(2025, 0, 15, 9, 0).toISOString(),
      status: 'CONFIRMED',
      duration: 30,
      Location: 'Room 101',
      AppointmentType: 'Consultation',
      notes: 'Bring previous reports',
      Mode: 'In-Person',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    },
    {
      id: '2',
      patientId: 'patient_2',
      doctorId: 'doctor_2',
      appointmentDate: new Date(2025, 0, 15, 10, 30).toISOString(),
      status: 'COMPLETED',
      duration: 15,
      Location: 'Room 102',
      AppointmentType: 'Follow-up',
      Mode: 'In-Person',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    },
    {
      id: '3',
      patientId: 'patient_3',
      doctorId: 'doctor_3',
      appointmentDate: new Date(2024, 0, 16, 14, 0).toISOString(),
      status: 'PENDING',
      duration: 45,
      Location: 'Room 103',
      AppointmentType: 'Check-up',
      Mode: 'In-Person',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    },
    {
      id: '4',
      patientId: 'patient_4',
      doctorId: 'doctor_4',
      appointmentDate: new Date(2024, 0, 17, 11, 0).toISOString(),
      status: 'PENDING',
      duration: 30,
      Location: 'Online',
      AppointmentType: 'Emergency',
      Mode: 'Online',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }
  ];

  // Mock patient name mapping for display
  const patientMap: Record<string, string> = {
    patient_1: 'Alice Cooper',
    patient_2: 'Bob Smith',
    patient_3: 'Carol White',
    patient_4: 'David Lee'
  };

  function getDaysInMonth(date: Date) {
    const days: (Date | null)[] = [];
    const year = date.getFullYear();
    const month = date.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDate = new Date(year, month + 1, 0).getDate();

    for (let i = 0; i < firstDay.getDay(); i++) {
      days.push(null);
    }

    for (let d = 1; d <= lastDate; d++) {
      days.push(new Date(year, month, d));
    }

    return days;
  }

  function filterAppointments(date: Date) {
    return appointments.filter(
      a =>
        new Date(a.appointmentDate).toDateString() === date.toDateString() &&
        (filterStatus === 'all' || a.status.toLowerCase() === filterStatus.toLowerCase()) &&
        (!searchQuery ||
          patientMap[a.patientId]?.toLowerCase().includes(searchQuery.toLowerCase()))
    );
  }

  function getStatusColor(status: AppointmentStatus | string) {
    const colors: Record<string, string> = {
      PENDING: 'bg-yellow-100 text-yellow-800',
      CONFIRMED: 'bg-blue-100 text-blue-800',
      COMPLETED: 'bg-green-100 text-green-800',
      CANCELLED: 'bg-red-100 text-red-800'
    };
    return colors[status.toUpperCase()] ?? 'bg-gray-100 text-gray-800';
  }

  function formatTime(isoDate: string) {
    return new Date(isoDate).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  function previousMonth() {
    currentMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1);
    onfoo?.('previousMonth');
  }

  function nextMonth() {
    currentMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1);
    onfoo?.('nextMonth');
  }

  function selectDate(date: Date) {
    selectedDate = date;
    onfoo?.('selectDate', date);
  }

  $: days = getDaysInMonth(currentMonth);
  $: monthLabel = currentMonth.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  $: selectedAppointments = filterAppointments(selectedDate);
</script>

<div class="bg-white rounded-2xl border border-gray-200 overflow-hidden">
  <!-- Header -->
  <div class="flex justify-between items-center px-6 py-4 border-b">
    <h2 class="text-lg font-semibold text-gray-900">{monthLabel}</h2>
    <div class="flex gap-2">
      <button on:click={previousMonth} class="p-2 hover:bg-gray-100 rounded-lg">
        <ChevronLeft class="size-4" />
      </button>
      <button on:click={nextMonth} class="p-2 hover:bg-gray-100 rounded-lg">
        <ChevronRight class="size-4" />
      </button>
    </div>
  </div>

  <!-- Weekday Headings -->
  <div class="grid grid-cols-7 px-6 pt-4 text-sm font-medium text-gray-500">
    {#each ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as day}
      <div class="text-center">{day}</div>
    {/each}
  </div>

  <!-- Days Grid -->
  <div class="grid grid-cols-7 gap-1 p-6 pt-2">
    {#each days as date}
      {#if date}
        <button
          type="button"
          class="min-h-[100px] w-full text-left border border-gray-100 rounded-lg p-2 hover:bg-gray-50 cursor-pointer focus:outline-none focus:ring-2 focus:ring-blue-400"
          on:click={() => selectDate(date)}
        >
          <div class="flex justify-between items-center mb-1">
            <span class="text-sm font-semibold {selectedDate.toDateString() === date.toDateString() ? 'text-blue-600' : 'text-gray-800'}">
              {date.getDate()}
            </span>
            {#if filterAppointments(date).length}
              <span class="text-xs bg-blue-100 text-blue-800 px-1 rounded-full">
                {filterAppointments(date).length}
              </span>
            {/if}
          </div>
          <div class="space-y-1">
            {#each filterAppointments(date).slice(0, 2) as a}
              <div class="text-xs truncate px-1 rounded {getStatusColor(a.status)}">
                {formatTime(a.appointmentDate)} - {patientMap[a.patientId]}
              </div>
            {/each}
            {#if filterAppointments(date).length > 2}
              <div class="text-xs text-gray-400">+{filterAppointments(date).length - 2} more</div>
            {/if}
          </div>
        </button>
      {:else}
        <div class="min-h-[100px] border border-gray-100 rounded-lg p-2"></div>
      {/if}
    {/each}
  </div>

  <!-- Selected Date Appointments -->
  <div class="border-t px-6 py-4">
    <h3 class="text-lg font-semibold mb-3">
      Appointments for {selectedDate.toLocaleDateString('en-US', {
        weekday: 'long',
        month: 'long',
        day: 'numeric'
      })}
    </h3>

    {#if selectedAppointments.length === 0}
      <p class="text-gray-500">No appointments.</p>
    {:else}
      <div class="space-y-3">
        {#each selectedAppointments as a}
          <div class="p-3 border border-gray-100 rounded-lg flex justify-between items-center hover:bg-gray-50">
            <div class="flex items-start gap-4">
              <div class="text-sm flex flex-col items-center">
                <Clock class="size-4 text-gray-400 mb-1" />
                <span>{formatTime(a.appointmentDate)}</span>
              </div>
              <div>
                <div class="flex items-center gap-2 text-sm font-medium">
                  <User class="size-4 text-gray-400" />
                  {patientMap[a.patientId]}
                </div>
                <div class="text-xs text-gray-500 mt-1 flex gap-3">
                  <span>{a.AppointmentType}</span>
                  <span class="flex items-center gap-1">
                    <MapPin class="size-3" /> {a.Location}
                  </span>
                  <span>{a.duration} min</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <span class="text-xs px-2 py-1 rounded-full font-medium {getStatusColor(a.status)}">
                {a.status}
              </span>
              <button class="text-sm text-blue-600 hover:underline">Details</button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
