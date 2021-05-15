
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Constants } from 'src/app/util/constants';
import { Response } from 'src/app/models/response.model'
import { Router } from '@angular/router';
import { Views } from 'src/app/util/views.enum';
import { User } from 'src/app/models/user.model';

@Injectable()
export class AuthService {

  constructor(
    private router: Router,
    private httpClient: HttpClient
    ) { }

  login(username: string, password: string): Observable<any> {
    return this.httpClient.post<Response>(`${Constants.BASE_URL}login`,
    JSON.stringify({ username: username, password: password }), { headers: this.getHeaders() })
    .pipe(
      map((response: any) => {
        const user = response;
        return user;
      }));
  }

  signUp(user: User) : Observable<User> {
    return this.httpClient.post<any>(`${Constants.BASE_URL}signup`, user, {headers: this.getHeaders()}).pipe(
      map((response: Response) => {
        const resp: User = response.data;
        return resp
      }));
  }

  recoverPassword(email: String) : Observable<any> {
    return this.httpClient.post<any>(`${Constants.BASE_URL}recover/${email}`, {headers: this.getHeaders()}).pipe(
      map((response: Response) => {
        const resp: any = response.data;
        return resp
      }));
  }

  redirectToLogin(){
    this.router.navigate([Views.login.url])
  }

  private getHeaders(): HttpHeaders {
    const headers = new HttpHeaders({'content-type': 'application/json', accept: 'application/json'});
    return headers;
  }

}
