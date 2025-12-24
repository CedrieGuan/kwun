<template>
  <div class="min-h-screen bg-[#0a0a0b] text-white">
    <div class="container mx-auto px-4 py-12">
      <!-- Header -->
      <div class="max-w-6xl mx-auto space-y-8">
        <div class="text-center space-y-4">
          <div class="inline-flex items-center justify-center p-4 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-2xl shadow-lg shadow-indigo-600/20">
            <Languages class="w-12 h-12 text-white" />
          </div>
          <h1 class="text-5xl font-extrabold bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
            Language Translator
          </h1>
          <p class="text-gray-400 text-lg">
            Translate text between multiple languages powered by DeepL
          </p>
        </div>

        <!-- Usage Quota Display -->
        <div v-if="usageData && usageData.character_limit > 0" class="bg-white/5 backdrop-blur-xl rounded-xl border border-white/10 p-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <Activity class="w-5 h-5 text-indigo-400" />
              <span class="text-sm font-medium text-gray-300">Monthly Usage</span>
            </div>
            <div class="text-sm">
              <span class="font-semibold text-white">{{ formatNumber(usageData.character_count) }}</span>
              <span class="text-gray-400"> / </span>
              <span class="text-gray-300">{{ formatNumber(usageData.character_limit) }}</span>
              <span class="text-gray-400 ml-1">characters</span>
            </div>
          </div>
          <div class="mt-2">
            <div class="w-full bg-black/30 rounded-full h-2">
              <div 
                class="bg-gradient-to-r from-indigo-600 to-purple-600 h-2 rounded-full transition-all duration-500"
                :style="{ width: usagePercentage + '%' }"
              ></div>
            </div>
          </div>
        </div>

        <!-- Main Translation Interface -->
        <div class="grid md:grid-cols-2 gap-6">
          <!-- Source Language Section -->
          <div class="space-y-4">
            <div class="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6 space-y-4">
              <!-- Source Language Selector -->
              <div class="flex items-center justify-between">
                <select
                  v-model="sourceLang"
                  class="px-4 py-2 bg-black/30 border border-white/10 rounded-xl text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                >
                  <option value="auto">Auto Detect</option>
                  <option v-for="lang in languages" :key="lang.code" :value="lang.code">
                    {{ lang.name }}
                  </option>
                </select>
                
                <button
                  @click="swapLanguages"
                  :disabled="sourceLang === 'auto'"
                  :class="[
                    'p-2 rounded-xl transition-all duration-300',
                    sourceLang === 'auto'
                      ? 'bg-white/5 text-gray-500 cursor-not-allowed'
                      : 'bg-white/10 hover:bg-white/20 text-indigo-400 hover:text-indigo-300'
                  ]"
                >
                  <ArrowLeftRight class="w-5 h-5" />
                </button>
              </div>

              <!-- Source Text Input -->
              <div class="space-y-2">
                <label class="text-sm font-medium text-gray-300">Source Text</label>
                <textarea
                  v-model="sourceText"
                  @input="handleSourceTextChange"
                  placeholder="Enter text to translate..."
                  class="w-full h-64 px-4 py-3 bg-black/30 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
                ></textarea>
                <div class="flex items-center justify-between text-xs text-gray-400">
                  <span>{{ sourceText.length }} characters</span>
                  <button
                    v-if="sourceText"
                    @click="clearAll"
                    class="flex items-center gap-1 px-2 py-1 hover:text-red-400 transition-colors"
                  >
                    <X class="w-3 h-3" />
                    Clear
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Target Language Section -->
          <div class="space-y-4">
            <div class="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6 space-y-4">
              <!-- Target Language Selector -->
              <div class="flex items-center justify-between">
                <select
                  v-model="targetLang"
                  class="px-4 py-2 bg-black/30 border border-white/10 rounded-xl text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                >
                  <option v-for="lang in languages" :key="lang.code" :value="lang.code">
                    {{ lang.name }}
                  </option>
                </select>

                <button
                  v-if="translatedText && !isTranslating"
                  @click="copyTranslation"
                  class="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-xl text-sm font-medium transition-colors"
                >
                  <Copy class="w-4 h-4" />
                  {{ copied ? 'Copied!' : 'Copy' }}
                </button>
              </div>

              <!-- Translation Output -->
              <div class="space-y-2">
                <div class="flex items-center justify-between">
                  <label class="text-sm font-medium text-gray-300">Translation</label>
                  <span v-if="detectedLanguage" class="text-xs text-indigo-400">
                    Detected: {{ getLanguageName(detectedLanguage) }}
                  </span>
                </div>
                <div class="relative">
                  <textarea
                    v-model="translatedText"
                    readonly
                    :placeholder="isTranslating ? 'Translating...' : errorMessage || 'Translation will appear here...'"
                    :class="[
                      'w-full h-64 px-4 py-3 border rounded-xl focus:outline-none resize-none',
                      errorMessage
                        ? 'bg-red-500/10 border-red-500/30 text-red-400'
                        : 'bg-black/30 border-white/10 text-purple-300 placeholder-gray-500'
                    ]"
                  ></textarea>
                  
                  <!-- Loading Spinner -->
                  <div v-if="isTranslating" class="absolute inset-0 flex items-center justify-center bg-black/20 backdrop-blur-sm rounded-xl">
                    <div class="flex flex-col items-center gap-3">
                      <div class="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
                      <span class="text-sm text-gray-300">Translating...</span>
                    </div>
                  </div>
                </div>
                <div class="text-xs text-gray-400">
                  {{ translatedText.length }} characters
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-center gap-4">
          <button
            @click="translate"
            :disabled="!sourceText || !targetLang || isTranslating"
            :class="[
              'flex items-center gap-2 px-8 py-3 rounded-xl font-semibold transition-all duration-300',
              sourceText && targetLang && !isTranslating
                ? 'bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-white shadow-lg shadow-indigo-600/20'
                : 'bg-white/5 text-gray-500 cursor-not-allowed'
            ]"
          >
            <Languages class="w-5 h-5" />
            Translate
          </button>
        </div>

        <!-- Supported Languages Info -->
        <div class="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6">
          <h3 class="text-lg font-semibold text-white mb-3">About DeepL Translation</h3>
          <p class="text-gray-400 text-sm leading-relaxed mb-4">
            DeepL uses advanced neural networks to provide highly accurate translations. 
            The free API allows up to 500,000 characters per month. 
            Supported languages include English, Chinese, Japanese, Korean, Spanish, French, German, and many more.
          </p>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="lang in languages.slice(0, 10)"
              :key="lang.code"
              class="px-3 py-1 bg-white/5 rounded-lg text-xs text-gray-300"
            >
              {{ lang.name }}
            </span>
            <span class="px-3 py-1 bg-white/5 rounded-lg text-xs text-gray-400">
              +{{ languages.length - 10 }} more
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { Languages, ArrowLeftRight, Copy, X, Activity } from 'lucide-vue-next';

// Supported languages based on DeepL API
// 基于DeepL API的支持语言
const languages = [
  { code: 'EN', name: 'English' },
  { code: 'ZH', name: 'Chinese (Simplified)' },
  { code: 'JA', name: 'Japanese' },
  { code: 'KO', name: 'Korean' },
  { code: 'ES', name: 'Spanish' },
  { code: 'FR', name: 'French' },
  { code: 'DE', name: 'German' },
  { code: 'IT', name: 'Italian' },
  { code: 'PT', name: 'Portuguese' },
  { code: 'RU', name: 'Russian' },
  { code: 'AR', name: 'Arabic' },
  { code: 'NL', name: 'Dutch' },
  { code: 'PL', name: 'Polish' },
  { code: 'TR', name: 'Turkish' },
  { code: 'SV', name: 'Swedish' },
  { code: 'DA', name: 'Danish' },
  { code: 'FI', name: 'Finnish' },
  { code: 'NO', name: 'Norwegian' },
  { code: 'CS', name: 'Czech' },
  { code: 'RO', name: 'Romanian' },
];

// State management
// 状态管理
const sourceLang = ref('auto');
const targetLang = ref('EN');
const sourceText = ref('');
const translatedText = ref('');
const detectedLanguage = ref('');
const isTranslating = ref(false);
const errorMessage = ref('');
const copied = ref(false);

// Usage data
// 使用量数据
const usageData = ref(null);

// Debounce timer for auto-translate
// 自动翻译的防抖定时器
let translateDebounce = null;

/**
 * Fetch DeepL usage statistics
 * 获取DeepL使用统计信息
 */
const fetchUsage = async () => {
  try {
    const apiUrl = import.meta.env.DEV 
      ? 'http://localhost:8080/api/usage'
      : '/api/usage';
    
    const response = await fetch(apiUrl, {
      method: 'GET',
    });

    if (response.ok) {
      const data = await response.json();
      usageData.value = data;
    }
  } catch (error) {
    console.error('Failed to fetch usage data:', error);
    // Graceful degradation - just don't show usage
    // 优雅降级 - 不显示使用量
  }
};

/**
 * Calculate usage percentage
 * 计算使用百分比
 */
const usagePercentage = computed(() => {
  if (!usageData.value || usageData.value.character_limit === 0) return 0;
  return Math.min(100, (usageData.value.character_count / usageData.value.character_limit) * 100);
});

/**
 * Format number with thousand separators
 * 格式化数字，添加千位分隔符
 */
const formatNumber = (num) => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

/**
 * Get language name from code
 * 从代码获取语言名称
 */
const getLanguageName = (code) => {
  const lang = languages.find(l => l.code.toUpperCase() === code.toUpperCase());
  return lang ? lang.name : code;
};

/**
 * Handle source text changes with debounce
 * 处理源文本变化（带防抖）
 */
const handleSourceTextChange = () => {
  // Clear previous translation if text changes
  // 如果文本改变，清除之前的翻译
  translatedText.value = '';
  errorMessage.value = '';
  detectedLanguage.value = '';
};

/**
 * Perform translation using DeepL API
 * 使用DeepL API执行翻译
 */
const translate = async () => {
  if (!sourceText.value || !targetLang.value) return;

  isTranslating.value = true;
  errorMessage.value = '';
  translatedText.value = '';
  detectedLanguage.value = '';

  try {
    const apiUrl = import.meta.env.DEV 
      ? 'http://localhost:8080/api/translate'
      : '/api/translate';

    const response = await fetch(apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        text: sourceText.value,
        source_lang: sourceLang.value === 'auto' ? '' : sourceLang.value,
        target_lang: targetLang.value,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Translation failed');
    }

    const data = await response.json();
    translatedText.value = data.translated_text;
    if (data.detected_language) {
      detectedLanguage.value = data.detected_language;
    }

    // Refresh usage data after translation
    // 翻译后刷新使用量数据
    fetchUsage();
  } catch (error) {
    console.error('Translation error:', error);
    errorMessage.value = error.message || 'Failed to translate. Please try again.';
  } finally {
    isTranslating.value = false;
  }
};

/**
 * Swap source and target languages
 * 交换源语言和目标语言
 */
const swapLanguages = () => {
  if (sourceLang.value === 'auto') return;
  
  const temp = sourceLang.value;
  sourceLang.value = targetLang.value;
  targetLang.value = temp;

  // Also swap text if translation exists
  // 如果翻译存在，也交换文本
  if (translatedText.value) {
    const tempText = sourceText.value;
    sourceText.value = translatedText.value;
    translatedText.value = tempText;
  }
};

/**
 * Copy translation to clipboard
 * 复制翻译到剪贴板
 */
const copyTranslation = async () => {
  try {
    await navigator.clipboard.writeText(translatedText.value);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch (error) {
    console.error('Copy failed:', error);
  }
};

/**
 * Clear all inputs and outputs
 * 清除所有输入和输出
 */
const clearAll = () => {
  sourceText.value = '';
  translatedText.value = '';
  detectedLanguage.value = '';
  errorMessage.value = '';
  copied.value = false;
};

/**
 * Initialize component
 * 初始化组件
 */
onMounted(() => {
  fetchUsage();
});
</script>
