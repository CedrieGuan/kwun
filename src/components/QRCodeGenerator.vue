<script setup>
import { ref, watch, nextTick } from 'vue';
import { QrCode, Download, Copy, Check, Settings, RefreshCw, Palette } from 'lucide-vue-next';
import QRCode from 'qrcode';

// State management for input and generated QR code
const inputText = ref('');
const qrCanvas = ref(null);
const qrDataURL = ref('');
const copySuccess = ref(false);
const generatedAt = ref(null);

// Advanced QR code options
const qrOptions = ref({
  errorCorrectionLevel: 'M', // L, M, Q, H
  type: 'image/png',
  quality: 0.95,
  margin: 4,
  width: 400,
  color: {
    dark: '#000000',
    light: '#FFFFFF'
  }
});

// Statistics about the generated QR code
const statistics = ref({
  textLength: 0,
  qrSize: 0,
  errorCorrection: 'Medium',
  estimatedCapacity: 0
});

// Preset templates for quick access
const presets = ref([
  { name: 'URL', value: 'https://example.com', icon: '🌐' },
  { name: 'Email', value: 'mailto:contact@example.com', icon: '📧' },
  { name: 'Phone', value: 'tel:+1234567890', icon: '📱' },
  { name: 'WiFi', value: 'WIFI:T:WPA;S:NetworkName;P:Password;;', icon: '📶' },
  { name: 'SMS', value: 'sms:+1234567890?body=Hello', icon: '💬' },
  { name: 'VCard', value: 'BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+1234567890\nEMAIL:john@example.com\nEND:VCARD', icon: '👤' }
]);

// Color presets for customization
const colorPresets = ref([
  { name: 'Classic', dark: '#000000', light: '#FFFFFF' },
  { name: 'Indigo', dark: '#4F46E5', light: '#EEF2FF' },
  { name: 'Purple', dark: '#7C3AED', light: '#F5F3FF' },
  { name: 'Blue', dark: '#2563EB', light: '#EFF6FF' },
  { name: 'Green', dark: '#059669', light: '#ECFDF5' },
  { name: 'Red', dark: '#DC2626', light: '#FEF2F2' },
  { name: 'Orange', dark: '#EA580C', light: '#FFF7ED' },
  { name: 'Dark Mode', dark: '#F9FAFB', light: '#111827' }
]);

// Show/hide advanced options panel
const showAdvancedOptions = ref(false);

/**
 * Generate QR code from input text
 * This function validates input, generates QR code using the qrcode library,
 * and updates statistics
 */
const generateQRCode = async () => {
  if (!inputText.value.trim()) {
    qrDataURL.value = '';
    statistics.value.textLength = 0;
    return;
  }

  try {
    // Generate QR code as data URL
    const dataURL = await QRCode.toDataURL(inputText.value, {
      errorCorrectionLevel: qrOptions.value.errorCorrectionLevel,
      type: qrOptions.value.type,
      quality: qrOptions.value.quality,
      margin: qrOptions.value.margin,
      width: qrOptions.value.width,
      color: {
        dark: qrOptions.value.color.dark,
        light: qrOptions.value.color.light
      }
    });

    qrDataURL.value = dataURL;
    generatedAt.value = new Date();

    // Update statistics
    updateStatistics();
  } catch (error) {
    console.error('Error generating QR code:', error);
  }
};

/**
 * Update statistics about the generated QR code
 * Includes text length, error correction level, and estimated capacity
 */
const updateStatistics = () => {
  statistics.value.textLength = inputText.value.length;
  
  // Map error correction level to human-readable format
  const ecLevelMap = {
    'L': 'Low (~7%)',
    'M': 'Medium (~15%)',
    'Q': 'Quartile (~25%)',
    'H': 'High (~30%)'
  };
  statistics.value.errorCorrection = ecLevelMap[qrOptions.value.errorCorrectionLevel] || 'Medium';

  // Estimate capacity based on error correction level
  // These are approximate values for numeric data
  const capacityMap = {
    'L': 7089,
    'M': 5596,
    'Q': 3993,
    'H': 3057
  };
  statistics.value.estimatedCapacity = capacityMap[qrOptions.value.errorCorrectionLevel] || 5596;

  // Calculate approximate QR code size in KB
  if (qrDataURL.value) {
    const base64Length = qrDataURL.value.length - 'data:image/png;base64,'.length;
    statistics.value.qrSize = Math.round((base64Length * 3 / 4) / 1024 * 100) / 100;
  }
};

/**
 * Download the generated QR code as PNG image
 */
const downloadQRCode = () => {
  if (!qrDataURL.value) return;

  const link = document.createElement('a');
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5);
  link.download = `qrcode-${timestamp}.png`;
  link.href = qrDataURL.value;
  link.click();
};

/**
 * Copy QR code image to clipboard
 */
const copyQRCode = async () => {
  if (!qrDataURL.value) return;

  try {
    // Convert data URL to blob
    const response = await fetch(qrDataURL.value);
    const blob = await response.blob();

    // Copy to clipboard
    await navigator.clipboard.write([
      new ClipboardItem({ 'image/png': blob })
    ]);

    copySuccess.value = true;
    setTimeout(() => {
      copySuccess.value = false;
    }, 2000);
  } catch (error) {
    console.error('Error copying QR code:', error);
    // Fallback: try to copy the data URL as text
    try {
      await navigator.clipboard.writeText(qrDataURL.value);
      copySuccess.value = true;
      setTimeout(() => {
        copySuccess.value = false;
      }, 2000);
    } catch (fallbackError) {
      console.error('Fallback copy failed:', fallbackError);
    }
  }
};

/**
 * Apply a preset template to the input
 */
const applyPreset = (preset) => {
  inputText.value = preset.value;
  generateQRCode();
};

/**
 * Apply a color preset to the QR code
 */
const applyColorPreset = (preset) => {
  qrOptions.value.color.dark = preset.dark;
  qrOptions.value.color.light = preset.light;
  if (inputText.value) {
    generateQRCode();
  }
};

/**
 * Clear all input and reset the form
 */
const clearAll = () => {
  inputText.value = '';
  qrDataURL.value = '';
  generatedAt.value = null;
  statistics.value.textLength = 0;
  statistics.value.qrSize = 0;
};

/**
 * Format file size in human-readable format
 */
const formatSize = (sizeInKB) => {
  if (sizeInKB < 1) {
    return `${Math.round(sizeInKB * 1024)} B`;
  }
  return `${sizeInKB.toFixed(2)} KB`;
};

/**
 * Format timestamp for display
 */
const formatTime = (date) => {
  if (!date) return '';
  return date.toLocaleTimeString();
};

// Watch for input changes and regenerate QR code
watch(inputText, () => {
  generateQRCode();
});

// Watch for option changes and regenerate QR code
watch(() => qrOptions.value, () => {
  if (inputText.value) {
    generateQRCode();
  }
}, { deep: true });
</script>

<template>
  <div class="min-h-screen bg-[#0a0a0b] py-8 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- Header Section -->
      <div class="mb-8 text-center animate-in fade-in slide-in-from-top duration-700">
        <div class="inline-flex items-center justify-center gap-3 mb-4">
          <div class="p-3 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-2xl shadow-lg shadow-indigo-600/30">
            <QrCode class="w-8 h-8 text-white" />
          </div>
          <h1 class="text-4xl md:text-5xl font-extrabold tracking-tight bg-gradient-to-r from-white via-indigo-100 to-purple-200 bg-clip-text text-transparent">
            QR Code Generator
          </h1>
        </div>
        <p class="text-gray-400 text-lg max-w-2xl mx-auto">
          Create beautiful QR codes for URLs, text, WiFi credentials, and more with advanced customization options
        </p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Left Column: Input and Options -->
        <div class="space-y-6 animate-in fade-in slide-in-from-left duration-700">
          <!-- Input Section -->
          <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6 shadow-xl transition-all duration-300 hover:border-indigo-500/30">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-xl font-bold text-white">Input Content</h2>
              <button
                @click="clearAll"
                class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-gray-400 hover:text-white transition-all duration-200 text-sm"
              >
                <RefreshCw class="w-4 h-4" />
                Clear
              </button>
            </div>
            
            <textarea
              v-model="inputText"
              placeholder="Enter URL, text, WiFi credentials, or any content to generate QR code..."
              class="w-full h-48 bg-black/40 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500/50 focus:ring-2 focus:ring-indigo-500/20 transition-all duration-200 resize-none font-mono text-sm"
            ></textarea>

            <!-- Character Count -->
            <div class="mt-3 flex items-center justify-between text-sm">
              <span class="text-gray-400">
                {{ statistics.textLength }} characters
                <span v-if="statistics.estimatedCapacity > 0">
                  / ~{{ statistics.estimatedCapacity }} max (numeric)
                </span>
              </span>
              <span
                v-if="statistics.textLength > statistics.estimatedCapacity"
                class="text-red-400 font-medium"
              >
                ⚠️ May exceed capacity
              </span>
            </div>
          </div>

          <!-- Quick Presets -->
          <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6 shadow-xl">
            <h3 class="text-lg font-bold text-white mb-4">Quick Templates</h3>
            <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
              <button
                v-for="preset in presets"
                :key="preset.name"
                @click="applyPreset(preset)"
                class="flex flex-col items-center gap-2 p-4 rounded-xl bg-white/5 hover:bg-indigo-600/20 border border-white/10 hover:border-indigo-500/30 transition-all duration-200 group"
              >
                <span class="text-2xl group-hover:scale-110 transition-transform duration-200">{{ preset.icon }}</span>
                <span class="text-sm font-medium text-gray-300 group-hover:text-white transition-colors duration-200">{{ preset.name }}</span>
              </button>
            </div>
          </div>

          <!-- Advanced Options -->
          <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6 shadow-xl">
            <button
              @click="showAdvancedOptions = !showAdvancedOptions"
              class="w-full flex items-center justify-between mb-4"
            >
              <div class="flex items-center gap-2">
                <Settings class="w-5 h-5 text-indigo-400" />
                <h3 class="text-lg font-bold text-white">Advanced Options</h3>
              </div>
              <span
                class="text-gray-400 transition-transform duration-200"
                :class="{ 'rotate-180': showAdvancedOptions }"
              >▼</span>
            </button>

            <div
              v-show="showAdvancedOptions"
              class="space-y-4 animate-in fade-in slide-in-from-top duration-300"
            >
              <!-- Error Correction Level -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">
                  Error Correction Level
                </label>
                <select
                  v-model="qrOptions.errorCorrectionLevel"
                  class="w-full bg-black/40 border border-white/10 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-indigo-500/50 focus:ring-2 focus:ring-indigo-500/20 transition-all duration-200"
                >
                  <option value="L">Low (7% recovery)</option>
                  <option value="M">Medium (15% recovery)</option>
                  <option value="Q">Quartile (25% recovery)</option>
                  <option value="H">High (30% recovery)</option>
                </select>
              </div>

              <!-- QR Code Size -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">
                  Size: {{ qrOptions.width }}px
                </label>
                <input
                  type="range"
                  v-model.number="qrOptions.width"
                  min="200"
                  max="800"
                  step="50"
                  class="w-full accent-indigo-600"
                />
              </div>

              <!-- Margin -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">
                  Margin: {{ qrOptions.margin }}
                </label>
                <input
                  type="range"
                  v-model.number="qrOptions.margin"
                  min="0"
                  max="10"
                  step="1"
                  class="w-full accent-indigo-600"
                />
              </div>

              <!-- Color Presets -->
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">
                  <Palette class="w-4 h-4 inline mr-1" />
                  Color Themes
                </label>
                <div class="grid grid-cols-4 gap-2">
                  <button
                    v-for="preset in colorPresets"
                    :key="preset.name"
                    @click="applyColorPreset(preset)"
                    class="aspect-square rounded-lg border-2 transition-all duration-200 hover:scale-105"
                    :style="{
                      background: `linear-gradient(135deg, ${preset.dark} 50%, ${preset.light} 50%)`,
                      borderColor: qrOptions.color.dark === preset.dark && qrOptions.color.light === preset.light ? '#6366F1' : 'rgba(255,255,255,0.1)'
                    }"
                    :title="preset.name"
                  ></button>
                </div>
              </div>

              <!-- Custom Colors -->
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-300 mb-2">
                    Dark Color
                  </label>
                  <input
                    type="color"
                    v-model="qrOptions.color.dark"
                    class="w-full h-10 bg-black/40 border border-white/10 rounded-lg cursor-pointer"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-300 mb-2">
                    Light Color
                  </label>
                  <input
                    type="color"
                    v-model="qrOptions.color.light"
                    class="w-full h-10 bg-black/40 border border-white/10 rounded-lg cursor-pointer"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column: QR Code Display and Statistics -->
        <div class="space-y-6 animate-in fade-in slide-in-from-right duration-700">
          <!-- QR Code Display -->
          <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6 shadow-xl">
            <h2 class="text-xl font-bold text-white mb-4">Generated QR Code</h2>
            
            <div
              v-if="qrDataURL"
              class="relative group"
            >
              <!-- QR Code Image -->
              <div class="bg-white rounded-2xl p-8 flex items-center justify-center mb-4">
                <img
                  :src="qrDataURL"
                  alt="Generated QR Code"
                  class="max-w-full h-auto transition-transform duration-300 group-hover:scale-105"
                  :style="{ width: qrOptions.width + 'px', maxWidth: '100%' }"
                />
              </div>

              <!-- Action Buttons -->
              <div class="grid grid-cols-2 gap-3">
                <button
                  @click="downloadQRCode"
                  class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-white font-semibold shadow-lg shadow-indigo-600/20 hover:shadow-indigo-600/40 transition-all duration-200"
                >
                  <Download class="w-5 h-5" />
                  Download
                </button>
                
                <button
                  @click="copyQRCode"
                  class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl bg-white/10 hover:bg-white/20 text-white font-semibold border border-white/10 hover:border-white/20 transition-all duration-200"
                >
                  <Check v-if="copySuccess" class="w-5 h-5 text-green-400" />
                  <Copy v-else class="w-5 h-5" />
                  {{ copySuccess ? 'Copied!' : 'Copy' }}
                </button>
              </div>

              <!-- Generation Time -->
              <div
                v-if="generatedAt"
                class="mt-3 text-center text-sm text-gray-400"
              >
                Generated at {{ formatTime(generatedAt) }}
              </div>
            </div>

            <!-- Placeholder when no QR code -->
            <div
              v-else
              class="flex flex-col items-center justify-center py-16 text-center"
            >
              <div class="w-24 h-24 rounded-full bg-white/5 flex items-center justify-center mb-4">
                <QrCode class="w-12 h-12 text-gray-600" />
              </div>
              <p class="text-gray-500 font-medium">
                Enter text to generate QR code
              </p>
            </div>
          </div>

          <!-- Statistics -->
          <div
            v-if="qrDataURL"
            class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6 shadow-xl"
          >
            <h3 class="text-lg font-bold text-white mb-4">Statistics</h3>
            
            <div class="grid grid-cols-2 gap-4">
              <!-- Text Length -->
              <div class="bg-black/30 rounded-xl p-4">
                <div class="text-sm text-gray-400 mb-1">Content Length</div>
                <div class="text-2xl font-bold text-white">
                  {{ statistics.textLength }}
                </div>
                <div class="text-xs text-gray-500 mt-1">characters</div>
              </div>

              <!-- File Size -->
              <div class="bg-black/30 rounded-xl p-4">
                <div class="text-sm text-gray-400 mb-1">Image Size</div>
                <div class="text-2xl font-bold text-white">
                  {{ formatSize(statistics.qrSize) }}
                </div>
                <div class="text-xs text-gray-500 mt-1">PNG format</div>
              </div>

              <!-- Error Correction -->
              <div class="bg-black/30 rounded-xl p-4">
                <div class="text-sm text-gray-400 mb-1">Error Correction</div>
                <div class="text-lg font-bold text-white">
                  {{ statistics.errorCorrection }}
                </div>
                <div class="text-xs text-gray-500 mt-1">recovery rate</div>
              </div>

              <!-- Dimensions -->
              <div class="bg-black/30 rounded-xl p-4">
                <div class="text-sm text-gray-400 mb-1">Dimensions</div>
                <div class="text-lg font-bold text-white">
                  {{ qrOptions.width }}×{{ qrOptions.width }}
                </div>
                <div class="text-xs text-gray-500 mt-1">pixels</div>
              </div>
            </div>

            <!-- Additional Info -->
            <div class="mt-4 p-4 bg-indigo-600/10 border border-indigo-500/20 rounded-xl">
              <div class="flex items-start gap-2">
                <span class="text-indigo-400 text-lg">💡</span>
                <div class="text-sm text-gray-300">
                  <strong class="text-white">Tip:</strong> Higher error correction levels allow the QR code to be scanned even if partially damaged or obscured, but reduce data capacity.
                </div>
              </div>
            </div>
          </div>

          <!-- Usage Examples -->
          <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl p-6 shadow-xl">
            <h3 class="text-lg font-bold text-white mb-4">Usage Examples</h3>
            
            <div class="space-y-3 text-sm">
              <div class="p-3 bg-black/30 rounded-lg">
                <div class="text-indigo-400 font-medium mb-1">📱 WiFi Connection</div>
                <code class="text-gray-300 text-xs">WIFI:T:WPA;S:YourNetwork;P:YourPassword;;</code>
              </div>
              
              <div class="p-3 bg-black/30 rounded-lg">
                <div class="text-purple-400 font-medium mb-1">📧 Email Message</div>
                <code class="text-gray-300 text-xs">mailto:email@example.com?subject=Hello&body=Message</code>
              </div>
              
              <div class="p-3 bg-black/30 rounded-lg">
                <div class="text-blue-400 font-medium mb-1">📞 Phone Call</div>
                <code class="text-gray-300 text-xs">tel:+1234567890</code>
              </div>
              
              <div class="p-3 bg-black/30 rounded-lg">
                <div class="text-green-400 font-medium mb-1">💬 SMS Message</div>
                <code class="text-gray-300 text-xs">sms:+1234567890?body=Your message here</code>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar for textarea */
textarea::-webkit-scrollbar {
  width: 6px;
}

textarea::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
}

textarea::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
}

textarea::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* Smooth animations */
@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slide-in-from-left {
  from {
    transform: translateX(-20px);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slide-in-from-right {
  from {
    transform: translateX(20px);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slide-in-from-top {
  from {
    transform: translateY(-10px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.animate-in {
  animation-fill-mode: both;
}

.fade-in {
  animation-name: fade-in;
}

.slide-in-from-left {
  animation-name: slide-in-from-left;
}

.slide-in-from-right {
  animation-name: slide-in-from-right;
}

.slide-in-from-top {
  animation-name: slide-in-from-top;
}

.duration-300 {
  animation-duration: 300ms;
}

.duration-700 {
  animation-duration: 700ms;
}
</style>
