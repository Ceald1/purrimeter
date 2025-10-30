import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAuth = function Auth():boolean {
    const isAuthenticated = false
    const validate = async (): Promise<boolean> => {
        
        
        return false
    }
    

    return isAuthenticated
}