<script lang="ts">
  import { Home, Calendar, FileText, Heart, Pill, MessageSquare, Phone, User, Activity } from "lucide-svelte";
  export let isOpen: boolean;

  const navItems = [
    { title: "Dashboard", icon: Home, url: "#" },
    { title: "My Appointments", icon: Calendar, url: "#" },
    { title: "Medical Records", icon: FileText, url: "#" },
    { title: "Prescriptions", icon: Pill, url: "#" },
    { title: "Health Monitoring", icon: Heart, url: "#" },
    { title: "Messages", icon: MessageSquare, url: "#" },
  ];

  const quickActions = [
  //  { title: "Book Appointment", icon: Calendar, url: "#" },
    { title: "Contact Doctor", icon: Phone, url: "#" },
    { title: "Update Profile", icon: User, url: "#" },
  ];
</script>
<aside
  class="h-full w-[250px] bg-white border-r shadow-md transition-all duration-300 ease-in-out flex flex-col"
  class:collapsed={!isOpen}
>
  
  <div class="flex items-center gap-2 px-4 py-4 border-b">
    <div class="h-8 w-8 bg-blue-600 text-white flex items-center justify-center rounded">
      <Activity class="w-4 h-4" />
    </div>
    {#if isOpen}
      <div class="text-md font-semibold">MediStream</div>
    {/if}
  </div>

 
  <nav class="flex-1 overflow-auto px-2 py-4 space-y-4">
    <div>
      <div class="text-xs text-gray-500 uppercase mb-2 px-2" class:hidden={!isOpen}>Navigation</div>
      {#each navItems as item}
        <a href={item.url}
          class="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-100 transition"
        >
          <item.icon class="w-5 h-5" />
          {#if isOpen}<span>{item.title}</span>{/if}
        </a>
      {/each}
    </div>

    <div>
      <div class="text-xs text-gray-500 uppercase mb-2 px-2" class:hidden={!isOpen}>Quick Actions</div>
      {#each quickActions as item}
        <a href={item.url}
          class="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-100 transition"
        >
          <item.icon class="w-5 h-5" />
          {#if isOpen}<span>{item.title}</span>{/if}
        </a>
      {/each}
    </div>
  </nav>

  <!-- Footer -->
  <div class="border-t p-4">
    <div class="flex items-center gap-2">
      <div class="h-8 w-8 bg-gray-300 rounded-full flex items-center justify-center">JD</div>
      {#if isOpen}
        <div class="flex-1">
          <div class="text-sm font-semibold">John Doe</div>
          <div class="text-xs text-gray-500">Patient ID: #12345</div>
        </div>
      {/if}
    </div>
  </div>
</aside>

<style>
  :global(.collapsed) {
    width: 64px !important;
  }
  aside {
    width: 250px;
  }
</style>
