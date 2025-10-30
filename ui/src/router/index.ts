import { createRouter, createWebHistory } from 'vue-router'

// Import your page components
import LoginPage from './../views/login.vue'
import IndexPage from './../views/index.vue'
import DashboardPage from './../views/dashboard.vue' // Define routes
import Auth from './auth' // authentication checks
const routes = [
  {path: '/', name: 'index', component: IndexPage},
  { path: '/login', name: 'Login', component: LoginPage },
  { path: '/dashboard', name: 'Dashboard', component: DashboardPage },
//   {path: '/', name: 'index', component: App}
]

// Create router instance
const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
