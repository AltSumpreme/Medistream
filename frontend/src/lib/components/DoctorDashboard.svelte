    <script lang="ts">
    import DoctorsSchedule from "./ui/widgets/DoctorsSchedule.svelte";
    import WelcomeCard from "./ui/widgets/WelcomeCard.svelte";
    import DoctorTasks from "./ui/widgets/DoctorTasks.svelte";
    import PatientList from "./ui/widgets/PatientList.svelte";


    type AppointmentStatus = "Upcoming" | "In Progress" | "Completed";

    type Appointment = {
    id: string;
    time: string;
    patient: string;
    patientId: string;
    type: string;
    duration: string;
    room: string;
    status: AppointmentStatus;
    avatar?: string;
    notes: string;
    };

    let appointments: Appointment[] = [
    {
        id: '1',
        time: '10:00 AM',
        patient: 'John Doe',
        patientId: 'p1',
        type: 'In-Person',
        duration: '30m',
        room: 'Room 101',
        status: "Upcoming",
        avatar: undefined,
        notes: 'First appointment',
    },
    {
        id: '2',
        time: '11:00 AM',
        patient: 'Jane Smith',
        patientId: 'p2',
        type: 'Video Consultation',
        duration: '30m',
        room: 'Online',
        status: "In Progress",
        avatar: undefined,
        notes: 'Follow-up',
    },
    ];
    

    
    </script>

    <div class="flex flex-1 flex-col gap-4 p-4 pt-0">
    <WelcomeCard name="Dr John Doe" appointmentsCount={2} prescriptionsCount={5} />

    <DoctorsSchedule appointments={appointments} />
    <DoctorTasks />
    <PatientList/>

    </div>

    

