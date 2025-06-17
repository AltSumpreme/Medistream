<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import Progress from "$lib/components/ui/Progress.svelte";
  import { Heart } from "lucide-svelte";

  export let goals = [
    { label: "Daily Steps", current: 7500, total: 10000 },
    { label: "Water Intake", current: 6, total: 8, unit: "glasses" },
    { label: "Sleep Hours", current: 7, total: 8, unit: "hours" },
    { label: "Exercise Minutes", current: 25, total: 30, unit: "min" }
  ];

  const getPercentage = (current: number, total: number) =>
    Math.min(100, (current / total) * 100);
</script>

<Card title="Health Goals" headerClass="flex items-center gap-2">
  <Heart class="h-5 w-5" />
  <div class="space-y-4 mt-4">
    {#each goals as { label, current, total, unit = '' }}
      <div>
        <div class="flex justify-between text-sm mb-2">
          <span>{label}</span>
          <span>{current} / {total} {unit}</span>
        </div>
        <Progress value={getPercentage(current, total)} className="h-2" />
      </div>
    {/each}
  </div>
</Card>
