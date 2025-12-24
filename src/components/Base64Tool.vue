<template>
  <div class="min-h-screen bg-[#0a0a0b] text-white">
    <div class="container mx-auto px-4 py-12">
      <!-- Header -->
      <div class="max-w-4xl mx-auto space-y-8">
        <div class="text-center space-y-4">
          <div class="inline-flex items-center justify-center p-4 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-2xl shadow-lg shadow-indigo-600/20">
            <Binary class="w-12 h-12 text-white" />
          </div>
          <h1 class="text-5xl font-extrabold bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
            Base64 Encoder/Decoder
          </h1>
          <p class="text-gray-400 text-lg">
            Encode and decode text to/from Base64 format
          </p>
        </div>

        <!-- Main Content -->
        <div class="grid md:grid-cols-2 gap-6">
          <!-- Encode Section -->
          <div class="space-y-4">
            <div class="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6 space-y-4">
              <div class="flex items-center justify-between">
                <h2 class="text-xl font-semibold text-white">Encode to Base64</h2>
                <ArrowRight class="w-5 h-5 text-indigo-400" />
              </div>
              
              <div class="space-y-2">
                <label class="text-sm font-medium text-gray-300">Plain Text Input</label>
                <textarea
                  v-model="plainText"
                  @input="encode"
                  placeholder="Enter text to encode..."
                  class="w-full h-40 px-4 py-3 bg-black/30 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none font-mono text-sm"
                ></textarea>
              </div>

              <div class="space-y-2">
                <div class="flex items-center justify-between">
                  <label class="text-sm font-medium text-gray-300">Base64 Output</label>
                  <button
                    v-if="encodedText"
                    @click="copyEncoded"
                    class="flex items-center gap-1 px-3 py-1 bg-indigo-600 hover:bg-indigo-500 rounded-lg text-xs font-medium transition-colors"
                  >
                    <Copy class="w-3 h-3" />
                    {{ encodedCopied ? 'Copied!' : 'Copy' }}
                  </button>
                </div>
                <textarea
                  v-model="encodedText"
                  readonly
                  placeholder="Base64 output will appear here..."
                  class="w-full h-40 px-4 py-3 bg-black/30 border border-white/10 rounded-xl text-indigo-300 placeholder-gray-500 focus:outline-none resize-none font-mono text-sm"
                ></textarea>
              </div>
            </div>
          </div>

          <!-- Decode Section -->
          <div class="space-y-4">
            <div class="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6 space-y-4">
              <div class="flex items-center justify-between">
                <h2 class="text-xl font-semibold text-white">Decode from Base64</h2>
                <ArrowLeft class="w-5 h-5 text-purple-400" />
              </div>
              
              <div class="space-y-2">
                <label class="text-sm font-medium text-gray-300">Base64 Input</label>
                <textarea
                  v-model="base64Input"
                  @input="decode"
                  placeholder="Enter Base64 to decode..."
                  class="w-full h-40 px-4 py-3 bg-black/30 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none font-mono text-sm"
                ></textarea>
              </div>

              <div class="space-y-2">
                <div class="flex items-center justify-between">
                  <label class="text-sm font-medium text-gray-300">Plain Text Output</label>
                  <button
                    v-if="decodedText && !decodeError"
                    @click="copyDecoded"
                    class="flex items-center gap-1 px-3 py-1 bg-purple-600 hover:bg-purple-500 rounded-lg text-xs font-medium transition-colors"
                  >
                    <Copy class="w-3 h-3" />
                    {{ decodedCopied ? 'Copied!' : 'Copy' }}
                  </button>
                </div>
                <textarea
                  v-model="decodedText"
                  readonly
                  :class="[
                    'w-full h-40 px-4 py-3 border rounded-xl focus:outline-none resize-none font-mono text-sm',
                    decodeError 
                      ? 'bg-red-500/10 border-red-500/30 text-red-400' 
                      : 'bg-black/30 border-white/10 text-purple-300 placeholder-gray-500'
                  ]"
                  :placeholder="decodeError || 'Decoded text will appear here...'"
                ></textarea>
              </div>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="flex items-center justify-center gap-4">
          <button
            @click="clearAll"
            class="flex items-center gap-2 px-6 py-3 bg-white/5 hover:bg-white/10 rounded-xl font-medium text-gray-300 hover:text-white transition-all duration-300"
          >
            <Trash2 class="w-4 h-4" />
            Clear All
          </button>
          
          <button
            @click="swapContent"
            :disabled="!encodedText"
            :class="[
              'flex items-center gap-2 px-6 py-3 rounded-xl font-medium transition-all duration-300',
              encodedText
                ? 'bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-white shadow-lg shadow-indigo-600/20'
                : 'bg-white/5 text-gray-500 cursor-not-allowed'
            ]"
          >
            <RefreshCw class="w-4 h-4" />
            Swap Encode ↔ Decode
          </button>
        </div>

        <!-- Info Section -->
        <div class="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6">
          <h3 class="text-lg font-semibold text-white mb-3">About Base64 Encoding</h3>
          <p class="text-gray-400 text-sm leading-relaxed">
            Base64 is a binary-to-text encoding scheme that represents binary data in an ASCII string format. 
            It's commonly used to encode data that needs to be stored or transferred over media designed to handle text. 
            This encoding helps ensure data remains intact without modification during transport.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { Binary, ArrowRight, ArrowLeft, Copy, Trash2, RefreshCw } from 'lucide-vue-next';

// State for encoding
const plainText = ref('');
const encodedText = ref('');
const encodedCopied = ref(false);

// State for decoding
const base64Input = ref('');
const decodedText = ref('');
const decodeError = ref('');
const decodedCopied = ref(false);

/**
 * Encode plain text to Base64
 */
const encode = () => {
  try {
    if (!plainText.value) {
      encodedText.value = '';
      return;
    }
    // Use btoa for encoding (browser built-in function)
    encodedText.value = btoa(unescape(encodeURIComponent(plainText.value)));
  } catch (error) {
    console.error('Encoding error:', error);
    encodedText.value = '';
  }
};

/**
 * Decode Base64 to plain text
 */
const decode = () => {
  try {
    if (!base64Input.value) {
      decodedText.value = '';
      decodeError.value = '';
      return;
    }
    // Use atob for decoding (browser built-in function)
    decodedText.value = decodeURIComponent(escape(atob(base64Input.value)));
    decodeError.value = '';
  } catch (error) {
    console.error('Decoding error:', error);
    decodedText.value = '';
    decodeError.value = 'Invalid Base64 input';
  }
};

/**
 * Copy encoded text to clipboard
 */
const copyEncoded = async () => {
  try {
    await navigator.clipboard.writeText(encodedText.value);
    encodedCopied.value = true;
    setTimeout(() => {
      encodedCopied.value = false;
    }, 2000);
  } catch (error) {
    console.error('Copy failed:', error);
  }
};

/**
 * Copy decoded text to clipboard
 */
const copyDecoded = async () => {
  try {
    await navigator.clipboard.writeText(decodedText.value);
    decodedCopied.value = true;
    setTimeout(() => {
      decodedCopied.value = false;
    }, 2000);
  } catch (error) {
    console.error('Copy failed:', error);
  }
};

/**
 * Clear all inputs and outputs
 */
const clearAll = () => {
  plainText.value = '';
  encodedText.value = '';
  base64Input.value = '';
  decodedText.value = '';
  decodeError.value = '';
  encodedCopied.value = false;
  decodedCopied.value = false;
};

/**
 * Swap encoded output to decode input for easy round-trip testing
 */
const swapContent = () => {
  if (encodedText.value) {
    base64Input.value = encodedText.value;
    decode();
  }
};
</script>
