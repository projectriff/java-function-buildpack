uri() {
  sed 's|gs://|https://storage.googleapis.com/|' "${ROOT}"/dependency/url
}

sha256() {
  shasum -a 256 "${ROOT}"/dependency/java-function-invoker-*.jar | cut -f 1 -d ' '
}
