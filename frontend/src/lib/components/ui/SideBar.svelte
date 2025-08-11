<script lang="ts">
  import {
    Home,
    Calendar,
    FileText,
    Heart,
    Pill,
    Phone,
    User,
    Activity,
    Stethoscope,
    Users,
    ClipboardList
  } from "lucide-svelte";

  export let isOpen: boolean;
  export let role: "PATIENT" | "DOCTOR" | "ADMIN" = "PATIENT";
  export let className = "";

  const navItems = {
    PATIENT: [
      { title: "Dashboard", icon: Home, url: "#" },
      { title: "My Appointments", icon: Calendar, url: "#" },
      { title: "Medical Records", icon: FileText, url: "#" },
      { title: "Prescriptions", icon: Pill, url: "#" },
      { title: "Health Monitoring", icon: Heart, url: "#" },
    ],
    DOCTOR: [
      { title: "Dashboard", icon: Home, url: "#" },
      { title: "Appointments", icon: Calendar, url: "#" },
      { title: "Patient List", icon: Users, url: "#" },
      { title: "Tasks", icon: ClipboardList, url: "#" },
      { title: "Medical Reports", icon: FileText, url: "#" },
    ],
    ADMIN: [
      { title: "Dashboard", icon: Home, url: "#" },
      { title: "User Management", icon: Users, url: "#" },
      { title: "Reports", icon: FileText, url: "#" },
      { title: "System Health", icon: Activity, url: "#" },
    ],
  };

  const quickActions = {
    PATIENT: [
      { title: "Contact Doctor", icon: Phone, url: "#" },
      { title: "Update Profile", icon: User, url: "#" },
      { title: "View Health Records", icon: Stethoscope, url: "#" },
    
    ],
    DOCTOR: [
      { title: "Add Patient", icon: FileText, url: "#" },
      { title: "Schedule Appointment", icon: User, url: "#" },
      { title: "View Reports", icon: Stethoscope, url: "#" },
      { title: "Manage Tasks", icon: ClipboardList, url: "#" },
      { title: "Patient Records", icon: Users, url: "#" },
      
    ],
    ADMIN: [
      { title: "Add User", icon: User, url: "#" },
      { title: "System Settings", icon: Stethoscope, url: "#" },
      { title: "View Logs", icon: FileText, url: "#" },
    ],
  };
</script>
<div class={className}>
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
      <div class="text-xs text-gray-500 uppercase mb-2 px-2" class:hidden={!isOpen}>
        Navigation
      </div>
      {#each navItems[role] as item}
        <a
          href={item.url}
          class="flex items-center gap-3 px-3 py-2 rounded hover:bg-gray-100 transition"
        >
          <item.icon class="w-5 h-5" />
          {#if isOpen}<span>{item.title}</span>{/if}
        </a>
      {/each}
    </div>

    <div>
      <div class="text-xs text-gray-500 uppercase mb-2 px-2" class:hidden={!isOpen}>
        Quick Actions
      </div>
      {#each quickActions[role] as item}
        <a
          href={item.url}
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
          <div class="text-xs text-gray-500">
            {role === 'DOCTOR' ? 'Doctor ID: #D456' : 'Patient ID: #12345'}
          </div>
        </div>
      {/if}
    </div>
  </div>
</aside>
</div>

<style>
  :global(.collapsed) {
    width: 64px !important;
  }
  aside {
    width: 250px;
  }
</style>
