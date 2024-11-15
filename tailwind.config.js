/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./ui/**/*.tmpl', './cmd/web/**/*.go'],
  theme: {
    extend: {
      colors: {
        primary: '#7C3AED', // purple-500
        secondary: '#94A3B8', // slate-300
        accent: '#E2E8F0', // slate-100 for light accents
        background: '#1E293B', // slate-800 for a dark background
        surface: '#F1F5F9', // slate-50 for cards or sections
        success: '#10B981', // green-500 for success indications
        warning: '#F59E0B', // amber-500 for warnings or highlights
        danger: '#EF4444', // red-500 for error messages or alerts
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'monospace'], // using JetBrains Mono
        sans: ['Inter', 'sans-serif'], // a neutral sans-serif for contrast
      },
      spacing: {
        128: '32rem', // Custom larger spacing, good for section breaks
        144: '36rem',
      },
      borderRadius: {
        xl: '1.25rem', // Rounded for cards or buttons
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
};
