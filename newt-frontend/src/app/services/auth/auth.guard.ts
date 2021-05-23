import { Injectable } from '@angular/core';
import { Router, ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, CanLoad, Route } from '@angular/router';
import { Observable } from 'rxjs';
import { Views } from 'src/app/util/views.enum';


@Injectable()
export class AuthGuard implements CanActivate {

  constructor(private router: Router) { }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
    return this.verifyPermission();
  }

  canLoad(route: Route): boolean | Observable<boolean> | Promise<boolean> {
    return this.verifyPermission();
  }

  verifyPermission(): boolean {

    if (localStorage.getItem('currentUser')) {
      return true;
    }
    // not logged in so redirect to login page with the return url
    this.router.navigate([Views.login.url]);
    return false;

  }
}
