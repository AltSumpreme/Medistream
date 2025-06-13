<script lang="ts">
  export let variant: 'default' | 'outline' = 'default';
  export let size: 'sm' | 'md' = 'md';
  export let className = '';
  export let type: 'button' | 'submit' | 'reset' = 'button';
</script>

<button
  type={type}
  class={`inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2
    ${variant === 'default' ? 'bg-blue-600 text-white hover:bg-blue-700' : ''}
    ${variant === 'outline' ? 'border border-gray-300 text-gray-800 hover:bg-gray-50' : ''}
    ${size === 'sm' ? 'px-3 py-1.5' : 'px-4 py-2'}
    ${className}`}
>
  <slot />
</button>
