<script lang="ts">
  import {  Calendar, Clock, User } from 'lucide-svelte';
  import Modal from '$lib/components/ui/Modal.svelte';

  export let onClose: () => void;

 

  function handleSubmit(e: Event) {
    e.preventDefault();
    // your submit logic
  }
</script>

<Modal title="New Medical Record" {onClose}>
  <form on:submit={handleSubmit} class="space-y-4 p-6">
    <!-- Patient Info -->
    <div>
      <label class="mb-1 flex items-center gap-1 text-sm font-medium">
        <User class="w-4" /> Patient
      </label>
      <input type="text" required class="w-full rounded border p-2" placeholder="Enter patient name" />
    </div>

    <!-- Date -->
    <div>
      <label class="mb-1 flex items-center gap-1 text-sm font-medium">
        <Calendar class="w-4" /> Date
      </label>
      <input type="date" required class="w-full rounded border p-2" />
    </div>

    <!-- Time -->
    <div>
      <label class="mb-1 flex items-center gap-1 text-sm font-medium">
        <Clock class="w-4" /> Time
      </label>
      <input type="time" required class="w-full rounded border p-2" />
    </div>

    <!-- Type -->
    <div>
      <label for="type-select" class="mb-1 flex text-sm font-medium">Type</label>
      <select id="type-select" class="w-full rounded border p-2">
        <option value="consultation">Consultation</option>
        <option value="follow-up">Follow-up</option>
        <option value="check-up">Check-up</option>
        <option value="video-call">Video Call</option>
      </select>
    </div>

    <!-- Notes -->
    <div>
      <label for="notes-textarea" class="mb-1 flex text-sm font-medium">Notes</label>
      <textarea id="notes-textarea" rows="3" class="w-full resize-none rounded border p-2"></textarea>
    </div>

    <!-- Footer Buttons -->
    <div class="flex justify-end gap-2 pt-4 border-t">
      <button type="button" on:click={onClose} class="rounded border px-4 py-2">Cancel</button>
      <button type="submit" class="rounded bg-blue-600 px-4 py-2 text-white">Save Record</button>
    </div>
  </form>
</Modal>
