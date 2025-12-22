<script setup>
import { ref, computed } from 'vue';
import { Upload, FileImage, Info, X, Download } from 'lucide-vue-next';

// Reactive state for storing image information
const imageFile = ref(null);
const imageUrl = ref(null);
const imageInfo = ref(null);
const isDragging = ref(false);
const isAnalyzing = ref(false);

/**
 * Format bytes to human-readable string
 * @param {number} bytes - Size in bytes
 * @param {number} decimals - Number of decimal places
 * @returns {string} Formatted size string
 */
const formatBytes = (bytes, decimals = 2) => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
};

/**
 * Get color depth from image data
 * @param {ImageData} imageData - Canvas image data
 * @returns {number} Color depth in bits
 */
const getColorDepth = (imageData) => {
  // Check if image has alpha channel with varying transparency
  const { data } = imageData;
  let hasAlpha = false;
  
  for (let i = 3; i < data.length; i += 4) {
    if (data[i] < 255) {
      hasAlpha = true;
      break;
    }
  }
  
  // Return 32-bit for RGBA, 24-bit for RGB
  return hasAlpha ? 32 : 24;
};

/**
 * Calculate image entropy (complexity measure)
 * @param {ImageData} imageData - Canvas image data
 * @returns {number} Entropy value
 */
const calculateEntropy = (imageData) => {
  const { data } = imageData;
  const histogram = new Array(256).fill(0);
  
  // Build histogram for luminance values
  for (let i = 0; i < data.length; i += 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    // Calculate luminance using standard formula
    const luminance = Math.round(0.299 * r + 0.587 * g + 0.114 * b);
    histogram[luminance]++;
  }
  
  const totalPixels = data.length / 4;
  let entropy = 0;
  
  // Calculate Shannon entropy
  for (let count of histogram) {
    if (count > 0) {
      const probability = count / totalPixels;
      entropy -= probability * Math.log2(probability);
    }
  }
  
  return entropy.toFixed(4);
};

/**
 * Detect dominant colors in the image
 * @param {ImageData} imageData - Canvas image data
 * @param {number} sampleSize - Number of pixels to sample
 * @returns {Array} Array of dominant colors
 */
const getDominantColors = (imageData, sampleSize = 5) => {
  const { data, width, height } = imageData;
  const colorMap = new Map();
  
  // Sample pixels across the image
  const step = Math.floor((width * height) / 1000); // Sample ~1000 pixels
  
  for (let i = 0; i < data.length; i += step * 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    
    // Group similar colors (reduce precision)
    const colorKey = `${Math.floor(r / 32) * 32},${Math.floor(g / 32) * 32},${Math.floor(b / 32) * 32}`;
    colorMap.set(colorKey, (colorMap.get(colorKey) || 0) + 1);
  }
  
  // Sort by frequency and get top colors
  return Array.from(colorMap.entries())
    .sort((a, b) => b[1] - a[1])
    .slice(0, sampleSize)
    .map(([color]) => {
      const [r, g, b] = color.split(',').map(Number);
      return `rgb(${r}, ${g}, ${b})`;
    });
};

/**
 * Analyze uploaded image and extract metadata
 * @param {File} file - Image file to analyze
 */
const analyzeImage = async (file) => {
  isAnalyzing.value = true;
  
  try {
    // Create object URL for preview
    if (imageUrl.value) {
      URL.revokeObjectURL(imageUrl.value);
    }
    imageUrl.value = URL.createObjectURL(file);
    
    // Create image element to load the file
    const img = new Image();
    
    await new Promise((resolve, reject) => {
      img.onload = resolve;
      img.onerror = reject;
      img.src = imageUrl.value;
    });
    
    // Create canvas to analyze image data
    const canvas = document.createElement('canvas');
    canvas.width = img.naturalWidth;
    canvas.height = img.naturalHeight;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);
    
    // Get image data for analysis
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    
    // Calculate aspect ratio
    const gcd = (a, b) => b === 0 ? a : gcd(b, a % b);
    const divisor = gcd(img.naturalWidth, img.naturalHeight);
    const aspectWidth = img.naturalWidth / divisor;
    const aspectHeight = img.naturalHeight / divisor;
    
    // Get file extension and MIME type
    const extension = file.name.split('.').pop().toUpperCase();
    const mimeType = file.type;
    
    // Calculate megapixels
    const megapixels = ((img.naturalWidth * img.naturalHeight) / 1000000).toFixed(2);
    
    // Get color depth
    const colorDepth = getColorDepth(imageData);
    
    // Calculate entropy
    const entropy = calculateEntropy(imageData);
    
    // Get dominant colors
    const dominantColors = getDominantColors(imageData);
    
    // Get creation date if available from EXIF (basic approach)
    const lastModified = new Date(file.lastModified);
    
    // Store all analyzed information
    imageInfo.value = {
      // Basic information
      fileName: file.name,
      fileSize: file.size,
      fileSizeFormatted: formatBytes(file.size),
      mimeType: mimeType,
      extension: extension,
      
      // Image dimensions
      width: img.naturalWidth,
      height: img.naturalHeight,
      megapixels: megapixels,
      aspectRatio: `${aspectWidth}:${aspectHeight}`,
      
      // Color information
      colorDepth: colorDepth,
      dominantColors: dominantColors,
      
      // Advanced metrics
      entropy: entropy,
      
      // Metadata
      lastModified: lastModified.toLocaleString(),
      
      // Calculated properties
      totalPixels: img.naturalWidth * img.naturalHeight,
      estimatedMemorySize: formatBytes(img.naturalWidth * img.naturalHeight * 4), // RGBA
    };
    
    imageFile.value = file;
  } catch (error) {
    console.error('Error analyzing image:', error);
    alert('Failed to analyze image. Please make sure the file is a valid image.');
  } finally {
    isAnalyzing.value = false;
  }
};

/**
 * Handle file input change event
 * @param {Event} event - Input change event
 */
const handleFileSelect = async (event) => {
  const file = event.target.files?.[0];
  if (file && file.type.startsWith('image/')) {
    await analyzeImage(file);
  } else {
    alert('Please select a valid image file');
  }
};

/**
 * Handle drag over event
 * @param {DragEvent} event - Drag event
 */
const handleDragOver = (event) => {
  event.preventDefault();
  isDragging.value = true;
};

/**
 * Handle drag leave event
 */
const handleDragLeave = () => {
  isDragging.value = false;
};

/**
 * Handle file drop event
 * @param {DragEvent} event - Drop event
 */
const handleDrop = async (event) => {
  event.preventDefault();
  isDragging.value = false;
  
  const file = event.dataTransfer.files?.[0];
  if (file && file.type.startsWith('image/')) {
    await analyzeImage(file);
  } else {
    alert('Please drop a valid image file');
  }
};

/**
 * Reset analyzer state
 */
const resetAnalyzer = () => {
  if (imageUrl.value) {
    URL.revokeObjectURL(imageUrl.value);
  }
  imageFile.value = null;
  imageUrl.value = null;
  imageInfo.value = null;
};

/**
 * Download analyzed image
 */
const downloadImage = () => {
  if (imageFile.value) {
    const link = document.createElement('a');
    link.href = imageUrl.value;
    link.download = imageFile.value.name;
    link.click();
  }
};

// Computed property for file input trigger
const triggerFileInput = () => {
  document.getElementById('image-file-input').click();
};
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
            <FileImage class="w-8 h-8 text-white" />
          </div>
          <h1 class="text-5xl font-extrabold tracking-tight bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
            Image Analyzer
          </h1>
        </div>
        <p class="text-gray-400 text-lg max-w-2xl mx-auto">
          Upload a local image to instantly view detailed information including dimensions, format, file size, and more
        </p>
      </div>

      <!-- Upload Area -->
      <div 
        v-if="!imageInfo"
        class="max-w-2xl mx-auto animate-in slide-in-from-bottom duration-700"
        @dragover="handleDragOver"
        @dragleave="handleDragLeave"
        @drop="handleDrop"
      >
        <div 
          :class="[
            'relative border-2 border-dashed rounded-3xl p-16 text-center transition-all duration-300',
            isDragging 
              ? 'border-indigo-500 bg-indigo-500/10 scale-105' 
              : 'border-white/10 bg-white/5 hover:border-indigo-500/50 hover:bg-white/10'
          ]"
        >
          <input
            id="image-file-input"
            type="file"
            accept="image/*"
            class="hidden"
            @change="handleFileSelect"
          />
          
          <div class="space-y-6">
            <div class="flex justify-center">
              <div class="p-6 bg-indigo-600/20 rounded-full">
                <Upload class="w-16 h-16 text-indigo-400" />
              </div>
            </div>
            
            <div>
              <h3 class="text-2xl font-bold mb-2">
                {{ isDragging ? 'Drop your image here' : 'Upload Image' }}
              </h3>
              <p class="text-gray-400">
                Drag and drop or click to select an image file
              </p>
            </div>
            
            <button
              @click="triggerFileInput"
              class="px-8 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-xl font-semibold hover:from-indigo-500 hover:to-purple-500 transition-all duration-300 shadow-lg shadow-indigo-600/20 hover:shadow-indigo-600/40 hover:scale-105"
            >
              Select Image
            </button>
            
            <p class="text-xs text-gray-500">
              Supports: JPG, PNG, GIF, WebP, BMP, SVG
            </p>
          </div>
        </div>
      </div>

      <!-- Analysis Results -->
      <div v-else class="animate-in fade-in duration-700">
        <div class="grid lg:grid-cols-2 gap-8">
          <!-- Image Preview -->
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h2 class="text-2xl font-bold">Preview</h2>
              <div class="flex gap-2">
                <button
                  @click="downloadImage"
                  class="p-2 bg-white/5 hover:bg-white/10 rounded-lg transition-colors"
                  title="Download"
                >
                  <Download class="w-5 h-5" />
                </button>
                <button
                  @click="resetAnalyzer"
                  class="p-2 bg-white/5 hover:bg-white/10 rounded-lg transition-colors"
                  title="Upload New"
                >
                  <X class="w-5 h-5" />
                </button>
              </div>
            </div>
            
            <div class="relative bg-[#111113] border border-white/5 rounded-2xl overflow-hidden">
              <img
                :src="imageUrl"
                :alt="imageInfo.fileName"
                class="w-full h-auto max-h-[600px] object-contain"
              />
            </div>
          </div>

          <!-- Image Information -->
          <div class="space-y-4">
            <h2 class="text-2xl font-bold flex items-center gap-2">
              <Info class="w-6 h-6" />
              Detailed Information
            </h2>

            <div class="space-y-3">
              <!-- Basic Info Section -->
              <div class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-4">
                <h3 class="text-lg font-semibold text-indigo-400 border-b border-white/5 pb-2">
                  Basic Information
                </h3>
                
                <div class="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <div class="text-gray-400 mb-1">File Name</div>
                    <div class="font-medium break-all">{{ imageInfo.fileName }}</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-1">File Size</div>
                    <div class="font-medium">{{ imageInfo.fileSizeFormatted }}</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-1">Format</div>
                    <div class="font-medium">{{ imageInfo.extension }}</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-1">MIME Type</div>
                    <div class="font-medium text-xs">{{ imageInfo.mimeType }}</div>
                  </div>
                  
                  <div class="col-span-2">
                    <div class="text-gray-400 mb-1">Last Modified</div>
                    <div class="font-medium text-xs">{{ imageInfo.lastModified }}</div>
                  </div>
                </div>
              </div>

              <!-- Dimensions Section -->
              <div class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-4">
                <h3 class="text-lg font-semibold text-purple-400 border-b border-white/5 pb-2">
                  Dimensions
                </h3>
                
                <div class="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <div class="text-gray-400 mb-1">Width</div>
                    <div class="font-medium">{{ imageInfo.width.toLocaleString() }} px</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-1">Height</div>
                    <div class="font-medium">{{ imageInfo.height.toLocaleString() }} px</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-1">Aspect Ratio</div>
                    <div class="font-medium">{{ imageInfo.aspectRatio }}</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-1">Megapixels</div>
                    <div class="font-medium">{{ imageInfo.megapixels }} MP</div>
                  </div>
                  
                  <div class="col-span-2">
                    <div class="text-gray-400 mb-1">Total Pixels</div>
                    <div class="font-medium">{{ imageInfo.totalPixels.toLocaleString() }}</div>
                  </div>
                </div>
              </div>

              <!-- Color Information Section -->
              <div class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-4">
                <h3 class="text-lg font-semibold text-pink-400 border-b border-white/5 pb-2">
                  Color Information
                </h3>
                
                <div class="space-y-4 text-sm">
                  <div>
                    <div class="text-gray-400 mb-1">Color Depth</div>
                    <div class="font-medium">{{ imageInfo.colorDepth }} bits</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-2">Dominant Colors</div>
                    <div class="flex gap-2 flex-wrap">
                      <div
                        v-for="(color, index) in imageInfo.dominantColors"
                        :key="index"
                        class="flex items-center gap-2 bg-white/5 rounded-lg px-3 py-2"
                      >
                        <div
                          :style="{ backgroundColor: color }"
                          class="w-6 h-6 rounded border border-white/20"
                        ></div>
                        <span class="text-xs font-mono">{{ color }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Advanced Metrics Section -->
              <div class="bg-[#111113] border border-white/5 rounded-2xl p-6 space-y-4">
                <h3 class="text-lg font-semibold text-cyan-400 border-b border-white/5 pb-2">
                  Advanced Metrics
                </h3>
                
                <div class="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <div class="text-gray-400 mb-1">Entropy</div>
                    <div class="font-medium">{{ imageInfo.entropy }}</div>
                    <div class="text-xs text-gray-500 mt-1">Image complexity</div>
                  </div>
                  
                  <div>
                    <div class="text-gray-400 mb-1">Memory Size</div>
                    <div class="font-medium">{{ imageInfo.estimatedMemorySize }}</div>
                    <div class="text-xs text-gray-500 mt-1">Uncompressed RGBA</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Upload New Button -->
            <button
              @click="resetAnalyzer"
              class="w-full px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-xl font-semibold hover:from-indigo-500 hover:to-purple-500 transition-all duration-300 shadow-lg shadow-indigo-600/20 hover:shadow-indigo-600/40"
            >
              Analyze Another Image
            </button>
          </div>
        </div>
      </div>

      <!-- Loading Overlay -->
      <div
        v-if="isAnalyzing"
        class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50"
      >
        <div class="bg-[#111113] border border-white/10 rounded-2xl p-8 text-center">
          <div class="w-16 h-16 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p class="text-lg font-semibold">Analyzing Image...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
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
</style>
