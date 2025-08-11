<script lang="ts">
  import { Search, AlertTriangle, Users } from 'lucide-svelte';
  import Button from '../Button.svelte';
  import Card from '../Card.svelte';
  import Badge from '../Badge.svelte';
  import { writable, derived } from 'svelte/store';

  interface Patient {
    id: string;
    name: string;
    age: number;
    condition: string;
    lastVisit: string;
    status: 'Critical' | 'Stable' | 'Needs Attention' | 'Recovering';
    priority: 'High' | 'Medium' | 'Low';
    avatar: string;
  }

  const searchTerm = writable('');

  const patients = writable<Patient[]>([
    {
      id: 'p1',
      name: 'Aarav Singh',
      age: 47,
      condition: 'Hypertension',
      lastVisit: '2025-07-12',
      status: 'Critical',
      priority: 'High',
      avatar: 'https://randomuser.me/api/portraits/men/32.jpg'
    },
    {
      id: 'p2',
      name: 'Priya Sharma',
      age: 34,
      condition: 'Type 2 Diabetes',
      lastVisit: '2025-06-22',
      status: 'Needs Attention',
      priority: 'Medium',
      avatar: 'https://randomuser.me/api/portraits/women/44.jpg'
    },
    {
      id: 'p3',
      name: 'Ravi Mehta',
      age: 58,
      condition: 'Cardiac Arrhythmia',
      lastVisit: '2025-07-10',
      status: 'Recovering',
      priority: 'Low',
      avatar: 'https://randomuser.me/api/portraits/men/54.jpg'
    },
    {
      id: 'p4',
      name: 'Neha Verma',
      age: 29,
      condition: 'Asthma',
      lastVisit: '2025-07-20',
      status: 'Stable',
      priority: 'Low',
      avatar: 'https://randomuser.me/api/portraits/women/65.jpg'
    }
  ]);

  const filteredPatients = derived(
    [patients, searchTerm],
    ([$patients, $searchTerm]) =>
      $patients.filter((p) =>
        p.name.toLowerCase().includes($searchTerm.toLowerCase())
      )
  );
</script>


<Card>
  <!-- Header -->
  <div class="flex items-center justify-between p-4 border-b">
    <div class="flex items-center gap-2 text-lg font-semibold">
      <Users class="w-5 h-5" />
      Patient List
    </div>
    <div class="relative w-[180px]">
      <Search class="absolute left-2 top-2.5 h-4 w-4 text-gray-400" />
      <input
        type="text"
        class="w-full pl-8 pr-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
        placeholder="Search patients..."
        on:input={(e) => searchTerm.set((e.target as HTMLInputElement).value)}
      />
    </div>
  </div>

  <!-- Patient Items -->
  <div class="p-4 space-y-4">
    {#each $filteredPatients as patient}
      <div class="flex items-center justify-between border rounded-lg p-3 hover:bg-gray-100 transition-colors">
        <div class="flex items-center gap-4">
          <img src={patient.avatar} alt={patient.name} class="w-10 h-10 rounded-full object-cover" />
          <div>
            <div class="flex items-center gap-2">
              <p class="font-medium">{patient.name}</p>
              {#if patient.status === 'Critical'}
                <AlertTriangle class="h-4 w-4 text-red-500" />
              {/if}
            </div>
            <p class="text-sm text-gray-500">
              Age {patient.age} • {patient.condition}
            </p>
            <p class="text-xs text-gray-400">Last visit: {patient.lastVisit}</p>
          </div>
        </div>

        <div class="flex items-center gap-3 text-right">
          <div>
            <Badge
              variant={
                patient.status === 'Critical'
                  ? 'red'
                  : patient.status === 'Needs Attention'
                  ? 'yellow'
                  : patient.status === 'Recovering'
                  ? 'green'
                  : 'blue'
              }
            >
              {patient.status}
            </Badge>
            <p class="text-xs mt-1 text-gray-400">
              Priority:
              <span
                class={
                  patient.priority === 'High'
                    ? 'text-red-500'
                    : patient.priority === 'Medium'
                    ? 'text-orange-500'
                    : 'text-green-500'
                }
              >
                {' '}{patient.priority}
              </span>
            </p>
          </div>
          <Button variant="outline" size="sm">View Profile</Button>
        </div>
      </div>
    {/each}

    <!-- View All CTA -->
    <div class="mt-4 flex justify-center">
      <Button variant="outline" size="sm">View All Patients</Button>
    </div>
  </div>
</Card>
