deploy-interno:
	scp -r ~/NEWT/newt-backend ottobitt@192.168.0.14:/home/ottobitt/NEWT
	ssh ottobitt@192.168.0.14 "cd NEWT/newt-backend && make build-and-run && docker ps"

deploy-externo:
	scp -r ~/NEWT/newt-backend ottobitt@newt.ottobittencourt.com:/home/ottobitt/NEWT
	ssh ottobitt@newt.ottobittencourt.com "cd NEWT/newt-backend && make build-and-run && docker ps"