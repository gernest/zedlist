rm -rf _book
gitbook build
cd _book
git init
git add -A 
git commit -m "Update docs"
git push -f git@github.com:zedio/zedlist.git master:gh-pages