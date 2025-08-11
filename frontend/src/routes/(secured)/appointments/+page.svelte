<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import Header from '$lib/components/ui/Header.svelte';
  import SideBar from '$lib/components/ui/SideBar.svelte';
  import AppointmentList from './AppointmentList.svelte';
  import CalendarView from './CalendarView.svelte';
  import AppointmentModal from './AppointmentModal.svelte';
  import type { Appointment } from '$lib/types/appointment';


export let data: {
    user: {
      role: 'PATIENT' | 'DOCTOR' | 'ADMIN';
      user_id: string;
    };
    appointments: Appointment[];
  };

  const { user, appointments } = data;
  const { role, user_id: userId } = user;


  let isSidebarOpen = true;
  let showModal = false;
  let selectedDate = new Date();
  let searchQuery = '';
  let filterStatus = 'all';

  const toggleSidebar = () => (isSidebarOpen = !isSidebarOpen);
  const dispatch = createEventDispatcher();
  const closeModal = () => {
    showModal = false;
    dispatch('close');
  };

function handleFoo(action: string, data?: any) {
  switch (action) {
    case 'previousMonth':
      console.log('Navigated to previous month');
      break;

    case 'nextMonth':
      console.log('Navigated to next month');
      break;

    case 'selectDate':
      console.log('Selected date:', data);
      selectedDate = data;
      break;

    default:
      console.warn('Unknown action in handleFoo:', action, data);
  }
}

</script>

<div class="flex flex-col h-screen bg-gradient-to-br from-slate-50 to-blue-50">
  <Header className="shrink-0" on:toggleSidebar={toggleSidebar} />

  <div class="flex flex-1 overflow-hidden">
    <SideBar isOpen={isSidebarOpen} role={role} className="shrink-0 h-full" />

    <main class="flex-1 overflow-y-auto p-6 max-w-7xl mx-auto space-y-8">
      <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
        <div class="flex flex-col lg:flex-row gap-4 items-stretch lg:items-center">
          <button
            on:click={() => (showModal = true)}
            class="w-full lg:w-auto bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white px-6 py-3 rounded-xl font-semibold shadow-lg hover:shadow-xl transition-all duration-200 flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            New Appointment
          </button>

          <div class="relative flex-1">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              bind:value={searchQuery}
              placeholder="Search appointments..."
              class="w-full pl-12 pr-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none transition-all bg-gray-50 hover:bg-white focus:bg-white"
            />
          </div>

          <div class="relative flex-shrink-0">
            <select
              bind:value={filterStatus}
              class="bg-white border border-gray-200 rounded-xl px-4 py-3 pr-10 focus:ring-2 focus:ring-blue-500 transition-all hover:border-gray-300 cursor-pointer min-w-[140px]"
            >
              <option value="all">All Status</option>
              <option value="upcoming">Upcoming</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
            <div class="absolute inset-y-0 right-0 flex items-center px-3 pointer-events-none">
              <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 xl:grid-cols-3 gap-8">
        <div class="xl:col-span-2">
          <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
            <div class="p-6 border-b bg-gradient-to-r from-gray-50 to-gray-100">
              <h2 class="text-xl font-semibold text-gray-900 flex items-center gap-2">
                <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v4a2 2 0 002 2h2m0-6v6m2-6h8a2 2 0 012 2v4a2 2 0 01-2 2h-8m0-6v6" />
                </svg>
                Appointment List
              </h2>
            </div>
            <AppointmentList {searchQuery} {filterStatus} {data} />
          </div>
        </div>

        <div class="xl:col-span-1">
          <div class="bg-white rounded-2xl shadow-lg border border-gray-100 sticky top-6 overflow-hidden">
            <div class="p-6 border-b bg-gradient-to-r from-gray-50 to-gray-100">
              <h2 class="text-xl font-semibold text-gray-900 flex items-center gap-2">
                <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                Calendar
              </h2>
            </div>
            <CalendarView bind:selectedDate {searchQuery} {filterStatus} onfoo={handleFoo} />
          </div>
        </div>
      </div>

      {#if showModal}
        <div class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <AppointmentModal onClose={closeModal} userRole={role} />
        </div>
      {/if}
    </main>
  </div>
</div>
