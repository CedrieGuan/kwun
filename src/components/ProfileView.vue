<script setup>
import { computed } from 'vue';
import { 
  Github, 
  Twitter, 
  Instagram, 
  ExternalLink,
  Linkedin,
  Globe
} from 'lucide-vue-next';

const props = defineProps({
  data: {
    type: Object,
    required: true
  }
});

const socialIcons = {
  github: Github,
  twitter: Twitter,
  instagram: Instagram,
  linkedin: Linkedin,
  website: Globe
};

const hasSocials = computed(() => {
  return props.data.socials && Object.values(props.data.socials).some(v => v);
});
</script>

<template>
  <div class="min-h-screen py-12 px-4 flex flex-col items-center">
    <div class="max-w-md w-full flex flex-col items-center text-center space-y-6">
      <!-- Avatar -->
      <div v-if="data.avatar" class="w-24 h-24 rounded-full overflow-hidden border-2 border-primary/20 shadow-xl">
        <img :src="data.avatar" :alt="data.name" class="w-full h-full object-cover" />
      </div>
      <div v-else class="w-24 h-24 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-3xl font-bold shadow-lg">
        {{ data.name ? data.name[0].toUpperCase() : '?' }}
      </div>

      <!-- Header -->
      <div class="space-y-2">
        <h1 class="text-3xl font-bold tracking-tight text-white">{{ data.name || $t('common.anonymous') }}</h1>
        <p class="text-gray-400 text-lg">{{ data.bio }}</p>
      </div>

      <!-- Socials -->
      <div v-if="hasSocials" class="flex flex-wrap justify-center gap-4">
        <template v-for="(handle, platform) in data.socials" :key="platform">
          <a v-if="handle" 
             :href="platform === 'website' ? handle : `https://${platform}.com/${handle}`"
             target="_blank"
             class="p-2 rounded-full bg-white/5 hover:bg-white/10 text-gray-300 hover:text-white transition-all transform hover:scale-110"
          >
            <component :is="socialIcons[platform] || ExternalLink" class="w-6 h-6" />
          </a>
        </template>
      </div>

      <!-- Links -->
      <div class="w-full space-y-4 pt-4">
        <a v-for="(link, index) in data.links" :key="index"
           :href="link.url"
           target="_blank"
           class="group block w-full p-4 rounded-xl bg-white/5 border border-white/10 hover:border-white/20 hover:bg-white/10 transition-all transform hover:-translate-y-1 active:scale-[0.98] shadow-sm"
        >
          <div class="flex items-center justify-between">
            <span class="text-lg font-medium text-white">{{ link.title }}</span>
            <ExternalLink class="w-5 h-5 text-gray-500 group-hover:text-white transition-colors" />
          </div>
        </a>
      </div>

      <!-- Footer -->
      <div class="pt-8">
        <a href="/" class="text-sm text-gray-500 hover:text-gray-300 transition-colors flex items-center gap-1">
          {{ $t('app.createOwn') }}
          <ExternalLink class="w-3 h-3" />
        </a>
      </div>
    </div>
  </div>
</template>
