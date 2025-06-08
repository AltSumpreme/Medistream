<script lang="ts">
  import Card from '$lib/components/Card.svelte';
  import { Activity } from 'lucide-svelte';
  import {registerUser} from '$lib/services/auth';

  let firstName = '';
  let lastName = '';
  let email = '';
  let password = '';
  let confirmPassword = '';
  let errorMessage = '';

  $: errorMessage = confirmPassword && confirmPassword !== password ? 'Passwords do not match' : '';
  
 const handleSignup = async () => {
  if (!firstName || !lastName || !email || !password || !confirmPassword) {
    errorMessage = 'All fields are required';
    return;
  }
   try {
    await registerUser(firstName, lastName, email, password);
    window.location.href = '/login';
  } catch (error) {
    errorMessage = 'Signup failed. Please try again.';
    console.error('Signup failed:', error);
  }
}

    $: {
        if (confirmPassword && confirmPassword !== password) {                  
            errorMessage = 'Passwords do not match';
        } else {
            errorMessage = '';
        }
    }
</script>

<div class="min-h-screen flex flex-col items-center justify-center px-4 py-8 bg-gradient-to-br from-blue-50 to-indigo-100">
  <div class="flex items-center justify-center mb-4">
    <div class="bg-blue-600 p-3 rounded-full">
      <Activity class="h-8 w-8 text-white" />
    </div>
  </div>

  <h1 class="text-center text-black text-3xl font-bold">MediStream</h1>
  <p class="text-center text-gray-600 mt-2">Healthcare management platform</p>

  <Card
    title="Create Account"
    description="Join us to manage your healthcare needs"
    className="mt-6 w-full max-w-md"
    headerClass="text-center"
    bodyClass="space-y-4"
  >
    <form on:submit|preventDefault={handleSignup} class="space-y-4">
      <div>
        <label for="firstName" class="block text-sm font-medium text-gray-700">First Name</label>
        <input
          id="firstName"
          type="text"
          bind:value={firstName}
          required
          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50"
        />
      </div>

      <div>
        <label for="lastName" class="block text-sm font-medium text-gray-700">Last Name</label>
        <input
          id="lastName"
          type="text"
          bind:value={lastName}
          required
          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50"
        />
      </div>

      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email Address</label>
        <input
          id="email"
          type="email"
          bind:value={email}
          required
          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-blue-500 focus:ring focus:ring-gray-500 focus:ring-opacity-50"
        />
      </div>

      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          required
          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50"
        />
      </div>

      <div>
        <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirm Password</label>
        <input
          id="confirmPassword"
          type="password"
          bind:value={confirmPassword}
          required
          class={`mt-1 block w-full border rounded-md shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 ${
            confirmPassword && confirmPassword !== password ? 'border-gray-300' : 'border-red-500'
          }`}
        />
        {#if errorMessage}
          <p class="text-red-500 text-sm mt-1">{errorMessage}</p>
        {/if}
      </div>

      <div>
        <button type="submit" class="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 disabled:opacity-50" disabled={!!errorMessage}>
          Sign Up
        </button>
      </div>
    </form>
  </Card>

  <div class="mt-6 text-center text-sm text-gray-500">
    <p>Your data is secure with us. We prioritize your privacy and security.</p>
    <p class="mt-1">HIPAA Compliant • SOC 2 Certified</p>
  </div>
</div>
