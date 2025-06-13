<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import  Badge  from "$lib/components/ui/Badge.svelte";
  import  Button  from "$lib/components/ui/Button.svelte";
  import  Progress  from "$lib/components/ui/Progress.svelte";
  import { Calendar, Heart, Pill, FileText, Clock, AlertCircle } from "lucide-svelte";
</script>

<div class="flex flex-1 flex-col gap-4 p-4 pt-0">
  <!-- Welcome Section -->
  <div class="grid auto-rows-min gap-4 md:grid-cols-3">
    <div class="aspect-video rounded-xl bg-gradient-to-br from-green-50 to-green-100 flex items-center justify-center">
      <div class="text-center">
        <h2 class="text-2xl font-bold text-green-700">Welcome Back, John!</h2>
        <p class="text-green-600">Your health is our priority</p>
      </div>
    </div>
    <div class="aspect-video rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 flex items-center justify-center">
      <div class="text-center">
        <div class="text-3xl font-bold text-blue-700">2</div>
        <p class="text-sm text-blue-600">Upcoming Appointments</p>
      </div>
    </div>
    <div class="aspect-video rounded-xl bg-gradient-to-br from-orange-50 to-orange-100 flex items-center justify-center">
      <div class="text-center">
        <div class="text-3xl font-bold text-orange-700">3</div>
        <p class="text-sm text-orange-600">Active Prescriptions</p>
      </div>
    </div>
  </div>

  <!-- Health Overview -->
  <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
    <Card title="Blood Pressure" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <Heart class="h-4 w-4 text-red-500" />
      <div class="text-2xl font-bold">120/80</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">Normal</span> • Last checked today
      </p>
    </Card>

    <Card title="Heart Rate" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <Heart class="h-4 w-4 text-pink-500" />
      <div class="text-2xl font-bold">72 BPM</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">Normal</span> • Resting rate
      </p>
    </Card>

    <Card title="Weight" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <FileText class="h-4 w-4 text-blue-500" />
      <div class="text-2xl font-bold">75 kg</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">-2kg</span> from last month
      </p>
    </Card>

    <Card title="BMI" headerClass="flex justify-between items-center" bodyClass="space-y-1">
      <FileText class="h-4 w-4 text-purple-500" />
      <div class="text-2xl font-bold">22.5</div>
      <p class="text-xs text-muted-foreground">
        <span class="text-green-600">Normal</span> • Healthy range
      </p>
    </Card>
  </div>

  <!-- Main Content Grid -->
  <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
    <!-- Appointments -->
    <div class="col-span-4">
      <Card title="Upcoming Appointments" headerClass="flex items-center gap-2">
        <Calendar class="h-5 w-5" />
        <div class="space-y-4 mt-4">
          <div class="flex justify-between items-center p-3 border rounded-lg">
            <div class="flex items-center gap-4">
              <div class="flex flex-col items-center">
                <Clock class="h-4 w-4 text-muted-foreground" />
                <span class="text-sm font-medium">10:00 AM</span>
              </div>
              <div>
                <p class="text-sm font-medium">Dr. Sarah Wilson</p>
                <p class="text-sm text-muted-foreground">Cardiology Consultation</p>
                <p class="text-xs text-muted-foreground">Tomorrow, Jan 16</p>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <Badge variant="secondary">Confirmed</Badge>
              <Button variant="outline" size="sm">Reschedule</Button>
            </div>
          </div>

          <div class="flex justify-between items-center p-3 border rounded-lg">
            <div class="flex items-center gap-4">
              <div class="flex flex-col items-center">
                <Clock class="h-4 w-4 text-muted-foreground" />
                <span class="text-sm font-medium">2:30 PM</span>
              </div>
              <div>
                <p class="text-sm font-medium">Dr. Michael Brown</p>
                <p class="text-sm text-muted-foreground">Follow-up Check</p>
                <p class="text-xs text-muted-foreground">Jan 18, 2024</p>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <Badge variant="outline">Pending</Badge>
              <Button variant="outline" size="sm">Confirm</Button>
            </div>
          </div>
        </div>
      </Card>
    </div>

    <!-- Health Goals -->
    <div class="col-span-3">
      <Card title="Health Goals" headerClass="flex items-center gap-2">
        <Heart class="h-5 w-5" />
        <div class="space-y-4 mt-4">
          <div>
            <div class="flex justify-between text-sm mb-2">
              <span>Daily Steps</span>
              <span>7,500 / 10,000</span>
            </div>
            <Progress value={75} className="h-2" />
          </div>

          <div>
            <div class="flex justify-between text-sm mb-2">
              <span>Water Intake</span>
              <span>6 / 8 glasses</span>
            </div>
            <Progress value={75} className="h-2" />
          </div>

          <div>
            <div class="flex justify-between text-sm mb-2">
              <span>Sleep Hours</span>
              <span>7 / 8 hours</span>
            </div>
            <Progress value={87.5} className="h-2" />
          </div>

          <div>
            <div class="flex justify-between text-sm mb-2">
              <span>Exercise Minutes</span>
              <span>25 / 30 min</span>
            </div>
            <Progress value={83} className="h-2" />
          </div>
        </div>
      </Card>
    </div>
  </div>

  <!-- Medications & Alerts -->
  <div class="grid gap-4 md:grid-cols-2">
    <Card title="Current Medications" headerClass="flex items-center gap-2">
      <Pill class="h-5 w-5" />
      <div class="space-y-3 mt-4">
        <div class="flex justify-between p-2 border rounded">
          <div>
            <p class="font-medium">Lisinopril 10mg</p>
            <p class="text-sm text-muted-foreground">Once daily, morning</p>
          </div>
          <Badge variant="secondary">Active</Badge>
        </div>

        <div class="flex justify-between p-2 border rounded">
          <div>
            <p class="font-medium">Metformin 500mg</p>
            <p class="text-sm text-muted-foreground">Twice daily, with meals</p>
          </div>
          <Badge variant="secondary">Active</Badge>
        </div>

        <div class="flex justify-between p-2 border rounded">
          <div>
            <p class="font-medium">Vitamin D3</p>
            <p class="text-sm text-muted-foreground">Once daily</p>
          </div>
          <Badge variant="outline">Supplement</Badge>
        </div>
      </div>
    </Card>

    <Card title="Health Alerts" headerClass="flex items-center gap-2">
      <AlertCircle class="h-5 w-5" />
      <div class="space-y-3 mt-4">
        <div class="flex items-start gap-3 p-2 bg-yellow-50 border border-yellow-200 rounded">
          <AlertCircle class="h-4 w-4 text-yellow-600 mt-0.5" />
          <div>
            <p class="text-sm font-medium">Prescription Renewal Due</p>
            <p class="text-xs text-muted-foreground">Lisinopril expires in 3 days</p>
          </div>
        </div>

        <div class="flex items-start gap-3 p-2 bg-blue-50 border border-blue-200 rounded">
          <AlertCircle class="h-4 w-4 text-blue-600 mt-0.5" />
          <div>
            <p class="text-sm font-medium">Lab Results Available</p>
            <p class="text-xs text-muted-foreground">Blood work from Jan 10</p>
          </div>
        </div>

        <div class="flex items-start gap-3 p-2 bg-green-50 border border-green-200 rounded">
          <Heart class="h-4 w-4 text-green-600 mt-0.5" />
          <div>
            <p class="text-sm font-medium">Great Progress!</p>
            <p class="text-xs text-muted-foreground">You've met your step goal 5 days in a row</p>
          </div>
        </div>
      </div>
    </Card>
  </div>
</div>
