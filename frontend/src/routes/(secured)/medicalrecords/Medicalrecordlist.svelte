<script lang="ts">
  import { FileText, User, Calendar, Eye, Edit, Download, MoreHorizontal, Stethoscope, TestTube, Pill, Camera } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import {createEventDispatcher} from 'svelte';

  const dispatch = createEventDispatcher();

  let searchQuery = '';
  let filterType = 'all';
  let filterDate = 'all';
  let selectedRecord: number | null = null;
  

  const allRecords = [
    {
      id: 1,
      type: 'consultation',
      title: 'Cardiology Consultation',
      patient: 'Alice Cooper',
      patientId: '#P-12345',
      date: new Date(2024, 0, 15),
      doctor: 'Dr. Sarah Wilson',
      status: 'completed',
      priority: 'normal',
      summary: 'Follow-up consultation for hypertension management',
      avatar: '/placeholder.svg?height=40&width=40'
    },
    // ... other records
  ];

  function getRecordIcon(type: string) {
    switch (type) {
      case 'consultation': return Stethoscope;
      case 'lab-result': return TestTube;
      case 'prescription': return Pill;
      case 'imaging': return Camera;
      default: return FileText;
    }
  }

  function getStatusColor(status: string) {
    switch (status) {
      case 'completed': return 'bg-green-100 text-green-800';
      case 'reviewed': return 'bg-blue-100 text-blue-800';
      case 'active': return 'bg-yellow-100 text-yellow-800';
      case 'pending': return 'bg-orange-100 text-orange-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }

  function getPriorityColor(priority: string) {
    switch (priority) {
      case 'high': return 'text-red-600';
      case 'normal': return 'text-gray-600';
      case 'low': return 'text-green-600';
      default: return 'text-gray-600';
    }
  }

  function toggleMenu(recordId: number) {
    selectedRecord = selectedRecord === recordId ? null : recordId;
  }

  function viewRecord(record: typeof allRecords[number]) {
    dispatch('view', record);
  }

  function viewPatientRecords(patientName: string) {
    console.log('View records for:', patientName);
  }

  $: filteredRecords = allRecords.filter(record => {
    const matchesSearch = searchQuery === '' ||
      record.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      record.patient.toLowerCase().includes(searchQuery.toLowerCase()) ||
      record.patientId.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesType = filterType === 'all' || record.type === filterType;
    const matchesDate = true;
    return matchesSearch && matchesType && matchesDate;
  });
</script>

<div class="min-h-screen bg-gray-50 p-8">
  <div class="max-w-6xl mx-auto">
    <div class="bg-white rounded-lg border border-gray-200 shadow">
      <div class="p-6 border-b border-gray-200 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-gray-900">Medical Records</h2>
        <span class="text-sm text-gray-500">{filteredRecords.length} records</span>
      </div>

      <div class="divide-y divide-gray-200">
        {#each filteredRecords as record}
          {@const IconComponent = getRecordIcon(record.type)}
          <div class="p-6 hover:bg-gray-50 transition-colors">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="flex flex-col items-center min-w-[60px]">
                  <div class="p-2 bg-blue-50 rounded-lg mb-1">
                    <IconComponent class="w-5 h-5 text-blue-600" />
                  </div>
                  <span class="text-xs text-gray-500 capitalize">{record.type.replace('-', ' ')}</span>
                </div>

                <img src={record.avatar} alt={record.patient} class="w-12 h-12 rounded-full" />
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <h3 class="font-medium text-gray-900">{record.title}</h3>
                    <span class="text-xs {getPriorityColor(record.priority)} font-medium">
                      {record.priority.toUpperCase()}
                    </span>
                  </div>

                  <div class="flex items-center gap-4 text-sm text-gray-500 mb-1">
                    <button on:click={() => viewPatientRecords(record.patient)} class="flex items-center gap-1 hover:text-blue-600">
                      <User class="w-3 h-3" />
                      <span>{record.patient} {record.patientId}</span>
                    </button>
                    <div class="flex items-center gap-1">
                      <Calendar class="w-3 h-3" />
                      <span>{record.date.toLocaleDateString()}</span>
                    </div>
                    <span>by {record.doctor}</span>
                  </div>

                  <p class="text-sm text-gray-600">{record.summary}</p>
                </div>
              </div>

              <div class="flex items-center gap-3">
                <span class="px-3 py-1 text-xs font-medium rounded-full {getStatusColor(record.status)}">
                  {record.status}
                </span>

                <button on:click={() => viewRecord(record)} class="flex items-center gap-1 px-3 py-2 text-sm font-medium text-blue-600 hover:bg-blue-50 rounded-lg">
                  <Eye class="w-4 h-4" /> View
                </button>

                <div class="relative">
                  <button on:click={() => toggleMenu(record.id)} class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg">
                    <MoreHorizontal class="w-4 h-4" />
                  </button>
                  {#if selectedRecord === record.id}
                    <div class="absolute right-0 top-full mt-1 w-48 bg-white border border-gray-200 rounded-lg shadow-lg z-10">
                      <div class="p-1">
                        <button class="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md">
                          <Edit class="w-4 h-4" /> Edit Record
                        </button>
                        <button class="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md">
                          <Download class="w-4 h-4" /> Download PDF
                        </button>
                        <button class="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md">
                          <User class="w-4 h-4" /> View Patient Profile
                        </button>
                      </div>
                    </div>
                  {/if}
                </div>
              </div>
            </div>
          </div>
        {/each}

        {#if filteredRecords.length === 0}
          <div class="p-12 text-center">
            <FileText class="w-12 h-12 text-gray-300 mx-auto mb-4" />
            <h3 class="text-lg font-medium text-gray-900 mb-2">No records found</h3>
            <p class="text-gray-500">Try adjusting your search or filter criteria.</p>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
