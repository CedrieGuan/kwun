# Kwun Tools

A modern, modular web application featuring AI-powered utilities and tools.

## 🚀 Features

### 🛠️ Tools & Modules

1. **Image Analyzer** - Upload and analyze images with AI-powered vision capabilities
2. **JSON Formatter** - Format, validate, minify, and analyze JSON data
3. **QR Code Generator** - Create customizable QR codes with various options
4. **Base64 Tool** - Encode and decode text to/from Base64 format

### 🤖 Global AI Assistant

A persistent AI chat assistant accessible from any tool via a floating button in the bottom-right corner. Ask questions, get help, or chat with AI while using any module.

## 🎨 Design

- **Modern Dark Theme** - Sleek, professional dark mode interface
- **Glassmorphism Effects** - Frosted glass aesthetics with backdrop blur
- **Gradient Accents** - Vibrant indigo-purple gradients throughout
- **Smooth Animations** - Polished micro-interactions and transitions
- **Responsive Layout** - Works seamlessly on desktop and mobile devices

## 🏗️ Architecture

### Modular System

The application uses a configuration-driven modular architecture that makes adding new tools incredibly simple:

```javascript
// Just add a config entry in src/modules/index.js
{
  id: 'my-tool',
  name: 'My Tool',
  route: 'my-tool',
  icon: MyIcon,
  component: MyToolComponent,
  description: 'What this tool does'
}
```

See [MODULE_SYSTEM.md](./docs/MODULE_SYSTEM.md) for detailed instructions.

### Tech Stack

- **Frontend Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Icons**: Lucide Vue Icons
- **Backend**: Go (for AI chat API)
- **AI Integration**: OpenRouter API

## 📦 Installation

```bash
# Install dependencies
npm install

# Run development server (frontend only)
npm run dev

# Run development with Go backend
npm run dev:all
```

## 🔧 Development

### Project Structure

```
kwun/
├── src/
│   ├── components/          # Vue components
│   │   ├── ImageAnalyzer.vue
│   │   ├── JSONFormatter.vue
│   │   ├── QRCodeGenerator.vue
│   │   ├── Base64Tool.vue
│   │   └── ChatDialog.vue   # Global AI assistant
│   ├── modules/
│   │   └── index.js         # Module configuration
│   ├── App.vue              # Main app component
│   ├── main.js              # App entry point
│   └── style.css            # Global styles
├── api/                     # Go backend (AI chat)
├── docs/                    # Documentation
│   ├── MODULE_SYSTEM.md
│   └── CHATDIALOG_RESTORATION.md
└── REFACTORING_CHANGELOG.md
```

### Adding a New Tool

1. Create your Vue component in `src/components/`
2. Add a config entry in `src/modules/index.js`
3. Done! Your tool is now available

See [MODULE_SYSTEM.md](./docs/MODULE_SYSTEM.md) for a complete guide.

## 📖 Documentation

- [Module System Guide](./docs/MODULE_SYSTEM.md) - How to add new modules
- [ChatDialog Documentation](./docs/CHATDIALOG_RESTORATION.md) - AI assistant details
- [Refactoring Changelog](./REFACTORING_CHANGELOG.md) - Project history and changes

## 🎯 Development Scripts

```bash
# Frontend development server
npm run dev

# Frontend only (custom port)
npm run dev:fe

# Go backend only
npm run dev:go

# Run both frontend and backend concurrently
npm run dev:all

# Build for production
npm run build

# Preview production build
npm run preview
```

## 🌐 Routes

Access tools via hash-based routing:

- `/#image-analyzer` - Image analysis tool
- `/#json-formatter` - JSON formatting tool
- `/#qr-generator` - QR code generator
- `/#base64-tool` - Base64 encoder/decoder

## 🤖 AI Chat Setup

The AI chat assistant requires a backend API. Configure your Go server with:

1. OpenRouter API key in `.env`
2. Chat endpoint at `/api/chat`
3. CORS enabled for local development

See [ChatDialog Documentation](./docs/CHATDIALOG_RESTORATION.md) for API details.

## 📝 Recent Changes

### 2025-12-24: Major Refactoring

- ✅ Removed OneLink module
- ✅ Implemented modular configuration system
- ✅ Simplified App.vue (230 → 90 lines)
- ✅ Restored ChatDialog as global AI assistant
- ✅ Added Base64 Tool as demo module
- ✅ Created comprehensive documentation

See [REFACTORING_CHANGELOG.md](./REFACTORING_CHANGELOG.md) for details.

## 🎨 Design Philosophy

1. **Premium First** - Every interaction should feel polished and professional
2. **Developer Experience** - Adding features should be trivial
3. **User Delight** - Smooth animations and intuitive interfaces
4. **Consistency** - Unified design language across all tools
5. **Performance** - Fast load times and smooth interactions

## 🚀 Future Enhancements

- [ ] More tool modules (calculator, color picker, markdown editor)
- [ ] Module lazy loading for better performance
- [ ] Local storage for settings and preferences
- [ ] Export/import functionality for data
- [ ] Keyboard shortcuts for power users
- [ ] Module categories and search
- [ ] PWA support for offline usage
- [ ] Themes and customization options

## 📄 License

MIT License - feel free to use this project as you wish.

## 🙏 Credits

- Built with Vue 3 and Vite
- Icons by Lucide
- Styled with Tailwind CSS
- AI powered by OpenRouter

---

**Built with ❤️ by Kwun**
