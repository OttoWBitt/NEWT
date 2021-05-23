
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Response } from 'src/app/models/response.model';
import { Subject } from 'src/app/models/subject.model';
import { Constants } from 'src/app/util/constants';

@Injectable()
export class SubjectService {

  constructor(private httpClient: HttpClient) { }
    
  getAllSubjects() : Observable<Subject[]> {
    return this.httpClient.get<any>(`${Constants.BASE_URL}subject/all`, {headers: this.getHeaders()}).pipe(
      map((response: Response) => {
        const resp: Subject[] = response.data;
        return resp;
      }));
  }

  private getHeaders(): HttpHeaders {
    const headers = new HttpHeaders();
    return headers;
  }

}