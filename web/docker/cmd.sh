if [ ! -z "$API" ]
then
  INDEX_PATH="/usr/share/nginx/html/index.html"
  INDEX_TMP_PATH="${INDEX_PATH}.tmp"
  sed 's/<script>.*<\/script>/<script>env = { API: '$API' }<\/script>/g' $INDEX_PATH > $INDEX_TMP_PATH
  mv $INDEX_TMP_PATH $INDEX_PATH
fi

nginx -g 'daemon off;'
