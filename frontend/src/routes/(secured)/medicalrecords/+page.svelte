<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import Medicalrecordlist from './Medicalrecordlist.svelte';
  import Medicalrecorddetail from './Medicalrecorddetail.svelte';
  import Medicalrecordmodal from './Medicalrecordmodal.svelte';
  import Header from '$lib/components/ui/Header.svelte';
  import SideBar from '$lib/components/ui/SideBar.svelte';

  let selectedRecord: any = null;
  let showModal = false;
  let isSidebarOpen = true;
  let searchQuery = '';
  const dispatch = createEventDispatcher();

  const toggleSidebar = () => (isSidebarOpen = !isSidebarOpen);

  function viewRecord(record: any) {
    selectedRecord = record;
  }

  function backToList() {
    selectedRecord = null;
  }

  function openModal() {
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    dispatch('close');
  }
</script>

<div class="flex flex-col h-screen bg-gradient-to-br from-slate-50 to-blue-50">
  <!-- Header -->
  <header class="shrink-0">
    <Header on:toggleSidebar={toggleSidebar} />
  </header>

  <!-- Body -->
  <div class="flex flex-1 overflow-hidden">
    <!-- Sidebar -->
    <div class="h-full shrink-0">
      <SideBar isOpen={isSidebarOpen} role="doctor" />
    </div>

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto p-6 max-w-7xl mx-auto space-y-8">
      <!-- Controls -->
      <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
        <div class="flex flex-col lg:flex-row gap-4 items-stretch lg:items-center">
          <!-- New Record Button -->
          <button
            on:click={openModal}
            class="w-full lg:w-auto bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white px-6 py-3 rounded-xl font-semibold shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200 flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            New Record
          </button>

          <!-- Search Input -->
          <div class="flex-1 relative">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              bind:value={searchQuery}
              placeholder="Search medical records..."
              class="w-full pl-12 pr-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all duration-200 bg-gray-50 hover:bg-white focus:bg-white"
            />
          </div>
        </div>
      </div>

      <!-- Main Grid -->
      <div class="grid grid-cols-1 gap-8">
        <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
          <div class="p-6 border-b border-gray-100 bg-gradient-to-r from-gray-50 to-gray-100 flex items-center gap-2">
            <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v4a2 2 0 002 2h2m0-6v6m2-6h8a2 2 0 012 2v4a2 2 0 01-2 2h-8m0-6v6" />
            </svg>
            <h2 class="text-xl font-semibold text-gray-900">Medical Records</h2>
          </div>
          <div class="p-6">
            {#if selectedRecord}
              <Medicalrecorddetail record={selectedRecord} backToList={backToList} />
            {:else}
              <Medicalrecordlist on:view={e => viewRecord(e.detail)} />
            {/if}
          </div>
        </div>
      </div>

      <!-- Modal -->
      {#if showModal}
        <div class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <Medicalrecordmodal onClose={closeModal} />
        </div>
      {/if}
    </main>
  </div>
</div>
