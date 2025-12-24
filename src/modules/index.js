/**
 * Module Configuration System
 * 
 * This file defines all available modules in the application.
 * To add a new module, simply add an entry to the modules array below.
 * 
 * Each module should have:
 * - id: unique identifier for the module
 * - name: display name in navigation
 * - route: URL hash route (e.g., 'image-analyzer')
 * - icon: Lucide icon component
 * - component: Vue component to render
 */

import { FileImage, Braces, QrCode, Binary } from 'lucide-vue-next';
import ImageAnalyzer from '../components/ImageAnalyzer.vue';
import JSONFormatter from '../components/JSONFormatter.vue';
import QRCodeGenerator from '../components/QRCodeGenerator.vue';
import Base64Tool from '../components/Base64Tool.vue';

/**
 * Module definitions
 * Add new modules here to automatically register them in the application
 */
export const modules = [
    {
        id: 'image-analyzer',
        name: 'Image Analyzer',
        route: 'image-analyzer',
        icon: FileImage,
        component: ImageAnalyzer,
        description: 'Analyze images with AI-powered vision capabilities'
    },
    {
        id: 'json-formatter',
        name: 'JSON Formatter',
        route: 'json-formatter',
        icon: Braces,
        component: JSONFormatter,
        description: 'Format, validate, and analyze JSON data'
    },
    {
        id: 'qr-generator',
        name: 'QR Generator',
        route: 'qr-generator',
        icon: QrCode,
        component: QRCodeGenerator,
        description: 'Generate QR codes for URLs and text'
    },
    {
        id: 'base64-tool',
        name: 'Base64 Tool',
        route: 'base64-tool',
        icon: Binary,
        component: Base64Tool,
        description: 'Encode and decode text to/from Base64 format'
    }
];

/**
 * Get module by route
 * @param {string} route - The route identifier
 * @returns {object|null} Module object or null if not found
 */
export function getModuleByRoute(route) {
    return modules.find(module => module.route === route) || null;
}

/**
 * Get all module routes
 * @returns {string[]} Array of all module routes
 */
export function getAllRoutes() {
    return modules.map(module => module.route);
}

/**
 * Check if a route is valid
 * @param {string} route - The route to check
 * @returns {boolean} True if route exists
 */
export function isValidRoute(route) {
    return modules.some(module => module.route === route);
}
