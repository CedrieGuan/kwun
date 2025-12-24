# Module System

This project uses a modular architecture that makes it easy to add new feature modules.

## Current Modules

- **Image Analyzer**: Analyze images with AI-powered vision capabilities
- **JSON Formatter**: Format, validate, and analyze JSON data
- **QR Generator**: Generate QR codes for URLs and text

## How to Add a New Module

Adding a new module is a simple 3-step process:

### Step 1: Create Your Component

Create a new Vue component in `src/components/`. For example:

```vue
<!-- src/components/MyNewTool.vue -->
<template>
  <div class="container mx-auto px-4 py-12">
    <h1>My New Tool</h1>
    <!-- Your tool logic here -->
  </div>
</template>

<script setup>
// Your Vue logic here
</script>
```

### Step 2: Register Your Module

Open `src/modules/index.js` and add your module to the `modules` array:

```javascript
import { FileImage, Braces, QrCode, YourIcon } from 'lucide-vue-next';
import ImageAnalyzer from '../components/ImageAnalyzer.vue';
import JSONFormatter from '../components/JSONFormatter.vue';
import QRCodeGenerator from '../components/QRCodeGenerator.vue';
import MyNewTool from '../components/MyNewTool.vue'; // Import your component

export const modules = [
  // ... existing modules ...
  {
    id: 'my-new-tool',              // Unique identifier
    name: 'My New Tool',            // Display name in navigation
    route: 'my-new-tool',           // URL hash route
    icon: YourIcon,                 // Lucide icon component
    component: MyNewTool,           // Your Vue component
    description: 'Description here' // Optional description
  }
];
```

### Step 3: Done!

That's it! Your new module will automatically:
- Appear in the navigation menu
- Be accessible via the URL hash `#my-new-tool`
- Use the same styling and layout as other tools

## Module Configuration

Each module object has the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | ✅ | Unique identifier for the module |
| `name` | string | ✅ | Display name shown in navigation |
| `route` | string | ✅ | URL hash route (e.g., 'my-tool') |
| `icon` | Component | ✅ | Lucide Vue icon component |
| `component` | Component | ✅ | Vue component to render |
| `description` | string | ⚪ | Optional module description |

## Helper Functions

The module system provides several helper functions in `src/modules/index.js`:

- **`getModuleByRoute(route)`**: Get a module by its route
- **`getAllRoutes()`**: Get all module routes
- **`isValidRoute(route)`**: Check if a route is valid

## Icon Options

You can use any icon from [Lucide Vue](https://lucide.dev/icons/). Import them like this:

```javascript
import { 
  FileImage, 
  Braces, 
  QrCode,
  Calculator,  // Math calculator
  Globe,       // Web tool
  Database,    // Data tool
  // ... and many more
} from 'lucide-vue-next';
```

## Best Practices

1. **Consistent Styling**: Use the existing dark theme with Tailwind classes
2. **Responsive Design**: Make sure your component works on mobile and desktop
3. **Unique Routes**: Ensure your route name doesn't conflict with existing modules
4. **Clear Names**: Use descriptive names that clearly indicate the tool's purpose
5. **Code Comments**: Add English comments to explain complex logic

## Example: Adding a Calculator

Here's a complete example of adding a calculator module:

```javascript
// 1. Create src/components/Calculator.vue
// 2. Add to src/modules/index.js:

import { Calculator } from 'lucide-vue-next';
import CalculatorTool from '../components/Calculator.vue';

export const modules = [
  // ... existing modules ...
  {
    id: 'calculator',
    name: 'Calculator',
    route: 'calculator',
    icon: Calculator,
    component: CalculatorTool,
    description: 'Simple calculator for basic math operations'
  }
];
```

Now users can access the calculator at `#calculator`!
