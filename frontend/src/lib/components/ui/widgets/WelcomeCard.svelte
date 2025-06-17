<script lang="ts">
  export let name: string = "John";
  export let appointmentsCount: number = 2;
  export let prescriptionsCount: number = 3;
</script>

  <div class="grid auto-rows-min gap-4 md:grid-cols-3">
    <div class="aspect-video rounded-xl bg-gradient-to-br from-green-50 to-green-100 flex items-center justify-center">
      <div class="text-center">
        <h2 class="text-2xl font-bold text-green-700">Welcome Back,{name}!</h2>
        <p class="text-green-600">Your health is our priority</p>
      </div>
    </div>
    <div class="aspect-video rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 flex items-center justify-center">
      <div class="text-center">
        <div class="text-3xl font-bold text-blue-700">{appointmentsCount}</div>
        <p class="text-sm text-blue-600">Upcoming Appointments</p>
      </div>
    </div>
    <div class="aspect-video rounded-xl bg-gradient-to-br from-orange-50 to-orange-100 flex items-center justify-center">
      <div class="text-center">
        <div class="text-3xl font-bold text-orange-700">{prescriptionsCount}</div>
        <p class="text-sm text-orange-600">Active Prescriptions</p>
      </div>
    </div>
  </div>
