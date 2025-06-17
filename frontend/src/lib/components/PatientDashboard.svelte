<script lang="ts">

  import WelcomeCard from "./ui/widgets/WelcomeCard.svelte";
  import HealthMetrics from "./ui/widgets/HealthMetrics.svelte";
  import UpcomingAppointments from "./ui/widgets/UpcomingAppointments.svelte";
  import MedicationCard from "./ui/widgets/MedicationCard.svelte";
  import HealthGoalsCard from "./ui/widgets/HealthGoalsCard.svelte";


  import { Calendar, Heart, Pill, FileText, Clock, AlertCircle } from "lucide-svelte";
</script>

<div class="flex flex-1 flex-col gap-4 p-4 pt-0">
  <WelcomeCard />
  <HealthMetrics />
  


  <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
    
    <div class="col-span-4">
      <UpcomingAppointments />
    </div>

    <!-- Health Goals -->
    <div class="col-span-3">
      <HealthGoalsCard />
    </div>
  </div>

  <!-- Medications & Alerts -->
  <div class="items-center">
    <MedicationCard />
    
  </div>
</div>
