/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.{templ,go}"],
  theme: {
    extend: {
      colors: {
        background: "oklch(5.5% 0.018 292)",
        foreground: "oklch(98% 0 0)",
        card: {
          DEFAULT: "oklch(10.5% 0.018 292)",
          foreground: "oklch(98% 0 0)",
        },
        primary: {
          DEFAULT: "oklch(76% 0.145 292)",
          foreground: "oklch(12% 0.03 292)",
        },
        muted: {
          DEFAULT: "oklch(18% 0.018 292)",
          foreground: "oklch(72% 0.035 292)",
        },
        accent: {
          DEFAULT: "oklch(22% 0.04 292)",
          foreground: "oklch(98% 0 0)",
        },
        border: "oklch(22% 0.035 292)",
        input: "oklch(22% 0.035 292)",
        ring: "oklch(76% 0.145 292)",
      },
      borderRadius: {
        "2xl": "1rem",
        "3xl": "1.5rem",
      },
    },
  },
  plugins: [],
};
