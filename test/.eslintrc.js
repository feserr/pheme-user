module.exports = {
  env: {
    browser: true,
    es2021: true,
  },

  extends: [
    'airbnb-base',
    'airbnb-typescript/base',
  ],

  parserOptions: {
    project: './tsconfig.json',
  },

  rules: {
    "prefer-destructuring": ["error", { "object": true, "array": false }]
  },
};
