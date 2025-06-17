<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import { Calendar, Clock } from "lucide-svelte";

  type Appointment = {
    time: string;
    doctor: string;
    type: string;
    date: string;
    status: "Confirmed" | "Pending" | "Urgent";
    action: "Reschedule" | "Confirm";
  };

  export let appointments: Appointment[] = [
    {
      time: "10:00 AM",
      doctor: "Dr. Sarah Wilson",
      type: "Cardiology Consultation",
      date: "Tomorrow, Jan 16",
      status: "Confirmed",
      action: "Reschedule"
    },
    {
      time: "2:30 PM",
      doctor: "Dr. Michael Brown",
      type: "Follow-up Check",
      date: "Jan 18, 2024",
      status: "Pending",
      action: "Confirm"
    },
    {
      time: "5:00 PM",
      doctor: "Dr. Emily Stone",
      type: "Emergency Checkup",
      date: "Today, Jan 15",
      status: "Urgent",
      action: "Confirm"
    }
  ];

  function getBadgeVariant(status: Appointment["status"]): "green" | "blue" | "red" | "yellow" {
    if (status === "Confirmed") return "green";
    if (status === "Pending") return "yellow";
    if (status === "Urgent") return "red";
    return "blue"; 
  }
</script>

<Card title="Upcoming Appointments" headerClass="flex items-center gap-2">
  <Calendar class="h-5 w-5" />
  <div class="space-y-4 mt-4">
    {#each appointments as appt}
      <div class="flex justify-between items-center p-3 border rounded-lg">
        <div class="flex items-center gap-4">
          <div class="flex flex-col items-center">
            <Clock class="h-4 w-4 text-muted-foreground" />
            <span class="text-sm font-medium">{appt.time}</span>
          </div>
          <div>
            <p class="text-sm font-medium">{appt.doctor}</p>
            <p class="text-sm text-muted-foreground">{appt.type}</p>
            <p class="text-xs text-muted-foreground">{appt.date}</p>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <Badge variant={getBadgeVariant(appt.status)}>{appt.status}</Badge>
          <Button variant="outline" size="sm">{appt.action}</Button>
        </div>
      </div>
    {/each}
  </div>
</Card>
