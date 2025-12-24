<script setup>
import { ref, onMounted, computed } from 'vue';
import { modules, getModuleByRoute, isValidRoute } from './modules';
import ChatDialog from './components/ChatDialog.vue';

// Route management for different app sections
const currentRoute = ref(null);

onMounted(() => {
  // Check for hash-based routing
  const hash = window.location.hash.slice(1);
  if (hash && isValidRoute(hash)) {
    currentRoute.value = hash;
  } else {
    currentRoute.value = modules.length > 0 ? modules[0].route : null;
  }

  // Listen for hash changes
  window.addEventListener('hashchange', () => {
    const newHash = window.location.hash.slice(1);
    if (newHash && isValidRoute(newHash)) {
      currentRoute.value = newHash;
    } else {
      currentRoute.value = modules.length > 0 ? modules[0].route : null;
    }
  });
});

// Get current active module
const currentModule = computed(() => {
  return getModuleByRoute(currentRoute.value);
});

// Navigation function
const navigateTo = (route) => {
  window.location.hash = route;
  currentRoute.value = route;
};
</script>

<template>
  <div class="min-h-screen bg-[#0a0a0b] text-white selection:bg-indigo-500/30">
    <!-- Background Decoration -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] bg-indigo-500/10 blur-[120px] rounded-full"></div>
      <div class="absolute -bottom-[10%] -right-[10%] w-[40%] h-[40%] bg-purple-500/10 blur-[120px] rounded-full"></div>
      <div class="absolute top-[20%] right-[10%] w-[20%] h-[20%] bg-blue-500/5 blur-[100px] rounded-full"></div>
    </div>

    <!-- Navigation Menu -->
    <nav class="relative z-20 border-b border-white/5 bg-[#0a0a0b]/80 backdrop-blur-xl sticky top-0">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-center gap-2 py-4 flex-wrap">
          <button
            v-for="module in modules"
            :key="module.id"
            @click="navigateTo(module.route)"
            :class="[
              'flex items-center gap-2 px-6 py-2.5 rounded-xl font-semibold transition-all duration-300',
              currentRoute === module.route
                ? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-lg shadow-indigo-600/20'
                : 'bg-white/5 text-gray-400 hover:bg-white/10 hover:text-white'
            ]"
          >
            <component :is="module.icon" class="w-5 h-5" />
            <span>{{ module.name }}</span>
          </button>
        </div>
      </div>
    </nav>

    <main class="relative z-10">
      <!-- Dynamic Module Rendering -->
      <component 
        v-if="currentModule" 
        :is="currentModule.component" 
        :key="currentModule.id"
      />
      
      <!-- Fallback if no module is found -->
      <div v-else class="container mx-auto px-4 py-20 text-center">
        <div class="max-w-lg mx-auto space-y-6">
          <div class="text-6xl">🔧</div>
          <h1 class="text-3xl font-bold text-white">Welcome to Kwun Tools</h1>
          <p class="text-gray-400">
            Select a tool from the navigation above to get started.
          </p>
        </div>
      </div>
    </main>

    <!-- Global AI Chat Assistant -->
    <ChatDialog />
  </div>
</template>

<style>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* Base resets in case style.css isn't fully loaded yet or overridden */
body {
  margin: 0;
  padding: 0;
  overflow-x: hidden;
}
</style>
