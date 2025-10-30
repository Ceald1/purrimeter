<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getCookie } from 'typescript-cookie'
import LoadingSpinner from './../components/LoadingSpinner.vue'

const router = useRouter()
const isLoading = ref(true)

onMounted(async () => {
  const prevRoute = router.currentRoute.value.fullPath
  console.log(prevRoute)
  const token = getCookie('token')

  if (!token) {
    await router.replace("/login")
  } else {
    try {
      const response = await fetch("/purr/api/v1/user/check", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token }),
      })

      if (response.ok) {
        if (!prevRoute){
          await router.replace("/dashboard")
        }else{
          await router.replace(prevRoute)
        }
      } else {
        await router.replace("/login")
      }
    } catch (error) {
      await router.replace("/login")
    }
  }
  isLoading.value = false
})
</script>

<template>
  <div id="app">
    <LoadingSpinner v-if="isLoading" />
    <router-view v-else />
  </div>
</template>