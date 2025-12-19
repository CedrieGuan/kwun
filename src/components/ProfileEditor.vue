<script setup>
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { 
  Plus, 
  Trash2, 
  Copy, 
  Check, 
  Github, 
  Twitter, 
  Instagram, 
  Linkedin, 
  Globe,
  Share2,
  Languages
} from 'lucide-vue-next';
import { encodeProfile } from '../utils/onelink';

const { t, locale } = useI18n();
const emit = defineEmits(['update']);

const profile = ref({
  name: '',
  bio: '',
  avatar: '',
  links: [
    { title: '', url: '' }
  ],
  socials: {
    github: '',
    twitter: '',
    instagram: '',
    linkedin: '',
    website: ''
  }
});

const copied = ref(false);

const generatedUrl = computed(() => {
  if (!profile.value.name) return '';
  const encoded = encodeProfile(profile.value);
  const baseUrl = window.location.origin;
  return `${baseUrl}?data=${encoded}`;
});

const addLink = () => {
  profile.value.links.push({ title: '', url: '' });
};

const removeLink = (index) => {
  profile.value.links.splice(index, 1);
};

const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(generatedUrl.value);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch (err) {
    console.error('Failed to copy:', err);
  }
};

watch(profile, (newVal) => {
  emit('update', newVal);
}, { deep: true });
</script>

<template>
  <div class="w-full max-w-2xl bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-2xl space-y-8">
    <div class="flex flex-col md:flex-row md:items-start justify-between gap-4">
      <div class="space-y-2">
        <h2 class="text-2xl font-bold text-white flex items-center gap-2">
          <Share2 class="w-6 h-6 text-indigo-400" />
          {{ t('app.subtitle') }}
        </h2>
        <p class="text-gray-400">{{ t('app.description') }}</p>
      </div>
      
      <!-- Language Switcher -->
      <div class="flex items-center gap-2 bg-white/5 border border-white/10 rounded-xl p-1">
        <button v-for="lang in ['en', 'zh-CN', 'zh-TW', 'ja']" :key="lang"
                @click="locale = lang"
                :class="[
                  'px-2 py-1 text-xs font-bold rounded-lg transition-all',
                  locale === lang ? 'bg-indigo-600 text-white shadow-lg' : 'text-gray-500 hover:text-white'
                ]"
        >
          {{ lang.toUpperCase() }}
        </button>
      </div>
    </div>

    <!-- Basic Info -->
    <div class="space-y-4">
      <h3 class="text-sm font-semibold uppercase tracking-wider text-gray-500">{{ t('editor.basicInfo') }}</h3>
      <div class="grid grid-cols-1 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-400 mb-1">{{ t('editor.nameLabel') }}</label>
          <input v-model="profile.name" type="text" :placeholder="t('editor.namePlaceholder')" 
                 class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-400 mb-1">{{ t('editor.bioLabel') }}</label>
          <textarea v-model="profile.bio" :placeholder="t('editor.bioPlaceholder')" rows="2"
                    class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all resize-none"></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-400 mb-1">{{ t('editor.avatarLabel') }}</label>
          <input v-model="profile.avatar" type="text" :placeholder="t('editor.avatarPlaceholder')" 
                 class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all" />
        </div>
      </div>
    </div>

    <!-- Socials -->
    <div class="space-y-4">
      <h3 class="text-sm font-semibold uppercase tracking-wider text-gray-500">{{ t('editor.socialPresence') }}</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex items-center gap-2 bg-white/5 border border-white/10 rounded-xl px-3 py-2">
          <Github class="w-5 h-5 text-gray-400" />
          <input v-model="profile.socials.github" type="text" placeholder="username" class="bg-transparent border-none focus:ring-0 text-white w-full placeholder-gray-600" />
        </div>
        <div class="flex items-center gap-2 bg-white/5 border border-white/10 rounded-xl px-3 py-2">
          <Twitter class="w-5 h-5 text-gray-400" />
          <input v-model="profile.socials.twitter" type="text" placeholder="username" class="bg-transparent border-none focus:ring-0 text-white w-full placeholder-gray-600" />
        </div>
        <div class="flex items-center gap-2 bg-white/5 border border-white/10 rounded-xl px-3 py-2">
          <Instagram class="w-5 h-5 text-gray-400" />
          <input v-model="profile.socials.instagram" type="text" placeholder="username" class="bg-transparent border-none focus:ring-0 text-white w-full placeholder-gray-600" />
        </div>
        <div class="flex items-center gap-2 bg-white/5 border border-white/10 rounded-xl px-3 py-2">
          <Linkedin class="w-5 h-5 text-gray-400" />
          <input v-model="profile.socials.linkedin" type="text" placeholder="username" class="bg-transparent border-none focus:ring-0 text-white w-full placeholder-gray-600" />
        </div>
        <div class="flex items-center gap-2 bg-white/5 border border-white/10 rounded-xl px-3 py-2 md:col-span-2">
          <Globe class="w-5 h-5 text-gray-400" />
          <input v-model="profile.socials.website" type="text" placeholder="https://yourwebsite.com" class="bg-transparent border-none focus:ring-0 text-white w-full placeholder-gray-600" />
        </div>
      </div>
    </div>

    <!-- Custom Links -->
    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-semibold uppercase tracking-wider text-gray-500">{{ t('editor.links') }}</h3>
        <button @click="addLink" class="text-xs flex items-center gap-1 text-indigo-400 hover:text-indigo-300 transition-colors uppercase font-bold tracking-widest">
          <Plus class="w-3 h-3" /> {{ t('editor.addLink') }}
        </button>
      </div>
      <div class="space-y-3">
        <div v-for="(link, index) in profile.links" :key="index" class="flex items-start gap-4 animate-in fade-in slide-in-from-top-4 duration-300">
          <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-3">
             <input v-model="link.title" type="text" :placeholder="t('editor.linkTitle')" class="bg-white/5 border border-white/10 rounded-xl px-4 py-2 text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50" />
             <input v-model="link.url" type="text" placeholder="https://..." class="bg-white/5 border border-white/10 rounded-xl px-4 py-2 text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50" />
          </div>
          <button @click="removeLink(index)" class="mt-2 text-gray-600 hover:text-red-400 transition-colors">
            <Trash2 class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- URL Generator -->
    <div v-if="generatedUrl" class="pt-6 border-t border-white/10 animate-in fade-in zoom-in duration-500">
      <label class="block text-sm font-medium text-indigo-300 mb-2">{{ t('editor.generatedLink') }}</label>
      <div class="flex items-center gap-2 bg-indigo-500/10 border border-indigo-500/20 rounded-2xl p-4">
        <input readonly :value="generatedUrl" class="bg-transparent border-none focus:ring-0 text-indigo-100 w-full text-sm truncate" />
        <button @click="copyUrl" class="flex-shrink-0 bg-indigo-600 hover:bg-indigo-500 text-white p-2 rounded-xl transition-all shadow-lg active:scale-95">
          <Check v-if="copied" class="w-5 h-5" />
          <Copy v-else class="w-5 h-5" />
        </button>
      </div>
      <p class="mt-2 text-xs text-center text-gray-500">{{ t('app.note') }}</p>
    </div>
  </div>
</template>
