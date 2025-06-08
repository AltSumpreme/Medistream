<script lang="ts">
    export let title = '';
    export let description = '';
    export let className = '';
    export let headerClass = '';
    export let bodyClass = '';
</script>

<div class={`bg-white shadow-md rounded-2xl p-6 ${className}`}>
  {#if title}
    <div class={`mb-4 ${headerClass}`}>
      <h2 class="text-2xl font-bold text-center">{title}</h2>
      {#if description}
        <p class="text-gray-600 text-sm mt-1 text-center">{description}</p>
      {/if}
    </div>
  {/if}

  <div class={bodyClass}>
    <slot />
  </div>
</div>