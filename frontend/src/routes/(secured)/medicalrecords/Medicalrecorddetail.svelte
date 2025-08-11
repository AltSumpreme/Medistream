<script lang="ts">
  import { ArrowLeft, Edit, Download, Share, Calendar, User, Stethoscope, FileText, Clock, Tag } from 'lucide-svelte';
  import { onMount } from 'svelte';

  export let record: {
    id: number;
    type: 'consultation' | 'prescription' | 'lab';
    title: string;
    patient: string;
    patientId: string;
    doctor: string;
    doctorId: string;
    department: string;
    createdAt: string;
    duration?: string;
    notes?: string;
  };

  export let backToList: () => void;

  function handleEdit() {
    // Navigate to edit form
  }

  function handleDownload() {
    // Trigger PDF or record download
  }

  function handleShare() {
    // Trigger share functionality
  }
</script>

<section class="p-6 bg-white shadow rounded-lg space-y-6 max-w-4xl mx-auto">
  <div class="flex items-center justify-between">
    <button on:click={backToList} class="text-sm text-gray-600 hover:text-black flex items-center gap-1">
      <ArrowLeft class="w-4 h-4" /> Back to Records
    </button>
    <div class="flex gap-3">
      <button on:click={handleEdit} class="text-blue-600 hover:underline flex items-center gap-1">
        <Edit class="w-4 h-4" /> Edit
      </button>
      <button on:click={handleDownload} class="text-green-600 hover:underline flex items-center gap-1">
        <Download class="w-4 h-4" /> Download
      </button>
      <button on:click={handleShare} class="text-purple-600 hover:underline flex items-center gap-1">
        <Share class="w-4 h-4" /> Share
      </button>
    </div>
  </div>

  <div class="space-y-1">
    <h2 class="text-2xl font-semibold">{record.title}</h2>
    <div class="text-sm text-gray-500 flex items-center gap-2">
      <Tag class="w-4 h-4" /> {record.type.charAt(0).toUpperCase() + record.type.slice(1)}
    </div>
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm text-gray-800">
    <div class="flex items-center gap-2">
      <User class="w-4 h-4 text-gray-500" /> <span>Patient:</span> <strong>{record.patient}</strong> ({record.patientId})
    </div>
    <div class="flex items-center gap-2">
      <Stethoscope class="w-4 h-4 text-gray-500" /> <span>Doctor:</span> <strong>{record.doctor}</strong> ({record.doctorId})
    </div>
    <div class="flex items-center gap-2">
      <FileText class="w-4 h-4 text-gray-500" /> <span>Department:</span> <strong>{record.department}</strong>
    </div>
    <div class="flex items-center gap-2">
      <Calendar class="w-4 h-4 text-gray-500" /> <span>Date:</span> <strong>{new Date(record.createdAt).toLocaleDateString()}</strong>
    </div>
    {#if record.duration}
      <div class="flex items-center gap-2">
        <Clock class="w-4 h-4 text-gray-500" /> <span>Duration:</span> <strong>{record.duration}</strong>
      </div>
    {/if}
  </div>

  {#if record.notes}
    <div class="mt-4 border-t pt-4">
      <h3 class="text-lg font-semibold mb-2">Notes</h3>
      <p class="text-gray-700 whitespace-pre-wrap">{record.notes}</p>
    </div>
  {/if}
</section>

<style>
  section {
    font-family: system-ui, sans-serif;
  }
</style>
