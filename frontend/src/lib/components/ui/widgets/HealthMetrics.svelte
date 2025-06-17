<script lang="ts">

import Card from '$lib/components/ui/Card.svelte';
import { Heart, FileText } from 'lucide-svelte';

export  let BloodPressure = '120/80';
export  let HeartRate = '72 BPM';
export  let Weight = '75 kg';
export  let BMI = '22.5';
</script>

<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
    <Card title="Blood Pressure" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <Heart class="h-4 w-4 text-red-500" />
      <div class="text-2xl font-bold">{BloodPressure}</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">Normal</span> • Last checked today
      </p>
    </Card>

    <Card title="Heart Rate" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <Heart class="h-4 w-4 text-pink-500" />
      <div class="text-2xl font-bold">{HeartRate}</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">Normal</span> • Resting rate
      </p>
    </Card>

    <Card title="Weight" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <FileText class="h-4 w-4 text-blue-500" />
      <div class="text-2xl font-bold">{Weight}</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">-2kg</span> from last month
      </p>
    </Card>

    <Card title="BMI" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <FileText class="h-4 w-4 text-purple-500" />
      <div class="text-2xl font-bold">{BMI}</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">Normal</span> • Healthy range
      </p>
    </Card>
  </div>