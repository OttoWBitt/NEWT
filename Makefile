deploy-interno-zuado:
	scp -r ~/NEWT/newt-backend ottobitt@192.168.0.14:/home/ottobitt/NEWT
	ssh ottobitt@192.168.0.14 "cd NEWT/newt-backend && make build-and-run && docker ps"

deploy-externo-zuado:
	scp -r ~/NEWT/newt-backend ottobitt@newt.ottobittencourt.com:/home/ottobitt/NEWT
	ssh ottobitt@newt.ottobittencourt.com "cd NEWT/newt-backend && make build-and-run && docker ps"

deploy-int:
	git fetch origin
	git reset --hard origin/main
	git push deploy-int master
	ssh ottobitt@192.168.0.14 "cd NEWT/newt-backend && make build-and-run && docker ps"

deploy-ext:
	git fetch origin
	git reset --hard origin/main
	git push deploy-int master
	ssh ottobitt@192.168.0.14 "cd NEWT/newt-backend && make build-and-run && docker ps"

#git remote add deploy-ext ssh://ottobitt@newt.ottobittencourt.com/home/ottobitt/NEWT.git