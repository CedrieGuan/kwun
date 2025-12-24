<script setup>
import { ref, nextTick } from 'vue';
import { MessageSquare, Send, X, Loader2, Bot } from 'lucide-vue-next';

const isOpen = ref(false);
const message = ref('');
const messages = ref([]);
const isLoading = ref(false);
const chatContainer = ref(null);

/**
 * Scroll chat container to bottom to show latest messages
 */
const scrollToBottom = async () => {
  await nextTick();
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};

/**
 * Send user message to AI and get response
 */
const sendMessage = async () => {
  if (!message.value.trim() || isLoading.value) return;

  const userMsg = message.value;
  messages.value.push({ role: 'user', content: userMsg });
  message.value = '';
  isLoading.value = true;
  scrollToBottom();

  try {
    const response = await fetch('/api/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: userMsg
      })
    });

    const data = await response.json();

    if (!response.ok) {
      // Display the actual error message from the API
      const errorMsg = data.error || 'Failed to get response from AI';
      messages.value.push({ role: 'ai', content: `Error: ${errorMsg}` });
    } else {
      messages.value.push({ role: 'ai', content: data.reply });
    }
  } catch (error) {
    messages.value.push({ 
      role: 'ai', 
      content: 'Network error: Unable to connect to the AI service. Please try again later.' 
    });
  } finally {
    isLoading.value = false;
    scrollToBottom();
  }
};

/**
 * Clear all chat messages
 */
const clearChat = () => {
  messages.value = [];
};
</script>

<template>
  <div class="fixed bottom-6 right-6 z-50">
    <!-- Chat Toggle Button -->
    <button 
      @click="isOpen = !isOpen" 
      class="p-4 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-white rounded-full shadow-2xl transition-all transform hover:scale-110 active:scale-95 group"
    >
      <MessageSquare v-if="!isOpen" class="w-6 h-6" />
      <X v-else class="w-6 h-6" />
      <span class="absolute -top-12 right-0 bg-gradient-to-r from-indigo-600 to-purple-600 text-white text-xs px-3 py-1.5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none shadow-lg">
        Ask AI Assistant
      </span>
    </button>

    <!-- Chat Window -->
    <div 
      v-if="isOpen" 
      class="absolute bottom-20 right-0 w-[350px] sm:w-[400px] h-[500px] bg-[#111113] border border-white/10 rounded-3xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in slide-in-from-bottom-10 duration-300"
    >
      <!-- Header -->
      <div class="p-4 border-b border-white/10 bg-gradient-to-r from-indigo-600/10 to-purple-600/10 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="p-2 bg-gradient-to-r from-indigo-500/20 to-purple-500/20 rounded-xl">
            <Bot class="w-5 h-5 text-indigo-400" />
          </div>
          <div>
            <h3 class="font-bold text-white text-sm">AI Assistant</h3>
            <p class="text-[10px] text-gray-500">Here to help you</p>
          </div>
        </div>
        <button
          v-if="messages.length > 0"
          @click="clearChat"
          class="text-xs text-gray-400 hover:text-white transition-colors px-2 py-1 rounded-lg hover:bg-white/5"
        >
          Clear
        </button>
      </div>

      <!-- Messages -->
      <div ref="chatContainer" class="flex-1 overflow-y-auto p-4 space-y-4 custom-scrollbar">
        <div v-if="messages.length === 0" class="h-full flex flex-col items-center justify-center text-center space-y-3 opacity-30">
          <MessageSquare class="w-12 h-12" />
          <div>
            <p class="text-sm font-medium">AI Assistant Ready</p>
            <p class="text-xs text-gray-500 mt-1">Ask me anything!</p>
          </div>
        </div>

        <div 
          v-for="(msg, index) in messages" 
          :key="index" 
          :class="['flex', msg.role === 'user' ? 'justify-end' : 'justify-start']"
        >
          <div 
            :class="[
              'max-w-[80%] p-3 rounded-2xl text-sm leading-relaxed',
              msg.role === 'user' 
                ? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-tr-none shadow-lg' 
                : 'bg-white/5 text-gray-300 border border-white/5 rounded-tl-none'
            ]"
          >
            {{ msg.content }}
          </div>
        </div>

        <!-- Loading Indicator -->
        <div v-if="isLoading" class="flex justify-start">
          <div class="bg-white/5 text-gray-300 border border-white/5 p-3 rounded-2xl rounded-tl-none">
            <Loader2 class="w-4 h-4 animate-spin" />
          </div>
        </div>
      </div>

      <!-- Input -->
      <div class="p-4 border-t border-white/10 bg-white/5">
        <div class="relative flex items-center gap-2">
          <input 
            v-model="message" 
            type="text" 
            placeholder="Type a message..."
            @keyup.enter="sendMessage"
            class="w-full bg-[#1a1a1c] border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all" 
          />
          <button 
            @click="sendMessage" 
            :disabled="isLoading || !message.trim()"
            class="p-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-white rounded-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Send class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
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
</style>
