<script lang="ts">
  export let variant: 'default' | 'secondary' | 'outline' = 'default';
  export let className = '';
</script>

<span
  class={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium
    ${variant === 'default' ? 'bg-gray-100 text-gray-800' : ''}
    ${variant === 'secondary' ? 'bg-blue-100 text-blue-800' : ''}
    ${variant === 'outline' ? 'border border-gray-300 text-gray-800' : ''}
    ${className}`}
>
  <slot />
</span>
