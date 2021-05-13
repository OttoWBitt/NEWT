
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Constants } from 'src/app/util/constants';
import { Response } from 'src/app/models/response.model'
import { Router } from '@angular/router';
import { Views } from 'src/app/util/views.enum';

@Injectable()
export class LoginService {

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

  redirectToLogin(){
    this.router.navigate([Views.login.url])
  }

  private getHeaders(): HttpHeaders {
    const headers = new HttpHeaders({'content-type': 'application/json', accept: 'application/json'});
    return headers;
  }

}
