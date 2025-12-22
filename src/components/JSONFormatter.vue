<script setup>
import { ref, computed, watch } from 'vue';
import { Code2, CheckCircle2, XCircle, Copy, Download, Trash2, Wand2, Minimize2, Info } from 'lucide-vue-next';

// Reactive state for JSON input and output
const inputJSON = ref('');
const outputJSON = ref('');
const isValid = ref(null);
const errorMessage = ref('');
const indentSize = ref(2);
const copySuccess = ref(false);
const activeTab = ref('format'); // format, minify, info

// Statistics about the JSON
const jsonStats = ref({
  size: 0,
  lines: 0,
  characters: 0,
  depth: 0,
  keys: 0,
  arrays: 0,
  objects: 0,
});

/**
 * Calculate the depth of a JSON object/array
 * @param {any} obj - Object to measure depth
 * @param {number} currentDepth - Current recursion depth
 * @returns {number} Maximum depth
 */
const calculateDepth = (obj, currentDepth = 1) => {
  if (typeof obj !== 'object' || obj === null) {
    return currentDepth;
  }
  
  let maxDepth = currentDepth;
  const values = Array.isArray(obj) ? obj : Object.values(obj);
  
  for (const value of values) {
    if (typeof value === 'object' && value !== null) {
      const depth = calculateDepth(value, currentDepth + 1);
      maxDepth = Math.max(maxDepth, depth);
    }
  }
  
  return maxDepth;
};

/**
 * Count keys, arrays, and objects in JSON structure
 * @param {any} obj - Object to analyze
 * @returns {Object} Counts of different JSON elements
 */
const analyzeStructure = (obj) => {
  const stats = { keys: 0, arrays: 0, objects: 0 };
  
  const traverse = (item) => {
    if (typeof item !== 'object' || item === null) {
      return;
    }
    
    if (Array.isArray(item)) {
      stats.arrays++;
      item.forEach(traverse);
    } else {
      stats.objects++;
      stats.keys += Object.keys(item).length;
      Object.values(item).forEach(traverse);
    }
  };
  
  traverse(obj);
  return stats;
};

/**
 * Format bytes to human-readable size
 * @param {number} bytes - Size in bytes
 * @returns {string} Formatted size string
 */
const formatBytes = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

/**
 * Validate and parse JSON input
 * Automatically converts JSON strings to JSON objects
 */
const validateJSON = () => {
  if (!inputJSON.value.trim()) {
    isValid.value = null;
    errorMessage.value = '';
    outputJSON.value = '';
    return;
  }
  
  try {
    let parsed = JSON.parse(inputJSON.value);
    
    // If the parsed result is a string, try to parse it again
    // This handles cases where JSON is double-encoded as a string
    if (typeof parsed === 'string') {
      try {
        const secondParse = JSON.parse(parsed);
        parsed = secondParse;
        // Update input to show the converted JSON
        inputJSON.value = JSON.stringify(parsed, null, indentSize.value);
      } catch (e) {
        // If second parse fails, the string is not JSON, keep original
      }
    }
    
    isValid.value = true;
    errorMessage.value = '';
    
    // Format based on active tab
    if (activeTab.value === 'format') {
      outputJSON.value = JSON.stringify(parsed, null, indentSize.value);
    } else if (activeTab.value === 'minify') {
      outputJSON.value = JSON.stringify(parsed);
    } else {
      outputJSON.value = JSON.stringify(parsed, null, indentSize.value);
    }
    
    // Calculate statistics
    const structure = analyzeStructure(parsed);
    jsonStats.value = {
      size: formatBytes(new Blob([outputJSON.value]).size),
      lines: outputJSON.value.split('\n').length,
      characters: outputJSON.value.length,
      depth: calculateDepth(parsed),
      keys: structure.keys,
      arrays: structure.arrays,
      objects: structure.objects,
    };
    
  } catch (error) {
    isValid.value = false;
    errorMessage.value = error.message;
    outputJSON.value = '';
    
    // Try to provide helpful error information
    const match = error.message.match(/position (\d+)/);
    if (match) {
      const position = parseInt(match[1]);
      const beforeError = inputJSON.value.substring(Math.max(0, position - 20), position);
      const afterError = inputJSON.value.substring(position, Math.min(inputJSON.value.length, position + 20));
      errorMessage.value += `\n\nNear: "${beforeError}⚠️${afterError}"`;
    }
  }
};

/**
 * Format JSON with specified indentation
 */
const formatJSON = () => {
  activeTab.value = 'format';
  validateJSON();
};

/**
 * Minify JSON (remove all whitespace)
 */
const minifyJSON = () => {
  activeTab.value = 'minify';
  validateJSON();
};

/**
 * Copy output to clipboard
 */
const copyToClipboard = async () => {
  if (!outputJSON.value) return;
  
  try {
    await navigator.clipboard.writeText(outputJSON.value);
    copySuccess.value = true;
    setTimeout(() => {
      copySuccess.value = false;
    }, 2000);
  } catch (error) {
    console.error('Failed to copy:', error);
  }
};

/**
 * Download JSON as file
 */
const downloadJSON = () => {
  if (!outputJSON.value) return;
  
  const blob = new Blob([outputJSON.value], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `formatted-${Date.now()}.json`;
  link.click();
  URL.revokeObjectURL(url);
};

/**
 * Clear all input and output
 */
const clearAll = () => {
  inputJSON.value = '';
  outputJSON.value = '';
  isValid.value = null;
  errorMessage.value = '';
  jsonStats.value = {
    size: 0,
    lines: 0,
    characters: 0,
    depth: 0,
    keys: 0,
    arrays: 0,
    objects: 0,
  };
};

/**
 * Load sample JSON for demonstration
 */
const loadSample = () => {
  inputJSON.value = JSON.stringify({
    "name": "Sample JSON",
    "version": "1.0.0",
    "description": "This is a sample JSON object for demonstration",
    "features": [
      "Validation",
      "Formatting",
      "Minification",
      "Statistics"
    ],
    "config": {
      "enabled": true,
      "timeout": 3000,
      "retries": 3,
      "endpoints": {
        "api": "https://api.example.com",
        "cdn": "https://cdn.example.com"
      }
    },
    "metadata": {
      "created": "2025-12-22",
      "author": "JSON Formatter Tool"
    }
  }, null, 2);
  validateJSON();
};

// Auto-validate on input change with debounce
let debounceTimer;
watch(inputJSON, () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    validateJSON();
  }, 500);
});

// Re-format when indent size changes
watch(indentSize, () => {
  if (isValid.value) {
    validateJSON();
  }
});

// Re-format when tab changes
watch(activeTab, () => {
  if (isValid.value) {
    validateJSON();
  }
});
</script>

<template>
  <div class="min-h-screen bg-[#0a0a0b] text-white py-12 px-4">
    <!-- Background Decoration -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] bg-indigo-500/10 blur-[120px] rounded-full"></div>
      <div class="absolute -bottom-[10%] -right-[10%] w-[40%] h-[40%] bg-purple-500/10 blur-[120px] rounded-full"></div>
    </div>

    <div class="relative z-10 max-w-7xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-12 animate-in fade-in duration-700">
        <div class="inline-flex items-center gap-3 mb-4">
          <div class="p-3 bg-indigo-600 rounded-2xl shadow-lg shadow-indigo-600/20">
            <Code2 class="w-8 h-8 text-white" />
          </div>
          <h1 class="text-5xl font-extrabold tracking-tight bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
            JSON Formatter
          </h1>
        </div>
        <p class="text-gray-400 text-lg max-w-2xl mx-auto">
          Validate, format, and minify JSON with automatic string-to-object conversion
        </p>
      </div>

      <!-- Main Content -->
      <div class="grid lg:grid-cols-2 gap-6 animate-in slide-in-from-bottom duration-700">
        <!-- Input Panel -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-2xl font-bold flex items-center gap-2">
              Input
              <span v-if="isValid !== null" class="text-sm">
                <CheckCircle2 v-if="isValid" class="w-5 h-5 text-green-500 inline" />
                <XCircle v-else class="w-5 h-5 text-red-500 inline" />
              </span>
            </h2>
            <div class="flex gap-2">
              <button
                @click="loadSample"
                class="px-3 py-1.5 text-sm bg-white/5 hover:bg-white/10 rounded-lg transition-colors flex items-center gap-1"
                title="Load Sample"
              >
                <Wand2 class="w-4 h-4" />
                Sample
              </button>
              <button
                @click="clearAll"
                class="px-3 py-1.5 text-sm bg-white/5 hover:bg-white/10 rounded-lg transition-colors flex items-center gap-1"
                title="Clear All"
              >
                <Trash2 class="w-4 h-4" />
                Clear
              </button>
            </div>
          </div>

          <div class="bg-[#111113] border border-white/5 rounded-2xl overflow-hidden">
            <textarea
              v-model="inputJSON"
              placeholder='Paste your JSON here... e.g., {"name": "value"}'
              class="w-full h-[500px] bg-transparent p-6 font-mono text-sm resize-none focus:outline-none custom-scrollbar"
              spellcheck="false"
            ></textarea>
          </div>

          <!-- Error Message -->
          <div v-if="!isValid && errorMessage" class="bg-red-500/10 border border-red-500/20 rounded-xl p-4">
            <div class="flex items-start gap-2">
              <XCircle class="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
              <div>
                <div class="font-semibold text-red-400 mb-1">Invalid JSON</div>
                <div class="text-sm text-red-300 whitespace-pre-wrap font-mono">{{ errorMessage }}</div>
              </div>
            </div>
          </div>

          <!-- Controls -->
          <div class="bg-[#111113] border border-white/5 rounded-2xl p-4">
            <div class="flex flex-wrap items-center gap-4">
              <div class="flex items-center gap-2">
                <label class="text-sm text-gray-400">Indent Size:</label>
                <select
                  v-model.number="indentSize"
                  class="bg-white/5 border border-white/10 rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:border-indigo-500 transition-colors"
                >
                  <option :value="2">2 spaces</option>
                  <option :value="4">4 spaces</option>
                  <option :value="8">8 spaces</option>
                </select>
              </div>

              <div class="flex gap-2 flex-1 justify-end">
                <button
                  @click="formatJSON"
                  :disabled="!isValid"
                  :class="[
                    'px-4 py-2 rounded-lg font-semibold transition-all duration-300 flex items-center gap-2',
                    isValid
                      ? 'bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 shadow-lg shadow-indigo-600/20'
                      : 'bg-white/5 text-gray-600 cursor-not-allowed'
                  ]"
                >
                  <Code2 class="w-4 h-4" />
                  Format
                </button>
                <button
                  @click="minifyJSON"
                  :disabled="!isValid"
                  :class="[
                    'px-4 py-2 rounded-lg font-semibold transition-all duration-300 flex items-center gap-2',
                    isValid
                      ? 'bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 shadow-lg shadow-purple-600/20'
                      : 'bg-white/5 text-gray-600 cursor-not-allowed'
                  ]"
                >
                  <Minimize2 class="w-4 h-4" />
                  Minify
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Output Panel -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-2xl font-bold">Output</h2>
            <div class="flex gap-2">
              <button
                @click="copyToClipboard"
                :disabled="!outputJSON"
                :class="[
                  'px-3 py-1.5 text-sm rounded-lg transition-all duration-300 flex items-center gap-1',
                  copySuccess
                    ? 'bg-green-500 text-white'
                    : outputJSON
                      ? 'bg-white/5 hover:bg-white/10'
                      : 'bg-white/5 text-gray-600 cursor-not-allowed'
                ]"
                title="Copy to Clipboard"
              >
                <Copy class="w-4 h-4" />
                {{ copySuccess ? 'Copied!' : 'Copy' }}
              </button>
              <button
                @click="downloadJSON"
                :disabled="!outputJSON"
                :class="[
                  'px-3 py-1.5 text-sm rounded-lg transition-colors flex items-center gap-1',
                  outputJSON
                    ? 'bg-white/5 hover:bg-white/10'
                    : 'bg-white/5 text-gray-600 cursor-not-allowed'
                ]"
                title="Download JSON"
              >
                <Download class="w-4 h-4" />
                Download
              </button>
            </div>
          </div>

          <div class="bg-[#111113] border border-white/5 rounded-2xl overflow-hidden">
            <textarea
              v-model="outputJSON"
              readonly
              placeholder="Formatted JSON will appear here..."
              class="w-full h-[500px] bg-transparent p-6 font-mono text-sm resize-none focus:outline-none custom-scrollbar"
              spellcheck="false"
            ></textarea>
          </div>

          <!-- Statistics -->
          <div v-if="isValid && outputJSON" class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-4">
            <h3 class="text-lg font-semibold text-indigo-400 flex items-center gap-2 border-b border-white/5 pb-2">
              <Info class="w-5 h-5" />
              Statistics
            </h3>
            
            <div class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <div class="text-gray-400 mb-1">File Size</div>
                <div class="font-medium">{{ jsonStats.size }}</div>
              </div>
              
              <div>
                <div class="text-gray-400 mb-1">Lines</div>
                <div class="font-medium">{{ jsonStats.lines.toLocaleString() }}</div>
              </div>
              
              <div>
                <div class="text-gray-400 mb-1">Characters</div>
                <div class="font-medium">{{ jsonStats.characters.toLocaleString() }}</div>
              </div>
              
              <div>
                <div class="text-gray-400 mb-1">Max Depth</div>
                <div class="font-medium">{{ jsonStats.depth }}</div>
              </div>
              
              <div>
                <div class="text-gray-400 mb-1">Total Keys</div>
                <div class="font-medium">{{ jsonStats.keys.toLocaleString() }}</div>
              </div>
              
              <div>
                <div class="text-gray-400 mb-1">Objects</div>
                <div class="font-medium">{{ jsonStats.objects }}</div>
              </div>
              
              <div>
                <div class="text-gray-400 mb-1">Arrays</div>
                <div class="font-medium">{{ jsonStats.arrays }}</div>
              </div>
              
              <div>
                <div class="text-gray-400 mb-1">Mode</div>
                <div class="font-medium">{{ activeTab === 'format' ? 'Formatted' : 'Minified' }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Features Info -->
      <div class="mt-12 grid md:grid-cols-3 gap-6 animate-in fade-in duration-1000">
        <div class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-3">
          <div class="w-12 h-12 bg-indigo-600/20 rounded-xl flex items-center justify-center">
            <CheckCircle2 class="w-6 h-6 text-indigo-400" />
          </div>
          <h3 class="text-lg font-semibold">Auto Validation</h3>
          <p class="text-gray-400 text-sm">
            Automatically validates JSON as you type with detailed error messages
          </p>
        </div>

        <div class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-3">
          <div class="w-12 h-12 bg-purple-600/20 rounded-xl flex items-center justify-center">
            <Wand2 class="w-6 h-6 text-purple-400" />
          </div>
          <h3 class="text-lg font-semibold">Smart Conversion</h3>
          <p class="text-gray-400 text-sm">
            Automatically converts JSON strings to objects when detected
          </p>
        </div>

        <div class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-3">
          <div class="w-12 h-12 bg-pink-600/20 rounded-xl flex items-center justify-center">
            <Info class="w-6 h-6 text-pink-400" />
          </div>
          <h3 class="text-lg font-semibold">Detailed Stats</h3>
          <p class="text-gray-400 text-sm">
            View size, depth, key count, and structure analysis
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar styling */
.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
  height: 8px;
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

/* Custom animations */
@keyframes in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-in {
  animation: in 0.5s ease-out;
}

.slide-in-from-bottom {
  animation: in 0.7s ease-out;
}

/* Textarea focus styling */
textarea:focus {
  outline: none;
}

/* Smooth transitions */
button {
  transition: all 0.3s ease;
}

button:disabled {
  opacity: 0.5;
}
</style>
