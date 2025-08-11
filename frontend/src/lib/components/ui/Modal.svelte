<script lang="ts">

import {X} from 'lucide-svelte';

export let title: string;
export let onClose: () => void;

function handleBackdropClick(event: MouseEvent){
    if(event.target === event.currentTarget) {
        onClose();
    }
}
 function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Escape') onClose();
  }

</script>



<!-- Modal Backdrop -->
<div
  class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4"
  on:click={handleBackdropClick}
  on:keydown={handleKeyDown}
  tabindex="0"
  role="dialog"
  aria-modal="true"
  aria-labelledby="modal-title"
>
  <!-- Modal Content -->
  <section class="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto animate-fadeIn"
       role="document"
       aria-labelledby="modal-title">
    <!-- Header -->
    <div class="flex items-center justify-between p-6 border-b border-gray-200">
      <h2 id="modal-title" class="text-xl font-semibold text-gray-900">{title}</h2>
      <button
        on:click={onClose}
        class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
      >
        <X class="size-5" />
      </button>
    </div>

    <!-- Body -->
    <div class="p-6 space-y-6">
      <slot />
    </div>
  </section>

   <div class="p-4 border-t border-gray-200">
      <slot name="footer" />
    </div>
</div>



<style>

    .animate-fadeIn{
        animation: fadeIn 0.3s ease-in-out;
    }
    @keyframes fadeIn {
    from { opacity: 0; transform: scale(0.98); }
    to { opacity: 1; transform: scale(1); }
  }

</style>
