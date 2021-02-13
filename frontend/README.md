https://tailwindcss.com/docs/installation

```bash
npm install @hotwired/turbo@latest stimulus@latest stimulus-use@latest tailwindcss@latest postcss@latest autoprefixer@latest
npm install -D clean-webpack-plugin@latest mini-css-extract-plugin@latest css-minimizer-webpack-plugin@latest css-loader@latest postcss-loader@latest @fullhuman/postcss-purgecss@latest purgecss-webpack-plugin@latest

cat <<EOF > postcss.config.js
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
EOF

npx tailwindcss init
```

```bash
npx webpack --mode='development'
```

## Typescript

https://webpack.js.org/guides/typescript/
https://www.typescriptlang.org/dt/search?search=

```bash
npm i --save-dev typescript ts-loader @types/node @types/webpack-env
```
