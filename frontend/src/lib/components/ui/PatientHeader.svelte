<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Menu, Search } from "lucide-svelte";
  const dispatch = createEventDispatcher();
</script>

<header class="flex items-center justify-between h-16 px-4 bg-white border-b shadow-sm">
 
  <button on:click={() => dispatch("toggleSidebar")} class="p-2 rounded hover:bg-gray-100">
    <Menu class="h-5 w-5" />
  </button>  
  <div class="relative w-full max-w-sm mx-4">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
    <input
      type="text"
      placeholder="Search..."
      class="w-full pl-10 pr-4 py-2 rounded border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
    /> 
    
  </div>

 
  <div class="relative">
    <button
      on:click={() => window.location.href = '/login'}
      class="flex items-center gap-2 px-3 py-1.5 border border-gray-300 rounded hover:bg-gray-100"
    >
     
      <span class="text-sm">Logout</span>
    </button>
  </div>
</header>
