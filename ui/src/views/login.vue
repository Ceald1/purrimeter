<script setup lang="ts">
import { ref } from "vue"
import { setCookie } from 'typescript-cookie';
import LoadingSpinner from "../components/LoadingSpinner.vue"; // import it

const username = ref('')
const password = ref('')
const isLoading = ref(false)
const errors = ref<{ username?: string; password?: string }>({})
const showPassword = ref<boolean>(false)

const validateForm = () => {
  errors.value = {}

  if (!username.value) {
    errors.value.username = 'Username is required!'
  }

  if (!password.value) {
    errors.value.password = 'Password is required!'
  }

  return Object.keys(errors.value).length === 0
}

const handleLogin = async () => {
  if (!validateForm()) return
  
  isLoading.value = true

  try {
    const response = await fetch('/purr/api/v1/user/login', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        username: username.value,
        password: password.value
      })
    })
    if (response.ok) {
      const data = await response.json()
      setCookie("token", JSON.parse(data)["token"]) // jwt will be stored in token field
    }

  }catch (error) {
    console.error("login failed!", error)
  }finally{
    isLoading.value = false
  }
}
</script>

<template>
<div class="w-full min-h-screen flex items-center justify-center bg-[length:200px_115.4px] bg-[repeating-conic-gradient(from_30deg_at_calc(50%+100px)_calc(50%+57.7px),transparent_0_120deg,#313244_0_180deg),repeating-conic-gradient(from_30deg,#1e1e2e_0_60deg,#45475a_0_120deg,#313244_0_180deg)] text-[#cdd6f4]">
    <form 
      @submit.prevent="handleLogin"
      class="flex flex-col gap-2.5 px-8 pb-6 pt-2 rounded-[25px] bg-ctp-base transition-all duration-[400ms] ease-in-out hover:scale-105 hover:border hover:border-ctp-crust w-full max-w-md"
    >
      <p class="text-center my-8 text-ctp-text text-xl font-medium">Login</p>
      
      <!-- Username Field -->
      <div class="flex items-center justify-center gap-2 rounded-[25px] px-2.5 py-2.5 bg-ctp-mantle shadow-[inset_2px_5px_10px_rgba(5,5,5,0.5)]">
        <svg class="h-5 w-5 fill-ctp-text flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
          <path d="M13.106 7.222c0-2.967-2.249-5.032-5.482-5.032-3.35 0-5.646 2.318-5.646 5.702 0 3.493 2.235 5.708 5.762 5.708.862 0 1.689-.123 2.304-.335v-.862c-.43.199-1.354.328-2.29.328-2.926 0-4.813-1.88-4.813-4.798 0-2.844 1.921-4.881 4.594-4.881 2.735 0 4.608 1.688 4.608 4.156 0 1.682-.554 2.769-1.416 2.769-.492 0-.772-.28-.772-.76V5.206H8.923v.834h-.11c-.266-.595-.881-.964-1.6-.964-1.4 0-2.378 1.162-2.378 2.823 0 1.737.957 2.906 2.379 2.906.8 0 1.415-.39 1.709-1.087h.11c.081.67.703 1.148 1.503 1.148 1.572 0 2.57-1.415 2.57-3.643zm-7.177.704c0-1.197.54-1.907 1.456-1.907.93 0 1.524.738 1.524 1.907S8.308 9.84 7.371 9.84c-.895 0-1.442-.725-1.442-1.914z"></path>
        </svg>
        <input 
          v-model="username"
          autocomplete="off" 
          placeholder="Username" 
          class="bg-transparent border-none outline-none w-full text-ctp-subtext0 placeholder:text-ctp-overlay0" 
          type="text"
          :disabled="isLoading"
        >
      </div>
      <p v-if="errors.username" class="text-ctp-red text-sm ml-2 -mt-1">{{ errors.username }}</p>
      
      <!-- Password Field -->
      <div class="flex items-center justify-center gap-2 rounded-[25px] px-2.5 py-2.5 bg-ctp-mantle shadow-[inset_2px_5px_10px_rgba(5,5,5,0.5)]">
        <svg class="h-5 w-5 fill-ctp-text flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 1a2 2 0 0 1 2 2v4H6V3a2 2 0 0 1 2-2zm3 6V3a3 3 0 0 0-6 0v4a2 2 0 0 0-2 2v5a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2z"></path>
        </svg>
        <input 
          v-model="password"
          placeholder="Password" 
          class="bg-transparent border-none outline-none w-full text-ctp-subtext0 placeholder:text-ctp-overlay0" 
          :type="showPassword ? 'text' : 'password'"
          :disabled="isLoading"
        >
        <button 
          type="button"
          @click="showPassword = !showPassword"
          class="text-ctp-overlay0 hover:text-ctp-text transition-colors flex-shrink-0"
        >
          <svg v-if="!showPassword" class="h-5 w-5" fill="currentColor" viewBox="0 0 16 16">
            <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8zM1.173 8a13.133 13.133 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.133 13.133 0 0 1 14.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5c-2.12 0-3.879-1.168-5.168-2.457A13.134 13.134 0 0 1 1.172 8z"/>
            <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5zM4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0z"/>
          </svg>
          <svg v-else class="h-5 w-5" fill="currentColor" viewBox="0 0 16 16">
            <path d="M13.359 11.238C15.06 9.72 16 8 16 8s-3-5.5-8-5.5a7.028 7.028 0 0 0-2.79.588l.77.771A5.944 5.944 0 0 1 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.134 13.134 0 0 1 14.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755-.165.165-.337.328-.517.486l.708.709z"/>
            <path d="M11.297 9.176a3.5 3.5 0 0 0-4.474-4.474l.823.823a2.5 2.5 0 0 1 2.829 2.829l.822.822zm-2.943 1.299.822.822a3.5 3.5 0 0 1-4.474-4.474l.823.823a2.5 2.5 0 0 0 2.829 2.829z"/>
            <path d="M3.35 5.47c-.18.16-.353.322-.518.487A13.134 13.134 0 0 0 1.172 8l.195.288c.335.48.83 1.12 1.465 1.755C4.121 11.332 5.881 12.5 8 12.5c.716 0 1.39-.133 2.02-.36l.77.772A7.029 7.029 0 0 1 8 13.5C3 13.5 0 8 0 8s.939-1.721 2.641-3.238l.708.709zm10.296 8.884-12-12 .708-.708 12 12-.708.708z"/>
          </svg>
        </button>
      </div>
      <p v-if="errors.password" class="text-ctp-red text-sm ml-2 -mt-1">{{ errors.password }}</p>
      
      <!-- Submit Button -->
       <div class="flex justify-center mt-10 mb-4">
  <template v-if="isLoading">
    <LoadingSpinner />
  </template>
  <template v-else>
    <button 
      type="submit"
      class="px-6 py-2 rounded-md transition-all duration-[400ms] ease-in-out bg-ctp-surface0 text-ctp-text hover:bg-ctp-surface1 font-medium"
    >
      Login
    </button>
  </template>
      <!-- <div class="flex justify-center mt-10 mb-4">
        <button 
          type="submit"
          :disabled="isLoading"
          class="px-6 py-2 rounded-md transition-all duration-[400ms] ease-in-out bg-ctp-surface0 text-ctp-text hover:bg-ctp-surface1 disabled:opacity-50 disabled:cursor-not-allowed font-medium"
        > -->
          <!-- {{ isLoading ? 'Loading...' : 'Login' }} -->
        <!-- </button> -->
      </div>
    </form>
  </div>
</template>