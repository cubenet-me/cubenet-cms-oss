/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.{templ,go}"],
  theme: {
    extend: {
      colors: {
        page: "#0a0514",
        foreground: "#f7f8ff",
        muted: "rgba(210, 213, 231, 0.7)",
        primary: "#7c5cff",
        "primary-2": "#5fd6ff",
        card: "#1a1325",
      },
      fontFamily: {
        body: ['"Minecraft Rus"', "Manrope", "Segoe UI", "sans-serif"],
        display: ['"Minecraft Rus"', '"Space Grotesk"', "Manrope", "sans-serif"],
      },
      borderRadius: {
        "2xl": "16px",
        "3xl": "28px",
      },
      boxShadow: {
        "neu": "18px 18px 36px rgba(6, 4, 12, 0.85), -18px -18px 36px rgba(44, 32, 66, 0.55)",
        "neu-sm": "8px 8px 20px rgba(6, 4, 12, 0.6), -8px -8px 20px rgba(44, 32, 66, 0.4)",
        "glass": "0 26px 70px rgba(0, 0, 0, 0.55)",
      },
    },
  },
  plugins: [],
};
