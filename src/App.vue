<script setup>
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import ProfileView from './components/ProfileView.vue';
import ProfileEditor from './components/ProfileEditor.vue';
import { decodeProfile } from './utils/onelink';
import { Sparkles } from 'lucide-vue-next';

const { t } = useI18n();
const encodedData = ref(null);
const decodedProfile = ref(null);
const previewData = ref({
  name: t('common.previewName'),
  bio: t('common.previewBio'),
  links: [],
  socials: {}
});

onMounted(() => {
  const params = new URLSearchParams(window.location.search);
  const data = params.get('data');
  if (data) {
    encodedData.value = data;
    decodedProfile.value = decodeProfile(data);
  }
});

const isViewMode = computed(() => !!encodedData.value && !!decodedProfile.value);

const handleUpdate = (newProfile) => {
  previewData.value = { ...newProfile };
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

    <main class="relative z-10">
      <!-- View Mode -->
      <div v-if="isViewMode" class="animate-in fade-in duration-1000">
        <ProfileView :data="decodedProfile" />
      </div>

      <!-- Editor Mode -->
      <div v-else class="container mx-auto px-4 py-12 md:py-20 lg:py-24">
        <div class="flex flex-col lg:flex-row gap-16 items-start justify-center">
          <!-- Left: Editor -->
          <div class="w-full lg:w-1/2 space-y-12 animate-in slide-in-from-left duration-700">
            <div class="flex items-center gap-3">
              <div class="p-3 bg-indigo-600 rounded-2xl shadow-lg shadow-indigo-600/20">
                <Sparkles class="w-8 h-8 text-white" />
              </div>
              <h1 class="text-4xl font-extrabold tracking-tight bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
                {{ t('app.title') }}
              </h1>
            </div>
            
            <ProfileEditor @update="handleUpdate" />
          </div>

          <!-- Right: Preview -->
          <div class="w-full lg:w-1/3 lg:sticky lg:top-12 animate-in slide-in-from-right duration-700">
            <div class="relative">
              <div class="absolute -inset-4 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-[2.5rem] blur opacity-10"></div>
              <div class="relative bg-[#111113] border border-white/5 rounded-[2.5rem] overflow-hidden shadow-2xl">
                <div class="p-4 border-b border-white/5 bg-white/5 flex items-center justify-between">
                  <div class="flex gap-1.5">
                    <div class="w-2.5 h-2.5 rounded-full bg-red-500/50"></div>
                    <div class="w-2.5 h-2.5 rounded-full bg-yellow-500/50"></div>
                    <div class="w-2.5 h-2.5 rounded-full bg-green-500/50"></div>
                  </div>
                  <span class="text-[10px] font-medium text-gray-500 uppercase tracking-widest">{{ t('app.livePreview') }}</span>
                  <div class="w-10"></div>
                </div>
                <div class="max-h-[700px] overflow-y-auto custom-scrollbar">
                  <ProfileView :data="previewData" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
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
