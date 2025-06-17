<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import { Pill } from "lucide-svelte";

  export let medications: {
    name: string;
    instructions: string;
    status: "Active" | "Pending" | "Urgent" | "Supplement";
    variant: string;
  }[] = [
    { name: "Lisinopril 10mg", instructions: "Once daily, morning", status: "Active", variant: "green" },
    { name: "Metformin 500mg", instructions: "Twice daily, with meals", status: "Active", variant: "green" },
    { name: "Vitamin D3", instructions: "Once daily", status: "Supplement", variant: "blue" }
  ];

   function getBadgeVariant(status: "Active" | "Pending" | "Urgent" | "Supplement"): "green" | "blue" | "red" | "yellow" {
    if (status === "Active") return "green";
    if (status === "Supplement") return "yellow";
    if (status === "Urgent") return "red";
    return "blue"; 
  }
</script>

<Card title="Current Medications" headerClass="flex items-center gap-2">
  <Pill class="h-5 w-5" />
  <div class="space-y-3 mt-4">
    {#each medications as med}
      <div class="flex justify-between p-2 border rounded">
        <div>
          <p class="font-medium">{med.name}</p>
          <p class="text-sm text-muted-foreground">{med.instructions}</p>
        </div>
        <Badge variant={getBadgeVariant(med.status)}>{med.status}</Badge>
      </div>
    {/each}
  </div>
</Card>
