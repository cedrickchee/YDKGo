submodule-update:
	git submodule foreach git pull origin master

setup: submodule-update

clean:
	rm -rf public resources

check-hugo:
	@./check_hugo.sh

serve: check-hugo
	hugo server \
	--buildDrafts \
	--buildFuture

preview-build: setup clean check-hugo
	hugo \
	--baseURL $(DEPLOY_PRIME_URL) \
	--buildDrafts \
	--buildFuture \
	--minify

production-build: setup clean check-hugo
	hugo \
	--minify

update: submodule-update
	git add -A && git commit -m "Update submodule" && git push
