<script lang="ts">
  export let value: number = 0;
  export let className = '';
</script>

<div class={`w-full bg-gray-200 rounded-full overflow-hidden h-2 ${className}`}>
  <div
    class="bg-green-500 h-full transition-all duration-300 ease-in-out"
    style={`width: ${Math.min(value, 100)}%`}
  ></div>
</div>
