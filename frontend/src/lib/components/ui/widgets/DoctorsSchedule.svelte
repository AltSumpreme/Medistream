<script lang="ts">
  import { Calendar, Clock, MapPin, Video, FileText, MoreHorizontal } from 'lucide-svelte';
  import Card from '../Card.svelte';
  import Badge from '../Badge.svelte';
  import Button from '../Button.svelte';
  import { writable } from 'svelte/store';

  export let appointments: Array<{
    id: string;
    time: string;
    patient: string;
    patientId: string;
    type: string;
    duration: string;
    room: string;
    status: 'Completed' | 'In Progress' | 'Upcoming';
    avatar?: string;
    notes: string;
  }> = [];

  const selectedAppointment = writable<string | null>(null);

  function getStatusColor(status: string) {
    switch (status) {
      case 'Completed': return 'green';
      case 'In Progress': return 'yellow';
      case 'Upcoming': return 'red';
     
    }
  }

  function getButtonText(a: typeof appointments[number]) {
    if (a.status === 'Completed') return 'View Notes';
    if (a.status === 'In Progress') return 'Continue';
    if (a.type === 'Video Consultation') return 'Join Call';
    return 'Start';
  }
</script>

<Card title="Today's Schedule" className="p-0 overflow-hidden">
  <div class="flex items-center justify-between p-6 border-b border-gray-200">
    <h2 class="flex items-center gap-2 text-lg font-semibold text-gray-900">
      <Calendar class="size-5" />
      Today's Schedule
    </h2>
    <Button variant="outline" size="sm">View Full Calendar</Button>
  </div>

  <div class="p-6 space-y-4">
    {#each appointments as a}
      <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
        <div class="flex items-center space-x-4">
          <div class="flex flex-col items-center">
            <Clock class="size-4 text-gray-400" />
            <span class="text-sm font-medium">{a.time}</span>
          </div>
          <img src={a.avatar || "/placeholder.svg"} alt={a.patient} class="size-12 rounded-full" />
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <p class="font-medium">{a.patient}</p>
              <span class="text-xs text-gray-500">{a.patientId}</span>
              {#if a.type === 'Video Consultation'}
                <Video class="size-4 text-blue-500" />
              {/if}
              {#if a.type === 'New Patient'}
                <Badge variant="blue">New</Badge>
              {/if}
            </div>
            <p class="text-sm text-gray-500">{a.type}</p>
            <div class="flex items-center gap-4 text-xs text-gray-500 mt-1">
              <div class="flex items-center gap-1">
                <MapPin class="size-3" />
                <span>{a.room}</span>
              </div>
              <div class="flex items-center gap-1">
                <FileText class="size-3" />
                <span>{a.notes}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex items-center space-x-3">
          <div class="text-right">
            <Badge variant={getStatusColor(a.status)}>{a.status}</Badge>
            <p class="text-xs text-gray-500 mt-1">{a.duration}</p>
          </div>

          <div class="flex flex-col gap-2">
            <Button
              variant={a.status === 'Upcoming' ? 'default' : 'outline'}
              size="sm"
            >
              {getButtonText(a)}
            </Button>

            <div class="relative">
              <button
                on:click={() =>
                  selectedAppointment.update((id) => id === a.id ? null : a.id)
                }
                class="p-1 text-gray-400 hover:text-gray-600"
              >
                <MoreHorizontal class="size-4" />
              </button>

              {#if $selectedAppointment === a.id}
                <div class="absolute right-0 top-full mt-1 w-48 bg-white border border-gray-200 rounded-md shadow-lg z-10">
                  <a href="#" class="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">View Patient Profile</a>
                  <a href="#" class="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">Access Medical Records</a>
                  <a href="#" class="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">Reschedule Appointment</a>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>
    {/each}
  </div>
</Card>
