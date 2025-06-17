<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import { AlertCircle, Heart } from "lucide-svelte";

  export let alerts = [
    {
      title: "Prescription Renewal Due",
      description: "Lisinopril expires in 3 days",
      color: "yellow",
      icon: AlertCircle
    },
    {
      title: "Lab Results Available",
      description: "Blood work from Jan 10",
      color: "blue",
      icon: AlertCircle
    },
    {
      title: "Great Progress!",
      description: "You've met your step goal 5 days in a row",
      color: "green",
      icon: Heart
    }
  ];
</script>

<Card title="Health Alerts" headerClass="flex items-center gap-2">
  <AlertCircle class="h-5 w-5" />
  <div class="space-y-3 mt-4">
    {#each alerts as { title, description, color, icon: Icon }}
      <div class={`flex items-start gap-3 p-2 bg-${color}-50 border border-${color}-200 rounded`}>
        <Icon class={`h-4 w-4 text-${color}-600 mt-0.5`} />
        <div>
          <p class="text-sm font-medium">{title}</p>
          <p class="text-xs text-muted-foreground">{description}</p>
        </div>
      </div>
    {/each}
  </div>
</Card>
