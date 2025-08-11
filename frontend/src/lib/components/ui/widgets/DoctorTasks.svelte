<script lang="ts">
  import { writable } from 'svelte/store';
  import { ClipboardList, Plus, Clock, AlertCircle, Trash } from 'lucide-svelte';

  // Define Task model
  interface Task {
    id: string;
    title: string;
    priority: 'High' | 'Medium' | 'Low';
    due: string;
    completed: boolean;
  }

  // Sample starter tasks (you can replace/add more)
  const tasks = writable<Task[]>([
    {
      id: crypto.randomUUID(),
      title: "Review patient lab results",
      priority: "High",
      due: "Today",
      completed: false
    },
    {
      id: crypto.randomUUID(),
      title: "Approve medical prescriptions",
      priority: "Medium",
      due: "Tomorrow",
      completed: false
    },
    {
      id: crypto.randomUUID(),
      title: "Prepare discharge summary",
      priority: "Low",
      due: "Friday",
      completed: true
    }
  ]);

  // Task actions
  function toggleTask(id: string) {
    tasks.update(current =>
      current.map(task =>
        task.id === id ? { ...task, completed: !task.completed } : task
      )
    );
  }

  function removeTask(id: string) {
    tasks.update(current => current.filter(task => task.id !== id));
  }

  function getPriorityColor(priority: string) {
    return {
      High: 'bg-red-100 text-red-800',
      Medium: 'bg-yellow-100 text-yellow-800',
      Low: 'bg-green-100 text-green-800',
    }[priority] ?? 'bg-gray-100 text-gray-800';
  }

  function getPriorityTextColor(priority: string) {
    return {
      High: 'text-red-500',
      Medium: 'text-orange-500',
      Low: 'text-green-500',
    }[priority] ?? 'text-gray-500';
  }
</script>

<div class="bg-white rounded-lg border border-gray-200">
  <div class="flex items-center justify-between p-6 border-b border-gray-200">
    <h2 class="flex items-center gap-2 text-lg font-semibold text-gray-900">
      <ClipboardList class="size-5" />
      Medical Tasks
    </h2>
    <button
      class="flex items-center gap-1 px-3 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
      on:click={() => {
        // Placeholder sample – change as needed
        tasks.update(current => [
          ...current,
          {
            id: crypto.randomUUID(),
            title: "New Follow-up Task",
            priority: "Low",
            due: "Next Week",
            completed: false
          }
        ]);
      }}
    >
      <Plus class="size-4" />
      Add Task
    </button>
  </div>

  <div class="p-6 space-y-3">
    {#each $tasks as task (task.id)}
      <div class="flex items-start gap-3 p-3 border border-gray-200 rounded-lg {task.completed ? 'bg-gray-50' : 'hover:bg-gray-20'} transition-colors">
        <input 
          type="checkbox" 
          checked={task.completed}
          on:change={() => toggleTask(task.id)}
          class="mt-1 size-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
          id={`task-${task.id}`}
        />
        <div class="flex-1">
          <label class="text-sm font-medium {task.completed ? 'line-through text-gray-500' : 'text-gray-900'}" for={`task-${task.id}`}>
            {task.title}
          </label>
          <div class="flex items-center gap-4 text-xs text-gray-500 mt-1">
            <div class="flex items-center gap-1">
              <Clock class="size-3" />
              <span>{task.due}</span>
            </div>
            <div class="flex items-center gap-1">
              <AlertCircle class="size-3" />
              <span class="{getPriorityTextColor(task.priority)}">{task.priority} Priority</span>
            </div>
          </div>
        </div>
        <span class="px-2 py-1 text-xs font-medium rounded-full {getPriorityColor(task.priority)}">
          {task.priority}
        </span>
        <button
          class="ml-2 text-gray-400 hover:text-red-500"
          on:click={() => removeTask(task.id)}
        >
          <Trash class="size-4" />
        </button>
      </div>
    {/each}
  </div>
</div>
