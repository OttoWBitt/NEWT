deploy-interno-zuado:
	scp -r ~/NEWT/newt-backend ottobitt@192.168.0.14:/home/ottobitt/NEWT
	ssh ottobitt@192.168.0.14 "cd NEWT/newt-backend && make build-and-run && docker ps"

deploy-externo-zuado:
	scp -r ~/NEWT/newt-backend ottobitt@newt.ottobittencourt.com:/home/ottobitt/NEWT
	ssh ottobitt@newt.ottobittencourt.com "cd NEWT/newt-backend && make build-and-run && docker ps"

deploy-int:
	ssh ottobitt@192.168.0.14 "rm -rf NEWT && mkdir NEWT"
	git fetch origin
	git reset --hard origin/main
	git push deploy-int main
	ssh ottobitt@192.168.0.14 "cd NEWT/newt-backend && make build-and-run-all && docker ps"

deploy-ext:
	ssh ottobitt@newt.ottobittencourt.com "rm -rf NEWT && mkdir NEWT"
	git fetch origin
	git reset --hard origin/main
	git push deploy-int main
	ssh ottobitt@newt.ottobittencourt.com "cd NEWT/newt-backend && make build-and-run-all && docker ps"

#git remote add deploy-ext ssh://ottobitt@newt.ottobittencourt.com/home/ottobitt/NEWT.git