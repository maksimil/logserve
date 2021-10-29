module.exports = {
  darkMode: "class",
  extract: {
    include: ["main.ts", "build.go"],
    extractors: [
      {
        extractor: (content) => ({
          classes: content.match(/[^<>"'`\s]*[^<>"'`\s:]/g) ?? [],
        }),
        extensions: ["ts", "go"],
      },
    ],
  },
};
