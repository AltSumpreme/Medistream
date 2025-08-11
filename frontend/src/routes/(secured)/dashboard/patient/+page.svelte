<script lang="ts">
  import Sidebar from '$lib/components/ui/SideBar.svelte';
  import Header from '$lib/components/ui/Header.svelte';
  import PatientDashboard from '$lib/components/PatientDashboard.svelte';
  

  let isSidebarOpen = true;
  const toggleSidebar = () => isSidebarOpen = !isSidebarOpen;
</script>


<div class="flex flex-col h-screen">

  <header class="shrink-0">
    <Header on:toggleSidebar={toggleSidebar} />
  </header>

  
  <div class="flex flex-1 overflow-hidden">
    
    <div class="h-full shrink-0">
      <Sidebar isOpen={isSidebarOpen} />
    </div>

    
    <main class="flex-1 overflow-y-auto bg-gray-50">

      <PatientDashboard />
    </main>
  </div>
</div>
