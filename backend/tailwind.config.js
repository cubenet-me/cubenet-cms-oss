/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.{templ,go}"],
  theme: {
    extend: {
      colors: {
        base: "oklch(18% 0.04 280deg)",
        surface: {
          DEFAULT: "oklch(24% 0.04 280deg)",
          hover: "oklch(28% 0.04 280deg)",
          dimmed: "oklch(65% 0.02 280deg)",
          muted: "oklch(45% 0.02 280deg)",
        },
        accent: {
          DEFAULT: "oklch(65% 0.2 300deg)",
          hover: "oklch(70% 0.2 300deg)",
        },
        border: "oklch(30% 0.04 280deg)",
      },
    },
  },
  plugins: [],
};
