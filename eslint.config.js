import js from "@eslint/js";
import globals from "globals";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import tseslint from "typescript-eslint";
import eslintPluginPrettier from "eslint-plugin-prettier/recommended";

export default tseslint
  .config(
    { ignores: ["dist"] },
    {
      extends: [js.configs.recommended, ...tseslint.configs.recommended],
      files: ["**/*.{ts,tsx}"],
      languageOptions: {
        ecmaVersion: 2020,
        globals: globals.browser,
      },
      plugins: {
        "react-hooks": reactHooks,
        "react-refresh": reactRefresh,
      },
      rules: {
        ...reactHooks.configs.recommended.rules,
        "react-refresh/only-export-components": [
          "warn",
          { allowConstantExport: true },
        ],
        semi: ["error", "never"],
        quotes: ["error", "single"],
        indent: ["error", 2],
        "prettier/prettier": [
          "error",
          {
            semi: false,
            singleQuote: true,
            tabWidth: 2,
            trailingComma: "none",
            endOfLine: "lf",
            printWidth: 100,
            arrowParens: "avoid",
            bracketSpacing: true,
            jsxSingleQuote: true,
            jsxBracketSameLine: false,
            embeddedLanguageFormatting: "auto",
          },
        ],
        "@typescript-eslint/explicit-module-boundary-types": "off",
      },
    },
  )
  .concat(eslintPluginPrettier);
