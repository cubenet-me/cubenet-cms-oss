/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.{templ,go}"],
  theme: {
    extend: {
      colors: {
        background: "#030207",
        foreground: "oklch(0.98 0 0)",
        card: "oklch(0.105 0.018 292)",
        "card-foreground": "oklch(0.98 0 0)",
        primary: "oklch(0.76 0.145 292)",
        "primary-foreground": "oklch(0.12 0.03 292)",
        secondary: "oklch(0.15 0.018 292)",
        muted: "oklch(0.18 0.018 292)",
        "muted-foreground": "oklch(0.72 0.035 292)",
        accent: "oklch(0.22 0.04 292)",
        border: "oklch(0.22 0.035 292)",
        destructive: "oklch(0.577 0.245 27.325)",
      },
      fontFamily: {
        body: ['"Unbounded"', "sans-serif"],
        display: ['"Unbounded"', "sans-serif"],
        mono: ['"Space Mono"', "monospace"],
      },
      borderRadius: {
        "2xl": "1rem",
        "3xl": "1.5rem",
      },
    },
  },
  plugins: [],
};
