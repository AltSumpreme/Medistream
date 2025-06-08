<script lang="ts">
  import Card from '$lib/components/Card.svelte';
  import { Activity } from 'lucide-svelte';
  import {loginUser} from '$lib/services/auth';
	import { goto } from '$app/navigation';

  let email = '';
  let password = '';
  let errorMessage = '';
  $: errorMessage = !email || !password ? 'Email and password are required' : '';
 
  const handleLogin = async () => {

    if (errorMessage) return;

    try{
      await loginUser(email, password);
      goto('/dashboard');
    } catch (error) {
      errorMessage = 'Login failed. Please check your credentials.';
      console.error('Login failed:', error);
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
    title="Welcome Back"
    description="Sign in to access your account"
    className="mt-6 w-full max-w-md"
    headerClass="text-center"
    bodyClass="space-y-4"
  >
    <form on:submit|preventDefault={handleLogin} class="space-y-4">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input
          id="email"
          type="email"
          name="email"
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
          name="password"
          bind:value={password}
          required
          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50"
        />
      </div>

      <div>
        <button type="submit" class="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600">
          Sign In
        </button>
      </div>
    </form>
  </Card>

 
  <div class="mt-6 text-center text-sm text-gray-500">
    <p>Protected by enterprise-grade security</p>
    <p class="mt-1">HIPAA Compliant • SOC 2 Certified</p>
  </div>
</div>
